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

	model "github.com/kotanetes/go-test-it/model2"
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
	logrus.SetLevel(logrus.InfoLevel)

	// intializing remote service
	service.InitHTTPClient(&http.Client{})
}

// variables used for flag arguments
var (
	filePath, fileName, scenarioName *string
	debugMode                        *bool
)

func main() {

	filePath = flag.String("file-path", "./", "Path to Test Files.")
	fileName = flag.String("file-name", "", "Name of Test Files.")
	scenarioName = flag.String("scenario-name", "all", "Tests a specific scenario.")
	debugMode = flag.Bool("d", false, "debug logs console")

	if *debugMode {
		logrus.SetLevel(logrus.DebugLevel)
	}

	args := os.Args[1:]

	if len(args) > 0 {
		if args[0] == "init" {
			utils.InitJSONFile()
			return
		}
	}

	flag.Parse()

	if strings.EqualFold(*filePath, "./") && *fileName != "" {
		printFileName(*fileName)
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
					logrus.WithField("file-name", file.Name()).Fatal(err)
				}
				result, err := handleTests(data, fileName)
				if err != nil {
					logrus.WithField("file-name", file.Name()).Fatal(err)
				}
				utils.PrintResults(file.Name(), result)

			}(file.Name(), *filePath)
		}
		wg.Wait()
	}
	utils.GenerateReport()
}

// handleSpecificFile reads specific json file from the given path
// thos function read the file and calls the underlying common function.
func handleSpecificFile(path, fileName string) {
	data, err := ioutil.ReadFile(path + fileName)
	if err != nil {
		logrus.Fatal(err)
	}
	result, err := handleTests(data, fileName)
	if err != nil {
		logrus.Fatal(err)
	}
	utils.PrintResults(fileName, result)
	utils.GenerateReport()
}

// handleTests unmarshals byte data to TestModel type and pass the scenarios
// to MakeHTTPCall function that makes calls to URL mentioned in tests.
func handleTests(data []byte, name string) (result model.TestModel, err error) {
	var scenarios model.TestModel
	err = json.Unmarshal(data, &scenarios)
	if err != nil {
		return result, err
	}
	scenarios.SetFileName(name)
	result, err = service.MakeHTTPCall(scenarios)
	return result, err
}

// printFileName prints the file name just before the execution starts
func printFileName(fn string) {
	fmt.Println("##########################################")
	fmt.Printf("Executing Test File: %v\n", fn)
	fmt.Println("##########################################")
}
