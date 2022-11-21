package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"

	tb "github.com/nsf/termbox-go"
	"github.com/pquerna/otp/totp"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
	"github.com/triptolemusew/gotp/encryption"
	"golang.design/x/clipboard"
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

func (p *PlainToken) GetLine() string {
	file, _ := os.Stat(p.name)
	return fmt.Sprintf("%s \t %s", file.Name(), p.secret)
}

func init() {
	showCmd.Flags().BoolP("all", "a", false, "Get all of the tokens with secret")
}

func renderUI(tokenList []PlainToken) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	tb.SetInputMode(tb.InputEsc)

	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	defer ui.Close()

	l := widgets.NewList()
	l.Title = "Available Accounts"

	for _, token := range tokenList {
		l.Rows = append(l.Rows, token.GetLine())
	}
	var copiedTokenList []PlainToken
	copy(copiedTokenList, tokenList)

	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = false
	l.SetRect(0, 0, 50, 14)

	p := widgets.NewParagraph()
	p.Title = "Search"
	p.Text = "> "
	p.SetRect(0, 14, 50, 17)
	p.TextStyle.Fg = ui.ColorWhite
	p.BorderStyle.Fg = ui.ColorCyan

	ui.Render(l, p)

	uiEvents := ui.PollEvents()

	var bufferText string

	for {
		e := <-uiEvents
		switch e.ID {
		case "<C-c>":
			return
		case "<Enter>":
			{
				t := copiedTokenList[l.SelectedRow]
				clipboard.Write(clipboard.FmtText, []byte(t.secret))

				// Check if current environment is running on Wayland
				if os.Getenv("WAYLAND_DISPLAY") != "" {
					exec.Command("wl-copy", t.secret).Run()
				}
				return
			}
		case "<Down>":
			l.ScrollDown()
		case "<Up>":
			l.ScrollUp()
		case "<C-d>":
			l.ScrollHalfPageDown()
		case "<C-u>":
			l.ScrollHalfPageUp()
		case "<Backspace>":
			{
				if last := len(bufferText) - 1; last >= 0 {
					bufferText = bufferText[:last]
				}
			}
		default:
			{
				bufferText += e.ID
			}
		}

		// TODO: Do it proper
		l.Rows, copiedTokenList = updateList(tokenList, bufferText)
		p.Text = "> " + bufferText

		ui.Render(l, p)
	}
}

func updateList(t []PlainToken, search string) ([]string, []PlainToken) {
	var output []string
	var outputOther []PlainToken

	items := funk.Filter(t, func(x PlainToken) bool {
		return strings.Contains(x.name, search)
	})
	if items, ok := items.([]PlainToken); ok {
		for _, token := range items {
			output = append(output, token.GetLine())
			outputOther = append(outputOther, token)
		}
	}
	return output, outputOther
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

			token, err = encryption.Decrypt([]byte(password), token)
			if err != nil {
				log.Fatal(err)
			}
		}

		secret, err := totp.GenerateCode(string(token), time.Now())
		if err != nil {
			log.Fatal(err)
		}

		plainTokenList = append(plainTokenList, PlainToken{
			secret: secret,
			name:   match,
		})
	}

	renderUI(plainTokenList)
}

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
