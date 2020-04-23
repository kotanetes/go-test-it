package utils

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/kotanetes/go-test-it/model"
	"github.com/sirupsen/logrus"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

// Stats is a type of struct and has map in it. where it stores all the
// test results by file name. Key holds file name and value holds results
type Stats struct {
	FileName string
	Total    int
	Passed   int32
	Failed   int32
	Ignored  int
	Status   string
}

// Store adds results to the FinalResults type
func (s *Stats) Store(key string) {
	d[key] = *s
}

// ReadAll return all results
func (s *Stats) ReadAll() map[string]Stats {
	return d
}

var (
	d              = make(map[string]Stats, 0)
	fr             Stats
	passed, failed int32
	ignored        int
)

// PrintResults evaluates and prints final results based on the given data.
func PrintResults(fileName string, result model.HTTPResult) {

	var (
		total      int
		pass, fail int32
	)

	if result.Ignored > 0 {
		ignored += result.Ignored
	}

	total = len(result.TestResults) + result.Ignored

	for _, v := range result.TestResults {
		if v {
			passed++
			pass++
		} else {
			failed++
			fail++
		}
	}
	fmt.Println("##############################################")
	fmt.Printf("Total Tests:%v, Passed:%v, Failed:%v, Ignored:%v\n", total, pass, fail, result.Ignored)
	fmt.Println("##############################################")

	fr.FileName = fileName
	fr.Total = total
	fr.Failed = fail
	fr.Ignored = result.Ignored
	fr.Passed = pass

	if fail > 0 {
		fr.Status = "FAILED"
	}

	fr.Store(fileName)

}

type data struct {
	PageTitle string
	Tests     []Stats
}

// GenerateReport will create a final report as pie chart and save them to a file
func GenerateReport() {

	var (
		status string
		d      data
	)
	path := "./results/"

	if _, err := os.Stat("./results/"); err != nil {
		logrus.Error(err)
		if err = os.Mkdir("results", os.ModePerm); err != nil {
			logrus.Fatal(err)
		}
	}

	defer func() {
		pie := chart.PieChart{
			Width:  250,
			Height: 250,
			Values: []chart.Value{
				{Style: chart.Style{FillColor: drawing.ColorGreen}, Value: float64(passed), Label: fmt.Sprintf("Passed:%v", passed)},
				{Style: chart.Style{FillColor: drawing.ColorRed}, Value: float64(failed), Label: fmt.Sprintf("Failed:%v", failed)},
				{Value: float64(ignored), Label: fmt.Sprintf("Ignored:%v", ignored)},
			},
		}

		resultFile := fmt.Sprintf(path + "results_pie_chart.png")

		f, err := os.Create(resultFile)
		if err != nil {
			logrus.Error(err)
		}
		defer f.Close()
		if err = pie.Render(chart.PNG, f); err != nil {
			logrus.Error(err)
		}
		failed, passed, ignored = 0, 0, 0
	}()

	tmpl := template.New("layout")

	tmpl.Parse(layout)

	mapData := fr.ReadAll()

	d.PageTitle = "Test Results_" + time.Now().Format("2006-01-02 15:04:05")
	for _, v := range mapData {
		d.Tests = append(d.Tests, v)
	}
	file := path + "Test Results.html"
	f, err := os.Create(file)
	if err != nil {
		logrus.Error(err)
	}
	defer f.Close()

	tmpl.Execute(f, d)

	if failed > 0 {
		status = "FAILED"
	} else {
		status = "SUCCESSFUL"
	}
	fmt.Printf("Result: %v\n", status)
}

// InitJSONFile - creates a sample json file to avoid overhead to user
func InitJSONFile() {
	var m = make(map[string]interface{}, 0)

	tests := []model.TestScenario{
		{
			Scenario: "First REST service Test 1",
			Ignore:   false,
			Type:     "REST",
			URL:      "http://sample-url.domain.com/users",
			Method:   "POST",
			Header: model.Header{
				Authorization: "Basic e1f4g5dew==",
				ContentType:   "application/json",
			},
			Body: map[string]interface{}{
				"id":   123456,
				"name": "sample",
			},
			ExpectedStatusCode: 200,
			ExpectedResult: map[string]interface{}{
				"message": "user crested",
			},
		},
		{
			Scenario: "First REST service Test 2",
			Ignore:   false,
			Type:     "REST",
			URL:      "http://sample-url.domain.com/users?id=123456",
			Method:   "GET",
			Header: model.Header{
				Authorization: "Basic e1f4g5dew==",
				ContentType:   "application/json",
			},
			ExpectedStatusCode: 200,
			ExpectedResult: map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"id":   123456,
						"name": "sample",
					},
				},
			},
		},
	}

	m["tests"] = tests

	data, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		logrus.Fatal(err)
	}

	f, err := os.Create("go_test_it_reg.json")
	defer f.Close()
	if err != nil {
		logrus.Fatal(err)
	}
	_, err = f.Write(data)
	if err != nil {
		logrus.Fatal(err)
	}
}
