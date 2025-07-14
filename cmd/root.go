package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "user-management",
	Short: "User Management API",
	Run: func(cmd *cobra.Command, args []string) {
		//cfg, err := config.LoadConfig()
		//if err != nil {
		//	log.Fatalf("Failed to load config: %v", err)
		//}
		//
		//db, err := database.NewMongoDB(cfg)
		//if err != nil {
		//	log.Fatalf("Failed to connect to MongoDB: %v", err)
		//}

		server, err := main.InitializeServer()
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
