package client

import (
	"github.com/mohammadiahmad/QuicMessageBroker/internal/client"
	"github.com/spf13/cobra"
)

const (
	use   = "client"
	short = "Run Quic Client"
)

func Client() *cobra.Command {
	// nolint: exhaustivestruct
	cmd := &cobra.Command{Use: use, Short: short, Run: main}

	return cmd
}

func main(cmd *cobra.Command, _ []string) {

	s := client.NewClient("localhost", "4040")
	for i := 0; i < 10000; i++ {
		go s.Run()
	}

	select {}

}
