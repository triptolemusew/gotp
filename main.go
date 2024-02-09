package main

import (
	"log"
	"os"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/triptolemusew/gotp/cmd"
	"github.com/triptolemusew/gotp/db"
	"github.com/triptolemusew/gotp/tui"
	"github.com/urfave/cli/v2"
)

const VERSION = "0.1"

func main() {
	app := &cli.App{
		Name:    "Gotp",
		Usage:   "Get your token from cli",
		Version: VERSION,
		Action: func(ctx *cli.Context) error {
			if err := ui.Init(); err != nil {
				log.Fatalf("failed to initialize termui: %v", err)
			}
			defer ui.Close()

			tuiManager := tui.NewManager()

			dbClient, _ := db.GetClient(".gotp")
			keys, err := db.GetAll(dbClient)
			if err != nil {
				return err
			}

			tuiManager.InitializeWidgets(keys)

			startApp(tuiManager)

			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "Add account",
				Action: func(ctx *cli.Context) error {
					return nil
				},
			},
			{
				Name:    "remove",
				Aliases: []string{"r"},
				Usage:   "Remove account",
				Action: func(ctx *cli.Context) error {
					return nil
				},
			},
			{
				Name:    "sync",
				Aliases: []string{"s"},
				Usage:   "Sync accounts with Gotp folder",
				Action: func(ctx *cli.Context) error {
					err := cmd.SyncCommandExecution(ctx)
					return err
				},
			},
			{
				Name:    "init",
				Aliases: []string{"i"},
				Usage:   "Initialize the application",
				Action: func(ctx *cli.Context) error {
					if err := db.Migrate(".gotp"); err != nil {
						return err
					}
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func startApp(tuiManager *tui.Manager) {
	uiRefreshTicker := time.NewTicker(1 * time.Millisecond)
	defer uiRefreshTicker.Stop()

	uiEvents := ui.PollEvents()

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<Backspace>":
				tuiManager.RemoveLastCharBuffer()
			default:
				tuiManager.HandleBuffer(e.ID)
			}
		case <-uiRefreshTicker.C:
			tuiManager.UpdateWidgets(1, 1)
		}
	}
}
