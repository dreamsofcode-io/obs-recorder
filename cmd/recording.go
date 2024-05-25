/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// recordingCmd represents the recording command
var recordingCmd = &cobra.Command{
	Use:   "recording",
	Short: "The parent command for all the recording subcommands",
}

func init() {
	clientCmd.AddCommand(recordingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// recordingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// recordingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
