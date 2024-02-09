package cmd

import (
	"image"
	_ "image/png"
	"os"
	"path/filepath"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/triptolemusew/gotp/db"
	"github.com/triptolemusew/gotp/otp"

	"github.com/urfave/cli/v2"
)

func SyncCommandExecution(ctx *cli.Context) error {
	homeDirectory, _ := os.UserHomeDir()
	appDirectory := ".gotp"
	directory := filepath.Join(homeDirectory, appDirectory, "qr")

	qrReader := qrcode.NewQRCodeReader()

	dbClient, err := db.GetClient(appDirectory)
	if err != nil {
		return err
	}

	err = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		img, _, _ := image.Decode(file)
		bmp, _ := gozxing.NewBinaryBitmapFromImage(img)
		decodedQr, _ := qrReader.Decode(bmp, nil)

		key, err := otp.ParseURL(decodedQr.String())

		db.Create(dbClient, key)

		return nil
	})

	return err
}
