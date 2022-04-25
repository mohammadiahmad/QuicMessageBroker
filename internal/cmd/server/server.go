package server

import (
	"fmt"

	"github.com/mohammadiahmad/QuicMessageBroker/internal/config"
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
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("cannot load config")
		panic(err)
	}

	s := server.NewQuicBroker(cfg.Server)
	err = s.Run()
	fmt.Println("cannot run quic server", err)

}
