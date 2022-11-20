package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/term"
)

const (
	TOKENFILES_DIR = ".gotp/tokens"
	ENCRYPTED_EXT  = ".enc"
)

type TokenFile struct {
	Path    string
	content []byte
}

func prompt(message string) string {
	var temp string
	fmt.Printf(message)
	fmt.Scan(&temp)
	return temp
}

func promptSecure(message string) (string, error) {
	fmt.Printf(message)
	byteText, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", nil
	}
	return strings.TrimSpace(string(byteText)), nil
}

func listTokens() ([]string, error) {
	var tokens []string

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("%s/%s", homeDir, TOKENFILES_DIR)
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		tokens = append(tokens, fileNameWithoutExtension(file.Name()))
	}

	return tokens, nil
}

func getTokenFile(tokenName string) (*TokenFile, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	pattern := fmt.Sprintf("%s/%s/%s*", homeDir, TOKENFILES_DIR, tokenName)

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	for _, match := range matches {
		file, _ := os.Stat(match)

		if fileNameWithoutExtension(file.Name()) == tokenName {
			content, err := ioutil.ReadFile(match)
			if err != nil {
				return nil, err
			}

			return &TokenFile{
				content: content,
				Path:    match,
			}, nil
		}
	}

	return nil, nil
}

func checkTokenFilesDir(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func getFileBasename(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

func getFileExtension(filename string) string {
	return filepath.Ext(filename)
}
