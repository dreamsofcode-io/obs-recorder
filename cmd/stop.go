/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/ncruces/zenity"
	"github.com/spf13/cobra"

	"github.com/dreamsofcode-io/obs-remote/internal/recording"
	"github.com/dreamsofcode-io/obs-remote/internal/server"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops any active recording and prompts to keep/rename the file",
	RunE: func(cmd *cobra.Command, args []string) error {
		baseURL, err := cmd.Flags().GetString("url")
		if err != nil {
			log.Println("failed to get the baseURL", err)
			return err
		}

		uri := fmt.Sprintf("%s/record/stop", baseURL)
		req, err := http.NewRequestWithContext(cmd.Context(), http.MethodPost, uri, nil)
		if err != nil {
			log.Println("failed to send stop request", err)
			return err
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("failed to do request", err)
			return err
		}

		if res.StatusCode >= 400 {
			log.Println("request failed", res.StatusCode)
			return errors.New("request failed")
		}

		rec := recording.Recording{}

		if err := json.NewDecoder(res.Body).Decode(&rec); err != nil {
			log.Println("failed to decode response body", err)
			return err
		}

		title, err := zenity.Entry("Enter recording title:", zenity.Title("Name recording"))
		if errors.Is(err, zenity.ErrCanceled) {
			// delete the file
		}

		url := fmt.Sprintf("%s/recordings/%s", baseURL, rec.ID)
		body := server.UpdateBody{
			Filename: title,
		}

		bs := &bytes.Buffer{}

		if err = json.NewEncoder(bs).Encode(body); err != nil {
			log.Println("failed to encode the filename request body", err)
			return err
		}

		req, err = http.NewRequestWithContext(cmd.Context(), http.MethodPut, url, bs)
		if err != nil {
			log.Println("failed to create new request", err)
			return err
		}

		res, err = http.DefaultClient.Do(req)
		if err != nil {
			log.Println("failed to send request", err)
		}

		return nil
	},
}

func init() {
	recordingCmd.AddCommand(stopCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stopCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stopCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
