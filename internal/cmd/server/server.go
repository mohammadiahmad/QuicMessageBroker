package server

import (
	"fmt"
	"github.com/mohammadiahmad/QuicMessageBroker/internal/server"
	"github.com/spf13/cobra"
)

const (
	use   = "server"
	short = "Run Quic Server"
)

// Server adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Server() *cobra.Command {
	// nolint: exhaustivestruct
	cmd := &cobra.Command{Use: use, Short: short, Run: main}

	return cmd
}

func main(cmd *cobra.Command, _ []string) {
	cfg := server.Config{
		"localhost",
		"4040",
	}
	s := server.NewQuicBroker(cfg)
	err := s.Run()
	fmt.Println(err)

}
