package main

import (
	"embed"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"slices"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

// simple
func copyToClipboard(text string) error {
	cmd := exec.Command("wl-copy")
	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	if _, err := in.Write([]byte(text)); err != nil {
		return err
	}
	if err := in.Close(); err != nil {
		return err
	}
	return cmd.Wait()
}

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "pmupdate",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}

	fileLoc := "/home/treeskin/Downloads/jobs-2025-06-06.csv"
	exceptions := []string{"Closed", "Cancelled"}

	csvFile, err := os.Open(fileLoc)
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(csvFile)

	var job []string
	var jobs string
	count := 0

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			fmt.Println(error)
		}
		if count == 1 {
			if !slices.Contains(exceptions, string(line[5])) {
				job = append(job, string(line[1]))
				jobs += string(line[1])
				jobs += " "
			}
		}
		if count == 0 {
			count = 1
		}
	}

	error := copyToClipboard(jobs)
	if error != nil {
		fmt.Println("Error copying to clipboard:", err)
		return
	}

	fmt.Println("Text Copied to Clipboard!")

}
