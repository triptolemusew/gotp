package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/triptolemusew/gotp/otp"
	"github.com/triptolemusew/gotp/tui"

	"golang.design/x/clipboard"
)

var pathFlag = flag.String("path", ".gotp", "Path to TOTP keys")

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: gotp [options]\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	defer ui.Close()

	tuiManager := tui.NewManager()

	keys, err := otp.GetAllKeys(*pathFlag)
	if err != nil {
		log.Fatalf("failed to get all totp keys: %v", err)
	}
	tuiManager.InitializeWidgets(keys)

	if err := startApp(tuiManager); err != nil {
		log.Fatalf("failed to start the app: %v", err)
	}
}

func startApp(tuiManager *tui.Manager) error {
	// Init clipboard
	clipboard.Init()

	uiRefreshTicker := time.NewTicker(1 * time.Millisecond)
	defer uiRefreshTicker.Stop()

	width, height := ui.TerminalDimensions()
	uiEvents := ui.PollEvents()

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return nil
			case "<Backspace>":
				tuiManager.RemoveLastCharBuffer()
			case "<Down>":
				tuiManager.ScrollDown()
			case "<Up>":
				tuiManager.ScrollUp()
			case "<Space>":
				tuiManager.HandleBuffer(" ")
			case "<Enter>":
				if err := tuiManager.SelectRow(); err != nil {
					return err
				}
				return nil
			default:
				tuiManager.HandleBuffer(e.ID)
			}
		case <-uiRefreshTicker.C:
			tuiManager.UpdateWidgets(width, height)
		}
	}
}
