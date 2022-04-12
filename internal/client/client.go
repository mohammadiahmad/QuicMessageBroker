package client

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/lucas-clemente/quic-go"
	"io"
	"math/rand"
	"strconv"
)

type Client struct {
	host string
	port string
}

func NewClient(host string, port string) Client {
	return Client{
		host: host,
		port: port,
	}
}

func (c *Client) Run() error {
	addr := fmt.Sprintf("%s:%s", c.host, c.port)
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-broker"},
	}
	conn, err := quic.DialAddr(addr, tlsConf, nil)
	if err != nil {
		return err
	}

	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		return err
	}

	randomClientId := rand.Int()
	_, err = stream.Write([]byte(strconv.Itoa(randomClientId) + "\n"))
	if err != nil {
		return err
	}

	//Listen to server messages
	r := bufio.NewReader(io.Reader(stream))
	for {
		s, err := r.ReadString('\n')
		if err != nil {
			return err
		}
		fmt.Printf("Client: Got %s", s)
	}
}
