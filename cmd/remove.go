package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove key",
	Run:   removeCmdExecution,
}

func removeCmdExecution(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		os.Exit(1)
	}

	tokenName := args[0]
	token, err := getTokenFile(tokenName)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	if token == nil {
		log.Fatal("Could not find the file")
		os.Exit(1)
	}

	err = os.Remove(token.Path)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully deleted: %s", token.Path)
}
