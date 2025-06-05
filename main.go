package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"slices"
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

func main() {
	fileLoc := "/home/treeskin/Downloads/jobs-2025-06-05 (3).csv"
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
