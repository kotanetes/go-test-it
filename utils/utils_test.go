package utils

import (
	"os"
	"testing"

	model "github.com/kotanetes/go-test-it/model2"
)

func init() {
	fr = Stats{FileName: "test.json", Total: 5, Passed: 2, Failed: 2, Ignored: 1, Status: "FAILED"}
}

func Test_GenerateReports(t *testing.T) {
	fr.Store("test.json")
	GenerateReport()
	_, err := os.Stat("./reports/")
	if err != nil {
		t.Fatal(err)
	}
	os.RemoveAll("./reports/")
}

func Test_PrintResults(t *testing.T) {
	fileName := "test.file"
	result := map[string]string{
		"test 1": model.Passed,
		"test 2": model.Passed,
	}

	hr := model.TestModel{}
	hr.TestResults.Passed = result

	PrintResults(fileName, hr)
	_, err := os.Stat("./file_results/")
	if err != nil {
		t.Fatal(err)
	}
	os.RemoveAll("./file_results/")
}

func Test_InitJSONFile(t *testing.T) {
	InitJSONFile()
	fileName := "go_test_it_reg.json"
	_, err := os.Stat(fileName)
	if err != nil {
		t.Fatal(err)
	}
	os.Remove(fileName)
}
