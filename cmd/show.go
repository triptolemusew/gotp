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
	"github.com/thoas/go-funk"
	"github.com/triptolemusew/gotp/encryption"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show key",
	Run:   showCmdExecution,
}

type PlainToken struct {
	secret string
	name   string
}

func (p *PlainToken) Print() {
	file, _ := os.Stat(p.name)
	fmt.Printf("%s - %s\n", file.Name(), p.secret)
}

func init() {
	showCmd.Flags().BoolP("all", "a", false, "Get all of the tokens with secret")
}

func showCmdExecution(cmd *cobra.Command, args []string) {
	var pattern string
	var tokenName string
	var plainTokenList []PlainToken

	showAll, _ := cmd.Flags().GetBool("all")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	if len(args) > 1 {
		fmt.Println("Not supporting mutliple keys at once")
		os.Exit(1)
	}

	if showAll {
		pattern = fmt.Sprintf("%s/%s/*", homeDir, TOKENFILES_DIR)
	} else {
		tokenName = args[0]
		pattern = fmt.Sprintf("%s/%s/%s*", homeDir, TOKENFILES_DIR, tokenName)
	}

	matches, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatal(err)
	}

	// Match depended
	if !showAll {
		matches = funk.FilterString(matches, func(x string) bool {
			file, _ := os.Stat(x)
			return fileNameWithoutExtension(file.Name()) == tokenName
		})
	}

	for _, match := range matches {
		ext := filepath.Ext(match)

		token, err := ioutil.ReadFile(match)
		if err != nil {
			log.Fatal(err)
		}

		if ext == ".enc" {
			password, _ := promptSecure("Enter password: ")

			token, err := encryption.Decrypt([]byte(password), token)
			if err != nil {
				log.Fatal(err)
			}

			secret, err := totp.GenerateCode(string(token), time.Time{})
			if err != nil {
				log.Fatal(err)
			}
			plainTokenList = append(plainTokenList, PlainToken{
				secret: secret,
				name:   match,
			})
		} else {
			secret, err := totp.GenerateCode(string(token), time.Time{})
			if err != nil {
				log.Fatal(err)
			}
			plainTokenList = append(plainTokenList, PlainToken{
				secret: secret,
				name:   match,
			})
		}
	}

	// Print all of the tokenList
	for _, token := range plainTokenList {
		token.Print()
	}
}

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
