package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/kotanetes/go-test-it/model"
	"github.com/kotanetes/go-test-it/service"
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
}

func main() {
	fmt.Println("hello")

	//var testFiles featuresByFiles

	files, err := ioutil.ReadDir("./tests")
	if err != nil {
		logrus.Fatal(err)
	}

	var wg sync.WaitGroup
	for _, file := range files {
		fmt.Println(file.Name())
		wg.Add(1)
		go func(fileNmae string) {
			defer wg.Done()
			data, err := ioutil.ReadFile("./tests/" + fileNmae)
			if err != nil {
				logrus.Fatal(err)
			}
			_ = handleTests(data)

		}(file.Name())
	}
	wg.Wait()
}

func handleTests(data []byte) (err error) {
	var scenarios model.TestModel
	err = json.Unmarshal(data, &scenarios)
	if err != nil {
		return err
	}
	service.MakeHttpCall(scenarios.Test)
	return nil
}
