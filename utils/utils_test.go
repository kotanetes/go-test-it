package utils

import (
	"os"
	"testing"

	"github.com/kotanetes/go-test-it/model"
)

func init() {
	fr = Stats{FileName: "test.json", Total: 5, Passed: 2, Failed: 2, Ignored: 1, Status: "FAILED"}
}

func Test_GenerateReports(t *testing.T) {
	fr.Store("test.json")
	GenerateReport()
	_, err := os.Stat("./results/")
	if err != nil {
		t.Fatal(err)
	}
	os.RemoveAll("./results/")
}

func Test_PrintResults(t *testing.T) {
	fileName := "test.file"
	ignored := 0
	result := map[string]bool{
		"test 1": true,
		"test 2": false,
	}

	hr := model.HTTPResult{TestResults: result, Ignored: ignored}

	PrintResults(fileName, hr)
}
