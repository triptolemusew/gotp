package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var encCmd = &cobra.Command{
	Use:   "enc",
	Short: "Encrypt/Decrypt available keys",
	Run:   enc,
}

func enc(cmd *cobra.Command, args []string) {
	fmt.Println("Encrypt/Decrypt")
	fmt.Println(args)
}
