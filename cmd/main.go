package main

import (
	"github.com/mohammadiahmad/QuicMessageBroker/internal/cmd/client"
	"log"

	"github.com/mohammadiahmad/QuicMessageBroker/internal/cmd/server"
	"github.com/spf13/cobra"
)

const (
	short = "QuicMessageBroker"
	long  = `Quic based message broker`
)

func main() {
	// nolint: exhaustivestruct
	root := &cobra.Command{Short: short, Long: long}
	root.AddCommand(server.Server())
	root.AddCommand(client.Client())

	if err := root.Execute(); err != nil {
		log.Println(err.Error())
	}
}
