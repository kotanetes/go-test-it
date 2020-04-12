package utils

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/wcharczuk/go-chart"
)

// FinalResult - prints final result and create pie chart
func FinalResult(result map[string]bool) {

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
	fmt.Println("##########################################")
	fmt.Printf("Total Tests:%v, Passed:%v, Failed:%v\n", total, passed, failed)
	fmt.Println("##########################################")

	generateReport(passed, failed)

	if failed > 0 {
		fmt.Errorf("Result: FAILED")
	}
}

func generateReport(passed, failed int32) {
	pie := chart.PieChart{
		Title:  "Test Result",
		Width:  512,
		Height: 512,
		Values: []chart.Value{
			{Value: float64(passed), Label: "Passed"},
			{Value: float64(failed), Label: "Failed"},
		},
	}

	f, err := os.Create("report.png")
	if err != nil {
		logrus.Error(err)
	}
	defer f.Close()
	pie.Render(chart.PNG, f)
}
