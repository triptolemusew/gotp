package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available keys",
	Run:   list,
}

func list(cmd *cobra.Command, args []string) {
	fmt.Println("List")
	fmt.Println(args)
}
