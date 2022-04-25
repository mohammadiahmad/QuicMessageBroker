package server

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/lucas-clemente/quic-go"
	"go.uber.org/ratelimit"
)

type Client struct {
	msgChan chan []byte
	status  bool
}

type QuicBroker struct {
	config  *Config
	clients map[string]Client
	mux     sync.RWMutex
}

type Config struct {
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	MessageCount int    `mapstructure:"msg_count"`
	RateLimit    int    `mapstructure:"rate_limit"`
}

func NewQuicBroker(cfg *Config) QuicBroker {
	return QuicBroker{
		config:  cfg,
		clients: make(map[string]Client),
	}
}

func (q *QuicBroker) Run() error {
	addr := fmt.Sprintf("%s:%s", q.config.Host, q.config.Port)
	listener, err := quic.ListenAddr(addr, q.generateTLSConfig(), nil)
	if err != nil {
		return err
	}
	fmt.Println("Server Started Successfully")

	for {
		sess, err := listener.Accept(context.Background())
		if err != nil {
			return err
		}
		log.Printf("session accepted: %s", sess.RemoteAddr().String())
		go func() {
			defer func() {
				_ = sess.CloseWithError(0, "bye")
				log.Printf("close session: %s", sess.RemoteAddr().String())
			}()
			err := q.communicate(sess)
			if err != nil {
				fmt.Println(err)
			}
		}()
	}
}

func (q *QuicBroker) randomDataProducer(clientId string) {
	q.mux.RLock()
	client := q.clients[clientId]
	q.mux.RUnlock()
	rt := ratelimit.New(q.config.RateLimit)
	for i := 0; i <= q.config.MessageCount; i++ {
		rt.Take()
		client.msgChan <- []byte(fmt.Sprintf("message from server to client %s \n", clientId))
	}

}

func (q *QuicBroker) communicate(sess quic.Connection) error {
	messageChan := make(chan []byte)
	for {
		stream, err := sess.AcceptStream(context.Background())
		if err != nil {
			return err
		}
		log.Printf("bidirectional stream accepted: %d", stream.StreamID())
		var clientID string
		for {
			clientID, err = receiveClientID(stream)
			if err != nil {
				return err
			}
			if clientID != "" {
				q.mux.Lock()
				q.clients[clientID] = Client{
					msgChan: messageChan,
					status:  true,
				}
				q.mux.Unlock()
				go q.randomDataProducer(clientID)
				break
			}

		}

		for {
			select {
			case s := <-messageChan:
				_, err := stream.Write(s)
				stream.SetWriteDeadline(time.Now().Add(1 * time.Second))
				if err != nil {
					q.mux.Lock()
					delete(q.clients, clientID)
					q.mux.Unlock()
					return err
				}
			}
		}

	}
}

func receiveClientID(stream quic.Stream) (string, error) {

	r := bufio.NewReader(io.Reader(stream))
	clientID, err := r.ReadString('\n')
	if len(clientID) > 0 && err == nil {
		return strings.TrimSuffix(clientID, "\n"), nil
	}

	return "", err

}

func (q *QuicBroker) generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-broker"},
	}
}
