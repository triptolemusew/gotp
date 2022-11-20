package cmd

import (
	// "fmt"
	// "os"

	"github.com/spf13/cobra"
)

var encCmd = &cobra.Command{
	Use:   "enc",
	Short: "Encrypt/Decrypt available keys",
	Run:   encCmdExecution,
}

func init() {
	encCmd.Flags().BoolP("encrypt", "e", true, "Encrypt the key")
	encCmd.Flags().BoolP("decrypt", "d", false, "Decrypt the key")
}

func encCmdExecution(cmd *cobra.Command, args []string) {
	// // encrypt, _ := cmd.Flags().GetString("encrypt")
	// // decrypt, _ := cmd.Flags().GetString("decrypt")
	// // encrypt, _ := cmd.Flags().GetBool("encrypt")
	// // decrypt, _ := cmd.Flags().GetBool("decrypt")
	// if len(args) > 1 {
	// 	fmt.Errorf("Only accepting single argument")
	// 	os.Exit(1)
	// }
	//
	// tokenName := args[0]
	// decrypt, _ := cmd.Flags().GetBool("decrypt")
	//
	// tokenfile, err := getTokenFile(tokenName)
	// if err != nil {
	// 	fmt.Errorf("Could not find the file: %s", tokenName)
	// 	os.Exit(1)
	// }
	//
	// if decrypt {
	//
	// }
}
