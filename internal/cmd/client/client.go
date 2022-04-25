package client

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/mohammadiahmad/QuicMessageBroker/internal/client"
	"github.com/mohammadiahmad/QuicMessageBroker/internal/http"
	"github.com/spf13/cobra"
)

const (
	use   = "client"
	short = "Run Quic Client"
)

func Client() *cobra.Command {
	// nolint: exhaustivestruct
	cmd := &cobra.Command{Use: use, Short: short, Run: main}
	cmd.PersistentFlags().Int("count", 10, "concurrent client count")
	cmd.PersistentFlags().String("address", "localhost", "quic server address ")
	cmd.PersistentFlags().String("port", "4040", "quic server port ")
	return cmd
}

func main(cmd *cobra.Command, _ []string) {
	count, err := cmd.PersistentFlags().GetInt("count")
	if err != nil {
		log.Fatal("client count is required")
	}

	addr, err := cmd.PersistentFlags().GetString("address")
	if err != nil {
		log.Fatal("server address is required")
	}

	port, err := cmd.PersistentFlags().GetString("port")
	if err != nil {
		log.Fatal("server port is required")
	}

	go http.Run("0.0.0.0", "9001")
	s := client.NewClient(addr, port)
	for i := 0; i < count; i++ {
		var err error
		clientId := i
		go func() {
			err = s.Run(clientId)
		}()
		if err != nil {
			fmt.Printf("error in run client number %s:%v\n", strconv.Itoa(i), err)
		}
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("signal", (<-signalChannel).String())

}
