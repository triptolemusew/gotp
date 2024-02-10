package main

import (
	"log"
	"os"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/triptolemusew/gotp/otp"
	"github.com/triptolemusew/gotp/tui"
	"github.com/urfave/cli/v2"
	"golang.design/x/clipboard"
)

const VERSION = "0.1"

func main() {
	app := &cli.App{
		Name:    "Gotp",
		Usage:   "Get your token from cli",
		Version: VERSION,
		Action: func(ctx *cli.Context) error {
			tuiManager := tui.NewManager()

			keys, err := otp.GetAllKeys(".gotp")
			if err != nil {
				return err
			}
			tuiManager.InitializeWidgets(keys)

			if err := startApp(tuiManager); err != nil {
				return err
			}

			return nil
		},
	}

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	defer ui.Close()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func startApp(tuiManager *tui.Manager) error {
	// Init clipboard
	clipboard.Init()

	uiRefreshTicker := time.NewTicker(1 * time.Millisecond)
	defer uiRefreshTicker.Stop()

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
			tuiManager.UpdateWidgets(1, 1)
		}
	}
}
