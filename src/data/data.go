package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Mega struct {
	User            string
	Password        string
	DownloadWorkers int
	UploadWorkers   int
	SkipSameSize    bool
	Verbose         int
}

type Sql struct {
	Host             string
	Port             uint16
	Database         string
	User             string
	DatabasePassword string
}
type Configuration struct {
	MEGA []Mega
	SQL  []Sql
}

var Config *Configuration

func Parse(file_location string) {
	fmt.Printf("Location of file: %s\n", file_location)
	data, err := ioutil.ReadFile(file_location)
	if err != nil {
		fmt.Println("File error:", err)
		return
	}
	err = json.Unmarshal(data, Config)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
}
