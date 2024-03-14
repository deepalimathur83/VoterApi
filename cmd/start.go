/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"drexel.edu/voter-api/pkg/adding"
	"drexel.edu/voter-api/pkg/changing"
	"drexel.edu/voter-api/pkg/http/rest"
	"drexel.edu/voter-api/pkg/listing"
	"drexel.edu/voter-api/pkg/storage/json"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts The Server",
	Long:  `The start command is used to start the server.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start called")

		storage, err := json.NewVoterDB("../data")
		if err != nil {
			log.Panic("The server couldn't start")
			panic(err)
		}
		add := adding.NewService(storage)
		list := listing.NewService(storage)
		change := changing.NewService(storage)

		rest.Handler(add, list, change)
		fmt.Println("The server is live now: http://localhost:3000")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
