package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/spf13/cobra"
	"github.com/triptolemusew/gotp/encryption"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show key",
	Run:   showCmdExecution,
}

func showCmdExecution(cmd *cobra.Command, args []string) {
	var pattern, secret string

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	if len(args) > 1 {
		fmt.Println("Not supporting mutliple keys at once")
		os.Exit(1)
	}

	tokenName := args[0]

	pattern = fmt.Sprintf("%s/%s/%s*", homeDir, TOKENFILES_DIR, tokenName)

	matches, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(matches)

	for _, match := range matches {
		file, _ := os.Stat(match)
		ext := filepath.Ext(match)

		if fileNameWithoutExtension(file.Name()) == tokenName {
			token, err := ioutil.ReadFile(match)
			if err != nil {
				log.Fatal(err)
			}

			var password string
			fmt.Printf("Enter password: ")
			fmt.Scan(&password)

			if ext == ".enc" {
				token, err := encryption.Decrypt([]byte(password), token)
				if err != nil {
					log.Fatal(err)
				}

				secret, err := totp.GenerateCode(string(token), time.Time{})
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println(string(secret))
			}
		}
	}

	fmt.Println(string(secret))
}

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
