package main

import (
	"time"
)

// ==================================================
// NOTE: PM Update Structs

// Reading CSV Paths
type CSVDatas struct {
	Path []string `json:"Dropped"`
}

// ==================================================
// NOTE: PM Update Structs
type Someshit struct {
	Path string
	Name string
}

type FileInput struct {
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Filetype string `json:"type"`
	Data     string `json:"data"`
}

type CSVData struct {
	Path string `json:"Dropped"`
}

type AllJobsStruct struct {
	Job    string
	Model  string
	TGL    string
	TID    string
	Status string
	Date   time.Time
}

type PM struct {
	PMWO   string
	TID    string
	Model  string
	LastCM time.Time
	LastPM time.Time
	PMV    time.Time
	PMReq  string
}
