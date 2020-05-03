package utils

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"sort"
	"time"

	model "github.com/kotanetes/go-test-it/model2"
	"github.com/sirupsen/logrus"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

// Stats is a type of struct and has map in it. where it stores all the
// test results by file name. Key holds file name and value holds results
type Stats struct {
	FileName string
	Total    int
	Passed   int
	Failed   int
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
	passed, failed int
	ignored        int
)

// PrintResults evaluates and prints final results based on the given data.
func PrintResults(fileName string, result model.TestModel) {

	var (
		total, pass, fail, ignore int
	)

	if len(result.TestResults.Ignored) > 0 {
		ignore = len(result.TestResults.Ignored)
		ignored += len(result.TestResults.Ignored)
	}

	pass = len(result.TestResults.Passed)
	passed += pass
	fail = len(result.TestResults.Failed)
	failed += fail

	total = pass + fail + ignore
	fmt.Println("##############################################")
	fmt.Printf("Total Tests:%v, Passed:%v, Failed:%v, Ignored:%v\n", total, pass, fail, ignore)
	fmt.Println("##############################################")

	fr.FileName = fileName
	fr.Total = total
	fr.Failed = fail
	fr.Ignored = ignore
	fr.Passed = pass
	fr.Status = "PASSED"

	if fail > 0 {
		fr.Status = "FAILED"
	}

	fr.Store(fileName)

	path := "./file_results/"

	if _, err := os.Stat("./file_results/"); err != nil {
		logrus.Error(err)
		if err = os.Mkdir("file_results", os.ModePerm); err != nil {
			logrus.Fatal(err)
		}
	}

	file := path + "result_" + result.FileName
	f, err := os.Create(file)
	if err != nil {
		logrus.Error(err)
	}
	defer f.Close()
	tResult, err := json.MarshalIndent(result.HTTPResult, "", "    ")
	if err != nil {
		logrus.Error(err)
	}
	_, err = f.Write(tResult)
	if err != nil {
		logrus.Error(err)
	}

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
	path := "./reports/"

	if _, err := os.Stat("./reports/"); err != nil {
		logrus.Error(err)
		if err = os.Mkdir("reports", os.ModePerm); err != nil {
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
	var files []string
	d.PageTitle = "Test Results_" + time.Now().Format("2006-01-02 15:04:05")

	for k := range mapData {
		files = append(files, k)
	}
	//sort test cases by name
	sort.Strings(files)

	for _, v := range files {
		d.Tests = append(d.Tests, mapData[v])
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
			Header: model.AuthHeader{
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
			Header: model.AuthHeader{
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
