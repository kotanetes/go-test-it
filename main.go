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

var errChan = make(chan error, 0)

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
	filePath, scenarioName *string
)

func main() {
	//var testFiles featuresByFiles

	filePath = flag.String("file-path", "./", "Path to Test Files.")
	scenarioName = flag.String("scenario-name", "all", "Tests a specific scenario.")
	//uniquePtr := flag.Bool("unique", false, "Measure unique values of a metric.")

	flag.Parse()

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
				result, err := handleTests(data)
				if err != nil {
					logrus.Fatal(err)
				}
				utils.FinalResult(result)

			}(file.Name(), *filePath)
		}
		wg.Wait()
	}
}

func handleTests(data []byte) (result map[string]bool, err error) {
	var scenarios model.TestModel
	err = json.Unmarshal(data, &scenarios)
	if err != nil {
		return nil, err
	}
	return service.MakeHTTPCall(scenarios.Test), nil
}

func printFileName(fn string) {
	fmt.Println("##########################################")
	fmt.Printf("Executing Test File: %v\n", fn)
	fmt.Println("##########################################")
}
