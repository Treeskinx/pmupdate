package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
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

// NOTE: Date Functions
func StringToDateF(dateInput string) time.Time {
	layout := "2006-01-02"
	regex := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)

	// Extract Dates
	sdate := regex.FindAllString(dateInput, 1)

	// Convert to date
	var date time.Time
	for _, s := range sdate {
		date, _ = time.Parse(layout, s)
		break
	}

	return date
}

func copyToClipboard(text string) error {
	// My linux copy command
	// cmd := exec.Command("wl-copy")

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		cmd = exec.Command("wl-copy")
	case "windows":
		cmd = exec.Command("clip")
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

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

func (a *App) PMDrop(data CSVData) string {

	file, err := os.Open(data.Path)
	if err != nil {
		return "Couldn't read the file"
	}
	reader := csv.NewReader(file)

	jobs, readErr := ReadCSV(reader)
	if readErr != "" {
		return "Incorrect PM List used"
	}
	error := copyToClipboard(jobs)
	if error != nil {
		fmt.Println("Error copying to clipboard:", error)
		return "Couldn't copy to clipboard"
	}

	fmt.Println("Text Copied to Clipboard!")
	return fmt.Sprintf("Text Copied to Clipboard!")
}

// This be the function to read the mcjiggy
func (a *App) PMList(data FileInput) string {
	reader := csv.NewReader(strings.NewReader(string(data.Data)))

	jobs, readErr := ReadCSV(reader)
	if readErr != "" {
		return "Incorrect PM List used"
	}

	error := copyToClipboard(jobs)
	if error != nil {
		fmt.Println("Error copying to clipboard:", error)
		return "Couldn't copy to clipboard"
	}

	fmt.Println("Text Copied to Clipboard!")
	return fmt.Sprintf("Text Copied to Clipboard!")
}

// This is going to be the PM List Creation function
func (a *App) PMCreate(data CSVDatas) string {
	var getbettername []Someshit

	stringClosedWO := "Closed WOs"
	stringPMWO := "All PM WOs"

	boolClosedWO := false
	boolPMWO := false

	// Determine which file is which
	for _, path := range data.Path {
		// Validate it is a CSV File
		matched, _ := regexp.MatchString(`^.*\.csv$`, path)
		if !matched {
			return fmt.Sprintf("%s is not a csv file, please select a csv file.", path)
		}
		file, err := os.Open(path)
		if err != nil {
			return "Couldn't read the file"
		}
		reader := csv.NewReader(file)

		for {
			line, error := reader.Read()
			if error == io.EOF {
				break
			} else if error != nil {
				fmt.Println(error)
			}
			if line[0] == "Job Category Short Name" && line[1] == "Job #" && line[7] == "Reported TGL" && line[12] == "Actual TID" {
				new := Someshit{Path: path, Name: stringClosedWO}
				getbettername = append(getbettername, new)
			} else if line[0] == "Job #" && line[4] == "Model" && line[11] == "Reported TGL" && line[12] == "Reported TID" {
				new := Someshit{Path: path, Name: stringPMWO}
				getbettername = append(getbettername, new)
			}
			break
		}
		defer file.Close()
	}

	// Validate that both file paths are present and set file pointers
	var fileClosedWO, fileAllPM Someshit
	for _, f := range getbettername {
		if f.Name == stringClosedWO {
			boolClosedWO = true
			fileClosedWO = f
		} else if f.Name == stringPMWO {
			boolPMWO = true
			fileAllPM = f
		}
	}

	if !boolPMWO && !boolClosedWO {
		return "Both files selected aren't formatted correctly"
	} else if !boolPMWO {
		return "All PM WOs file isn't formatted Correctly"
	} else if !boolClosedWO {
		return "Closed WOs file isn't formatted Correctly"
	}

	// NOTE: READ IN AND PROCESS CSV DATA
	// ==================================

	// Read in all Closed WOs
	initClosedWO, _ := ReadAllCSV(fileClosedWO)

	// Read in all PM WOs
	initPMWO, _ := ReadAllCSV(fileAllPM)

	// Filter out devices not GAC / EGK / CVM
	devicesClosedWO := RemoveDevices(initClosedWO)
	devicesPMWO := RemoveDevices(initPMWO)

	// Find the Unique TID and TGL from the Closed WO Slice
	closedTIDs, closedTGLs := returnClosedTidTgl(devicesClosedWO)

	// Filter out wo not in the same TGL as closedTGLs
	tglPMs := filterPMbyTGL(devicesPMWO, closedTGLs)

	// Separate open and closed pms
	openPM, closedPM := separateOpClPM(tglPMs)

	// Filter unique TID from the TGL and Open PM Slice
	pmTIDs, _ := returnClosedTidTgl(tglPMs)
	openTIDs, _ := returnClosedTidTgl(openPM)

	for _, pm := range openTIDs {

		if pm == "91600972" {
			fmt.Println("=========================================================")
			fmt.Println(pm)
			fmt.Println("=========================================================")
		}
	}
	// Create unique list of closed PMs based on pmTIDs and leave latest closed PM
	lastClosedPM := filterLastClosedPM(closedPM, pmTIDs, openPM)

	// Create unique list of closed WOs based on closedTIDs and leave latest closed wo
	lastClosedWO := filterLastClosedPM(devicesClosedWO, closedTIDs, openPM)

	// Filter Closed based on last closedPM
	tids := filterClosed(lastClosedWO, lastClosedPM)

	var pmList []string
	pmMap := make(map[string]AllJobsStruct)

	for _, pm := range openPM {
		pmMap[pm.TID] = pm
	}

	for _, tid := range tids {
		value, exists := pmMap[tid]
		if exists {
			pmList = append(pmList, value.Job)
		} else if !exists {
			fmt.Println(tid)
		}
	}

	var lestring string
	for _, id := range pmList {
		lestring = lestring + " " + id
	}
	copyToClipboard(lestring)

	// NOTE: Printing Statements
	fmt.Println("=========================================================")
	fmt.Println("### Step 1 ###")
	fmt.Println("Read File Closed WO Slice length: " + strconv.Itoa(len(initClosedWO)))
	fmt.Println("Read File All PM WO Slice length: " + strconv.Itoa(len(initPMWO)))
	fmt.Println("")
	fmt.Println("### Step 2 ###")
	fmt.Println("Removed Devices Closed WO Slice length: " + strconv.Itoa(len(devicesClosedWO)))
	fmt.Println("Removed Devices PM WO Slice length: " + strconv.Itoa(len(devicesPMWO)))
	fmt.Println("")
	fmt.Println("### Step 3 ###")
	fmt.Println("Unique Closed TID Slice length: " + strconv.Itoa(len(closedTIDs)))
	fmt.Println("Unique Closed TGL Slice length: " + strconv.Itoa(len(closedTGLs)))
	fmt.Println("")
	fmt.Println("### Step 4 ###")
	fmt.Println("Removed PMs via TGL Slice length: " + strconv.Itoa(len(tglPMs)))
	fmt.Println("")
	fmt.Println("### Step 5 ###")
	fmt.Println("Unique PM TID Slice length: " + strconv.Itoa(len(pmTIDs)))
	fmt.Println("")
	fmt.Println("### Step 6 ###")
	fmt.Println("Open PM Slice length: " + strconv.Itoa(len(openPM)))
	fmt.Println("Closed PM Slice length: " + strconv.Itoa(len(closedPM)))
	fmt.Println("")
	fmt.Println("### Step 7 ###")
	fmt.Println("Last Closed PM Slice length: " + strconv.Itoa(len(lastClosedPM)))
	fmt.Println("")
	fmt.Println("### Step 8 ###")
	fmt.Println("Last Closed WO Slice length: " + strconv.Itoa(len(lastClosedWO)))
	fmt.Println("")
	fmt.Println("### Step 9 ###")
	fmt.Println("TIDs for Open PM Slice length: " + strconv.Itoa(len(openTIDs)))
	fmt.Println("")
	fmt.Println("### Step 10 ###")
	fmt.Println("TIDs Eligible For PM Slice length: " + strconv.Itoa(len(tids)))
	fmt.Println("")
	fmt.Println("=========================================================")

	/*
		// NOTE: TROUBLESHOOOTING
		ts := 3
		fmt.Println(strconv.Itoa(ts))
		troubleshoot(data)

	*/
	return "All PMs copied to clipboard!"
}

// ==================================================
// NOTE: PM Update Functions

// function to read csv data and validate for correct file
func ReadCSV(reader *csv.Reader) (string, string) {

	exceptions := []string{"Closed", "Cancelled"}
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
			// Read the first line with the CSV data and validate it meets what is needed or
			// return an error to the frontend.

			if line[1] != "Job #" {
				return "", "Validation failed."
			}
			if line[3] != "Priority Short Name" {
				return "", "Validation failed."
			}
			if line[5] != "Status" {
				return "", "Validation failed."
			}

			count = 1
		}

	}
	return jobs, ""
}

// ==================================================
// NOTE: PM Create Functions

// Function to return 2 slices for the closedWO file
// Slice 1 will contain model, tgl and tid where tid is not duplicated
// Slice 2 will contain tgl where tgl isn't duplicated

// Read all CSV Data to relevant Slice based on File name
func ReadAllCSV(csvFile Someshit) ([]AllJobsStruct, error) {
	stringClosedWO := "Closed WOs"
	stringPMWO := "All PM WOs"

	file, err := os.Open(csvFile.Path)
	if err != nil {
		return nil, nil
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading csv: %v", err)
	}

	var records []AllJobsStruct
	if csvFile.Name == stringPMWO {
		for i, row := range lines {
			if i == 0 {
				continue // skip header or malformed rows
			}

			if row[6] != "Cancelled" {
				dateStr := strings.TrimSpace(row[13])
				parsedDate, _ := time.Parse("2006-01-02 03:04 PM", dateStr) // adjust layout if needed

				record := AllJobsStruct{
					Job:    strings.TrimSpace(row[0]),
					Model:  strings.TrimSpace(row[4]),
					TGL:    strings.TrimSpace(row[11]),
					TID:    strings.TrimSpace(row[12]),
					Date:   parsedDate,
					Status: strings.TrimSpace(row[6]),
				}
				records = append(records, record)
			}
		}
	} else if csvFile.Name == stringClosedWO {
		for i, row := range lines {
			if i == 0 || len(row) < 5 {
				continue // skip header or malformed rows
			}

			dateStr := strings.TrimSpace(row[8])
			parsedDate, err := time.Parse("2006-01-02 03:04 PM", dateStr) // adjust layout if needed
			if err != nil {
				return nil, fmt.Errorf("invalid date format in row %d: %v", i+1, err)
			}

			record := AllJobsStruct{
				Job:   strings.TrimSpace(row[1]),
				Model: strings.TrimSpace(row[13]),
				TGL:   strings.TrimSpace(row[7]),
				TID:   strings.TrimSpace(row[12]),
				Date:  parsedDate,
			}
			records = append(records, record)
		}
	}

	return records, nil

}

// Filter out devices not GAC / EGK / CVM
func RemoveDevices(jobs []AllJobsStruct) []AllJobsStruct {
	var filtered []AllJobsStruct
	devices := []string{"GAC", "AVM/CVM-s", "EGK", "FPD-s", "FPD-g", "SEM"}
	for _, job := range jobs {
		if slices.Contains(devices, job.Model) {
			filtered = append(filtered, job)
		}
	}

	return filtered
}

// Find the Unique TID and TGL from the Closed WO Slice
func returnClosedTidTgl(closedWOs []AllJobsStruct) ([]string, []string) {
	var TIDs, TGLs []string

	for _, wo := range closedWOs {
		if !slices.Contains(TIDs, wo.TID) {
			TIDs = append(TIDs, wo.TID)
		}
		if !slices.Contains(TGLs, wo.TGL) {
			TGLs = append(TGLs, wo.TGL)
		}
	}

	return TIDs, TGLs
}

// Filter out wo not in the same TGL as closedTGLs
func filterPMbyTGL(allPMs []AllJobsStruct, sliceTGL []string) []AllJobsStruct {
	var filtered []AllJobsStruct

	for _, pm := range allPMs {
		if slices.Contains(sliceTGL, pm.TGL) {
			filtered = append(filtered, pm)
		}
	}

	return filtered
}

// Separate open and closed pms
func separateOpClPM(pms []AllJobsStruct) ([]AllJobsStruct, []AllJobsStruct) {
	var closed, open []AllJobsStruct

	for _, pm := range pms {
		if pm.Status == "Closed" {
			closed = append(closed, pm)
		} else if pm.Status != "Closed" {
			open = append(open, pm)
		}
	}

	return open, closed
}

// Create unique list of closed PMs based on pmTIDs and leave latest closed PM
func filterLastClosedPM(
	closedPM []AllJobsStruct,
	tids []string,
	openPM []AllJobsStruct) []AllJobsStruct {
	var filtered []AllJobsStruct
	var noPMTID []string

	for _, tid := range tids {
		returned := findValue(closedPM, tid)
		if len(returned) == 0 {
			noPMTID = append(noPMTID, tid)
			continue // or handle this case differently
		}
		// sort dates so latest is first
		if len(returned) > 1 {
			sort.Slice(returned, func(i, j int) bool {
				return returned[i].Date.After(returned[j].Date)
			})
		}
		filtered = append(filtered, returned[0])
	}

	// Iterate through the openPM list and add PMs to filtered if TID in noPMTID
	// Only do this if noPMTID > 0
	if len(noPMTID) > 0 {
		for _, pm := range openPM {
			if slices.Contains(noPMTID, pm.TID) {
				dateS := "2000/01/01 00:00 AM"
				date, _ := time.Parse("2006/01/02 03:04 PM", dateS)
				beep := AllJobsStruct{
					TID:    pm.TID,
					Job:    pm.Job,
					Model:  pm.Model,
					TGL:    pm.TGL,
					Status: pm.Status,
					Date:   date,
				}
				filtered = append(filtered, beep)
			}
		}
	}

	return filtered
}

func findValue(closed []AllJobsStruct, tid string) []AllJobsStruct {
	var found []AllJobsStruct

	for _, pm := range closed {
		if pm.TID == tid {
			found = append(found, pm)
		}
	}
	return found
}

// Filter Closed WOs based on last closedPM
func filterClosed(closedWO []AllJobsStruct, lastClosedPM []AllJobsStruct) []string {
	var filtered []string

	// turn the closed pm slice into a map
	cwMap := make(map[string]AllJobsStruct)

	for _, cwo := range closedWO {
		cwMap[cwo.TID] = cwo
	}

	for _, pm := range lastClosedPM {
		if pm.Model == "AVM/CVM-s" {
			dateCheck := pm.Date.AddDate(0, 0, 30)
			if dateCheck.Before(cwMap[pm.TID].Date) {
				filtered = append(filtered, pm.TID)
			}
			continue
		} else {
			dateCheck := pm.Date.AddDate(0, 0, 90)
			if dateCheck.Before(cwMap[pm.TID].Date) {
				filtered = append(filtered, pm.TID)
			}
		}
	}

	return filtered
}

// ==================================================
// NOTE: TROUBLESHOOOTING FUNCTION

// Use this to read in 2 csv files into different slices and then compare the pair to find
// out why the output is different to the current list
func troubleshoot(data CSVDatas) string {
	if len(data.Path) != 2 {
		return "Please provide exactly two CSV file paths"
	}

	readCSV := func(path string) ([]string, error) {
		file, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %v", err)
		}
		defer file.Close()

		reader := csv.NewReader(file)
		lines, err := reader.ReadAll()
		if err != nil {
			return nil, fmt.Errorf("error reading csv: %v", err)
		}

		var ids []string
		for i, row := range lines {
			if i == 0 || len(row) <= 1 {
				continue // Skip header or malformed rows
			}
			ids = append(ids, strings.TrimSpace(row[1]))
		}
		return ids, nil
	}

	ids1, err := readCSV(data.Path[0])
	if err != nil {
		return err.Error()
	}
	ids2, err := readCSV(data.Path[1])
	if err != nil {
		return err.Error()
	}

	// Create sets
	set1 := make(map[string]bool)
	set2 := make(map[string]bool)
	for _, id := range ids1 {
		set1[id] = true
	}
	for _, id := range ids2 {
		set2[id] = true
	}

	var onlyInFirst, onlyInSecond []string
	for id := range set1 {
		if !set2[id] {
			onlyInFirst = append(onlyInFirst, id)
		}
	}
	for id := range set2 {
		if !set1[id] {
			onlyInSecond = append(onlyInSecond, id)
		}
	}

	// Print results
	fmt.Println("Only in first CSV:")
	for _, id := range onlyInFirst {
		fmt.Println(id)
	}

	fmt.Println("\nOnly in second CSV:")
	for _, id := range onlyInSecond {
		fmt.Println(id)
	}

	// Optionally copy new items to clipboard
	var lestring string
	for _, id := range onlyInSecond {
		lestring += " " + id
	}
	copyToClipboard(lestring)

	return "Done"
}
