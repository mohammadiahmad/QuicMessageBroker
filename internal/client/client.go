package client

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/mohammadiahmad/QuicMessageBroker/internal/client/metrics"
)

type Client struct {
	host   string
	port   string
	metric metrics.ClientMetrics
}

func NewClient(host string, port string) Client {
	m := metrics.NewClientMetrics("Dispatching")
	return Client{
		host:   host,
		port:   port,
		metric: m,
	}
}

func (c *Client) Run(clientID int) error {
	addr := fmt.Sprintf("%s:%s", c.host, c.port)
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-broker"},
	}
	conn, err := quic.DialAddr(addr, tlsConf, nil)
	if err != nil {
		fmt.Println("cannot connect to quic server", err)
		return err
	}

	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		fmt.Println("cannot open stream", err)
		return err
	}

	//send client id to quic server
	_, err = stream.Write([]byte(strconv.Itoa(clientID) + "\n"))
	if err != nil {
		fmt.Println("cannot send client id to Quic server")
		return err
	}

	//Listen to server messages
	r := bufio.NewReader(io.Reader(stream))
	for {
		stream.SetReadDeadline(time.Now().Add(5 * time.Second))
		s, err := r.ReadString('\n')
		if err != nil {
			fmt.Println("stream read deadline exceeded", err)
			return err
		}
		c.metric.IncReceive("TotalMessagesReceived")
		c.metric.IncReceive(strconv.Itoa(clientID))
		fmt.Printf("Client: Got %s", s)
	}
}
