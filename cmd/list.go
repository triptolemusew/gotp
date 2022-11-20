package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available keys",
	Run:   listCmdExecution,
}

func listCmdExecution(cmd *cobra.Command, args []string) {
	tokens, err := listTokens()
	if err != nil {
		log.Fatal(err)
	}

	for _, token := range tokens {
		fmt.Println(token)
	}
}
