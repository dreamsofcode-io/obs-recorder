/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/dreamsofcode-io/obs-remote/internal/server"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Server based top level command",
	Run: func(cmd *cobra.Command, args []string) {
		if err := server.Start(cmd.Context()); err != nil {
			log.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
