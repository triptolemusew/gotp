package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/triptolemusew/gotp/encryption"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new key",
	Run:   addCmdExecution,
}

type Token struct {
	name                 string
	key                  string
	password             string
	passwordConfirmation string
}

func (t *Token) IsPasswordSame() bool {
	return t.password == t.passwordConfirmation
}

func addCmdExecution(cmd *cobra.Command, args []string) {
	var token Token

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	token.name = prompt("Token name: ")
	token.key = prompt("Token key: ")

	fmt.Printf("\nAn empty password will not lock the file.\n")

	token.password, err = promptSecure("Password: ")
	if err != nil {
		log.Fatal(err)
	}

	filePath := fmt.Sprintf("%s/%s/%s", homeDir, TOKENFILES_DIR, token.name)

	fmt.Println() // Might not need this

	if token.password == "" {
		fmt.Println("Empty password. Will resume creating the secret unencrypted.")

		err := ioutil.WriteFile(filePath, []byte(token.key), 0644)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		token.passwordConfirmation, err = promptSecure("Confirm password: ")
		if err != nil {
			fmt.Errorf("Mismatched password")
			os.Exit(1)
		}

		if !token.IsPasswordSame() {
			fmt.Errorf("Mismatched password")
			os.Exit(1)
		}

		filePath += ENCRYPTED_EXT

		encryptedToken, err := encryption.Encrypt([]byte(token.password), []byte(token.key))
		if err != nil {
			log.Fatal(err)
		}

		err = ioutil.WriteFile(filePath, []byte(encryptedToken), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}
