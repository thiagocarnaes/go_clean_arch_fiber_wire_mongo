package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "initApiServer",
	Short: "User Management API",
	Run: func(cmd *cobra.Command, args []string) {

		server, err := InitializeServer()
		if err != nil {
			log.Fatalf("Failed to initialize server: %v", err)
		}
		if err := server.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
