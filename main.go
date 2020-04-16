package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/kotanetes/go-test-it/model"
	"github.com/kotanetes/go-test-it/service"
	"github.com/kotanetes/go-test-it/utils"
	"github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logrus.SetLevel(logrus.DebugLevel)

	// intializing remote service
	service.InitHTTPClient(&http.Client{})
}

var (
	filePath, fileName, scenarioName *string
)

func main() {

	filePath = flag.String("file-path", "./", "Path to Test Files.")
	fileName = flag.String("file-name", "", "Name Of Test Files.")
	scenarioName = flag.String("scenario-name", "all", "Tests a specific scenario.")
	//uniquePtr := flag.Bool("unique", false, "Measure unique values of a metric.")

	flag.Parse()

	if strings.Contains(*filePath, ".json") {
		handleSpecificFile(*filePath, *fileName)
		return
	}

	files, err := ioutil.ReadDir(*filePath)
	if err != nil {
		logrus.Fatal(err)
	}

	var wg sync.WaitGroup
	for _, file := range files {
		if strings.Contains(file.Name(), ".json") {
			wg.Add(1)
			go func(fileName string, path string) {
				printFileName(fileName)
				defer wg.Done()
				data, err := ioutil.ReadFile(path + fileName)
				if err != nil {
					logrus.Fatal(err)
				}
				result, ignored, err := handleTests(data)
				if err != nil {
					logrus.Fatal(err)
				}
				utils.FinalResult(file.Name(), ignored, result)

			}(file.Name(), *filePath)
		}
		wg.Wait()
	}
}

func handleSpecificFile(path, fileName string) {
	data, err := ioutil.ReadFile(path + fileName)
	if err != nil {
		logrus.Fatal(err)
	}
	result, ignored, err := handleTests(data)
	if err != nil {
		logrus.Fatal(err)
	}
	utils.FinalResult(fileName, ignored, result)
}

func handleTests(data []byte) (result map[string]bool, ignored int, err error) {
	var scenarios model.TestModel
	err = json.Unmarshal(data, &scenarios)
	if err != nil {
		return nil, 0, err
	}
	result, ignored = service.MakeHTTPCall(scenarios.Test)
	return result, ignored, nil
}

func printFileName(fn string) {
	fmt.Println("##########################################")
	fmt.Printf("Executing Test File: %v\n", fn)
	fmt.Println("##########################################")
}
