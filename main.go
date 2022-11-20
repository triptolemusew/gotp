package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/triptolemusew/gotp/encryption"
	"github.com/triptolemusew/gotp/cmd"
)

const (
	TOKENFILES_DIR = ".gotp/tokens"
	ENCRYPTED_EXT  = ".enc"
)

func listTokens(tokenFilesDir string) error {
	files, err := os.ReadDir(tokenFilesDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		code, _ := showToken(fmt.Sprintf("%s/%s", tokenFilesDir, file.Name()))
		fmt.Printf("%s -> %s", file.Name(), code)
		fmt.Println()
	}

	return nil
}

func showToken(tokenPath string) (string, error) {
	var pattern string
	var secret string

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}

	for _, match := range matches {
		file, _ := os.Stat(match)
		ext := filepath.Ext(match)

		if fileNameWithoutExtension(file.Name()) == "a" {
			token, err := ioutil.ReadFile(match)
			if err != nil {
				return "", err
			}

			var password string
			fmt.Printf("Enter password: ")
			fmt.Scan(&password)

			if ext == ENCRYPTED_EXT {
				token, err := encryption.Decrypt([]byte(password), token)
				if err != nil {
					return "", err
				}

				secret, err := totp.GenerateCode(string(token), time.Time{})
				if err != nil {
					return "", err
				}

				return string(secret), nil
			}
		}
	}

	return string(secret), nil
}

type Token struct {
	name     string
	key      string
	password string
}

func addToken(homeDir string) error {
	var token Token
	var tokenPasswordConfirmation string

	fmt.Printf("Token name: ")
	fmt.Scan(&token.name)
	fmt.Printf("Token key: ")
	fmt.Scan(&token.key)

	fmt.Println()
	fmt.Println("An empty password will not lock the file")

	fmt.Printf("Password: ")
	fmt.Scan(&token.password)

	filePath := fmt.Sprintf("%s/%s/%s", homeDir, TOKENFILES_DIR, token.name)

	if token.password == "" {
		fmt.Println("Empty password. Will resume creating the secret unencrypted")

		err := ioutil.WriteFile(filePath, []byte(token.key), 0644)
		if err != nil {
			return err
		}
	} else {
		fmt.Println()

		fmt.Printf("Confirm password: ")
		fmt.Scan(&tokenPasswordConfirmation)

		if token.password != tokenPasswordConfirmation {
			return &GotpError{}
		}

		filePath += ".enc"

		encryptedToken, err := encryption.Encrypt([]byte(token.password), []byte(token.key))
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(filePath, []byte(encryptedToken), 0644)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Created the token at [%s] \n", filePath)

	return nil
}

func main() {
	cmd.Execute()

	// tokenFilesDir := os.Getenv("TOKENFILES_DIR")
	// homeDir, _ := os.UserHomeDir()
	//
	// if tokenFilesDir == "" {
	// 	tokenFilesDir = fmt.Sprintf("%s/%s", homeDir, TOKENFILES_DIR)
	// }

	// fmt.Println(tokenFilesDir)
	// listTokens(tokenFilesDir)
	// fmt.Println(code)

	// Adding routine
	// fmt.Println("Adding subroutine")
	// addToken(homeDir)

	// Decrypt routine
	// password := decryptToken(homeDir, "a.enc")
	// fmt.Println("RESULT")
	// fmt.Println(password)
}
