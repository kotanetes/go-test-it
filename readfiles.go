package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

func readFilesFromDirecotry(path string) {

	if path == "" {
		path = "."
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Name())
	}
}
