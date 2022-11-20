package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove key",
	Run:   remove,
}

func remove(cmd *cobra.Command, args []string) {
	fmt.Println("Remove")
	fmt.Println(args)
}
