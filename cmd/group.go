package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Group subcommands",
	Run:   group,
}

func group(cmd *cobra.Command, args []string) {
	fmt.Println("Grouping")
	fmt.Println(args)
}
