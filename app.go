package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os/exec"
	"slices"
	"strings"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

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

type FileInput struct {
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	filetype string `json:"type"`
	Data     string `json:"data"`
}

type CSVData struct {
	Data []byte `json:"data"`
}

func (a *App) ProcessFiles(files []string) {
	for _, file := range files {
		fmt.Println("File:", file)
	}
}

// This be the function to read the mcjiggy
func (a *App) PMList(data FileInput) string {

	//	fileLoc := "/home/treeskin/Downloads/jobs-2025-06-06.csv"
	exceptions := []string{"Closed", "Cancelled"}

	//	csvFile, err := os.Open(fileLoc)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	reader := csv.NewReader(strings.NewReader(string(data.Data)))

	var job []string
	var jobs string
	count := 0

	for {
		line, error := reader.Read()
		fmt.Println("Am I doing the loop?")
		fmt.Println("=================")
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
			fmt.Println(line[5])
		}
		if count == 0 {
			count = 1
		}

		fmt.Println("=================")
	}

	error := copyToClipboard(jobs)
	if error != nil {
		fmt.Println("Error copying to clipboard:", error)
		return "Couldn't copy to clipboard"
	}

	fmt.Println("Text Copied to Clipboard!")
	return fmt.Sprintf("It has been done.")
}
