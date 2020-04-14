package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/wcharczuk/go-chart"
)

// FinalResult - prints final result and create pie chart
func FinalResult(fileName string, ignored int, result map[string]bool) {

	var (
		total          int
		passed, failed int32
	)

	total = len(result)

	for _, v := range result {
		if v {
			passed++
		} else {
			failed++
		}
	}
	fmt.Println("##############################################")
	fmt.Printf("Total Tests:%v, Passed:%v, Failed:%v, Ignored:%v\n", total, passed, failed, ignored)
	fmt.Println("##############################################")

	generateReport(fileName, ignored, passed, failed)

	if failed > 0 {
		fmt.Printf("Result: %v\n", "FAILED")
	}
}

func generateReport(fileName string, ignored int, passed, failed int32) {
	pie := chart.PieChart{
		Title:  "Test Result",
		Width:  512,
		Height: 512,
		Values: []chart.Value{
			{Style: chart.Style{Show: true}, Value: float64(passed), Label: "Passed"},
			{Style: chart.Style{Show: true}, Value: float64(failed), Label: "Failed"},
			{Style: chart.Style{Show: true}, Value: float64(ignored), Label: "Ignored"},
		},
	}

	fileStrings := strings.Split(fileName, ".")
	resultFile := fmt.Sprintf("%v.png", fileStrings[0])

	f, err := os.Create(resultFile)
	if err != nil {
		logrus.Error(err)
	}
	defer f.Close()
	if err = pie.Render(chart.PNG, f); err != nil {
		logrus.Error(err)
	}
}
