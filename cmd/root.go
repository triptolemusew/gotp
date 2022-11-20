package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gotp",
	Version: "0.0.1",
	Short: "Gotp - OTP authentication in golang",
	// Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	// cobra.OnInitialize(initConfig)

	// Add all of the commands here
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Adding available commands
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(encCmd)
	rootCmd.AddCommand(groupCmd)
	rootCmd.AddCommand(removeCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
