package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gosexy/exif"
)

const dateForm = "2006:01:02 15:04:05"

func updateDate(file string) {
	ext := filepath.Ext(file)
	if ext != ".JPG" && ext != ".jpg" {
		return
	}

	reader := exif.New()
	err := reader.Open(file)
	if err != nil {
		fmt.Printf("Error reading EXIF data of file: %s. Skipping...", err.Error())
		return
	}
	create_date, ok := reader.Tags["Date and Time (Original)"]

	if ok {
		newDate, err := time.Parse(dateForm, create_date)
		newDate = newDate.Add(time.Duration(3) * time.Hour)
		if err != nil {
			fmt.Printf("ERRROR: %s", err)
		} else {
			fmt.Printf("Updating file: %s --> %s\n", file, create_date)
			err = os.Chtimes(file, newDate, newDate)
		}

	}
}

func visit(path string, f os.FileInfo, err error) error {
	updateDate(path)
	return nil
}

func main() {
	if len(os.Args) == 1 {
		log.Fatal("Must supply directory or file as argument")
	}
	path := os.Args[1]
	file, err := os.Stat(path)
	if err == nil && file.IsDir() {
		filepath.Walk(path, visit)
	} else if err == nil && !file.IsDir() {
		updateDate(path)
	}
}
