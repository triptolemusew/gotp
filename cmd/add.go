package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new key",
	Run:   add,
}

func add(cmd *cobra.Command, args []string) {
	fmt.Println("Adding")
	fmt.Println(args)
}
