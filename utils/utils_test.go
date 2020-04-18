package utils

import (
	"os"
	"testing"
)

func Test_GenerateReports(t *testing.T) {
	fileName := "test.file"
	ignored := 0
	var passed, failed int32 = 10, 5
	generateReport(fileName, ignored, passed, failed)
	_, err := os.Stat("test.png")
	if err != nil {
		t.Fatal(err)
	}
	os.Remove("test.png")
}

func Test_FinalResult(t *testing.T) {
	fileName := "test.file"
	ignored := 0
	result := map[string]bool{
		"test 1": true,
		"test 2": false,
	}

	FinalResult(fileName, ignored, result)
	_, err := os.Stat("test.png")
	if err != nil {
		t.Fatal(err)
	}
	os.Remove("test.png")
}
