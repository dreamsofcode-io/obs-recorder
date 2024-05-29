/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts a recording on the remote server.",
	Long: `Starts a recording on the remote server provided that there
	is no recording already`,
	RunE: func(cmd *cobra.Command, args []string) error {
		baseURL, err := cmd.Flags().GetString("url")
		if err != nil {
			log.Println("failed to get the url")
			return err
		}

		url := fmt.Sprintf("%s/record/start", baseURL)

		req, err := http.NewRequestWithContext(cmd.Context(), http.MethodPost, url, nil)
		if err != nil {
			log.Println("failed to send request", err)
			return err
		}

		client := http.Client{}
		if _, err = client.Do(req); err != nil {
			log.Println("failed to do request", err)
		}

		return nil
	},
}

func init() {
	recordingCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
