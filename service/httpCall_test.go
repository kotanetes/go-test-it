package service

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/kotanetes/go-test-it/model"
)

type mockClient struct {
	DoReturn *http.Response
	DoError  error
	DoFn     func(*http.Request) (*http.Response, error)
}

func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFn != nil {
		return m.DoFn(req)
	}
	return m.DoReturn, m.DoError
}

func Test_CompareData(t *testing.T) {
	t.Run("should return true for valid map data", func(t *testing.T) {
		var result, expected map[string]interface{}

		result = map[string]interface{}{
			"id":              "16",
			"employee_name":   "Michael Silva",
			"employee_salary": "198500",
			"employee_age":    "66",
			"profile_image":   "",
		}

		expected = map[string]interface{}{
			"id":              "16",
			"employee_name":   "Michael Silva",
			"employee_salary": "198500",
			"employee_age":    "66",
			"profile_image":   "",
		}

		final := true

		r := compareData(t.Name(), result, expected)
		if r != final {
			t.Errorf("Test %v failed, expected not equeal to got", t.Name())
		}
	})

	t.Run("should return true for the data of slice of interface", func(t *testing.T) {
		var result, expected map[string]interface{}

		result = map[string]interface{}{
			"id":              "16",
			"employee_name":   "Michael Silva",
			"employee_salary": "198500",
			"employee_age":    "66",
			"additional_data": []interface{}{
				map[string]interface{}{"address": "line 1"},
				map[string]interface{}{"phone": "12345556"},
			},
		}

		expected = map[string]interface{}{
			"id":              "16",
			"employee_name":   "Michael Silva",
			"employee_salary": "198500",
			"employee_age":    "66",
			"additional_data": []interface{}{
				map[string]interface{}{"address": "line 1"},
				map[string]interface{}{"phone": "12345556"},
			},
		}

		final := true

		r := compareData(t.Name(), result, expected)
		if r != final {
			t.Errorf("%v failed, expected not equal to got", t.Name())
		}
	})

	t.Run("should return false for the data of slice of interface", func(t *testing.T) {
		var result, expected map[string]interface{}

		result = map[string]interface{}{
			"id":              "16",
			"employee_name":   "Michael Silva",
			"employee_salary": "198500",
			"employee_age":    "66",
			"additional_data": []interface{}{
				map[string]interface{}{"address": "line 1"},
				map[string]interface{}{"phone": "1234555"},
			},
		}

		expected = map[string]interface{}{
			"id":              "16",
			"employee_name":   "Michael Silva",
			"employee_salary": "198500",
			"employee_age":    "66",
			"additional_data": []interface{}{
				map[string]interface{}{"address": "line 1"},
				map[string]interface{}{"phone": "12345556"},
			},
		}

		final := false

		r := compareData(t.Name(), result, expected)
		if r != final {
			t.Errorf("%v failed, expected not equal to got", t.Name())
		}
	})

	t.Run("should return false for valid map data", func(t *testing.T) {
		var result, expected map[string]interface{}

		result = map[string]interface{}{
			"id":              "16",
			"employee_name":   "Michael Silva",
			"employee_salary": "198500",
			"employee_age":    "66",
			"profile_image":   "",
		}

		expected = map[string]interface{}{
			"id":              "15",
			"employee_name":   "Michael Silva",
			"employee_salary": "198500",
			"employee_age":    "66",
			"profile_image":   "",
		}

		final := false

		r := compareData(t.Name(), result, expected)
		if r != final {
			t.Errorf("Test %v failed, expected not equeal to got", t.Name())
		}
	})

	t.Run("should return true for the data of embedded of map", func(t *testing.T) {
		var result, expected map[string]interface{}

		result = map[string]interface{}{
			"id":              "16",
			"employee_name":   "Michael Silva",
			"employee_salary": "198500",
			"employee_age":    "66",
			"additional_data": []interface{}{
				map[string]interface{}{"address": map[string]interface{}{"addr 1": "line 1", "zip": "12345"}},
				map[string]interface{}{"phone": "12345556"},
			},
		}

		expected = map[string]interface{}{
			"id":              "16",
			"employee_name":   "Michael Silva",
			"employee_salary": "198500",
			"employee_age":    "66",
			"additional_data": []interface{}{
				map[string]interface{}{"address": map[string]interface{}{"addr 1": "line 1", "zip": "12345"}},
				map[string]interface{}{"phone": "12345556"},
			},
		}

		final := true

		r := compareData(t.Name(), result, expected)
		if r != final {
			t.Errorf("%v failed, expected not equal to got", t.Name())
		}
	})

	t.Run("should return false for the data of embedded of map", func(t *testing.T) {
		var result, expected map[string]interface{}

		result = map[string]interface{}{
			"id":              "16",
			"employee_name":   "Michael Silva",
			"employee_salary": "198500",
			"employee_age":    "66",
			"additional_data": []interface{}{
				map[string]interface{}{"address": map[string]interface{}{"addr 1": "line 1", "zip": "123456"}},
				map[string]interface{}{"phone": "12345556"},
			},
		}

		expected = map[string]interface{}{
			"id":              "16",
			"employee_name":   "Michael Silva",
			"employee_salary": "198500",
			"employee_age":    "66",
			"missing_field":   "",
			"additional_data": []interface{}{
				map[string]interface{}{"address": map[string]interface{}{"addr 1": "line 1", "zip": "12345"}},
				map[string]interface{}{"phone": "12345556"},
			},
		}

		final := false

		r := compareData(t.Name(), result, expected)
		if r != final {
			t.Errorf("%v failed, expected not equal to got", t.Name())
		}
	})

}

func Test_MakeHTTPCall(t *testing.T) {
	t.Run("should ignore the scenario", func(t *testing.T) {
		ignoredCount := 1
		scenario := []model.TestScenario{
			{
				Scenario:           "sample test 4",
				URL:                "https://jsonplaceholder.typicode.com/todos/11",
				Method:             "GET",
				Ignore:             true,
				ExpectedStatusCode: 200,
				ExpectedResult: map[string]interface{}{
					"userId":    1,
					"id":        11,
					"title":     "vero rerum temporibus dolor",
					"completed": true,
				},
			},
		}

		mock := &mockClient{}

		InitHTTPClient(mock)

		_, ignored, _ := MakeHTTPCall(scenario)
		if ignored != ignoredCount {
			t.Errorf("%v failed,scenario is ignored. expected %v, got %v", t.Name(), ignoredCount, ignored)
		}

	})

	t.Run("should make HTTP call and compare result", func(t *testing.T) {
		ignoredCount := 0
		scenarioResult := true
		scenario := []model.TestScenario{
			{
				Scenario:           "mock scenario",
				URL:                "https://mock.typicode.com/todos/11",
				Method:             "GET",
				Ignore:             false,
				ExpectedStatusCode: 200,
				ExpectedResult: map[string]interface{}{
					"userId":    "1",
					"id":        "11",
					"title":     "vero rerum temporibus dolor",
					"completed": true,
				},
			},
		}

		mock := &mockClient{
			DoReturn: &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader(`{"userId":"1","id":"11","title":"vero rerum temporibus dolor","completed": true}`)),
			},
		}

		InitHTTPClient(mock)

		result, ignored, _ := MakeHTTPCall(scenario)
		if ignored != ignoredCount {
			t.Errorf("%v failed,scenario is ignored. expected %v, got %v", t.Name(), ignoredCount, ignored)
		}

		if result[scenario[0].Scenario] != scenarioResult {
			t.Errorf("%v failed,scenario result expected %v, got %v", t.Name(), scenarioResult, result[scenario[0].Scenario])
		}

	})

	t.Run("should make HTTP call and compare result with Auth", func(t *testing.T) {
		ignoredCount := 0
		scenarioResult := false
		scenario := []model.TestScenario{
			{
				Scenario: "mock scenario",
				URL:      "https://mock.typicode.com/todos/11",
				Method:   "GET",
				Ignore:   false,
				Header: model.Header{
					Authorization: "basic Abc12344==",
				},
				ExpectedStatusCode: 200,
				ExpectedResult: map[string]interface{}{
					"userId":    "1",
					"id":        "11",
					"title":     "vero rerum temporibus dolor",
					"completed": true,
				},
			},
		}

		mock := &mockClient{
			DoReturn: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       ioutil.NopCloser(strings.NewReader(`{"userId":"1","id":"11","title":"vero rerum temporibus dolor","completed": true}`)),
			},
		}

		InitHTTPClient(mock)

		result, ignored, _ := MakeHTTPCall(scenario)
		if ignored != ignoredCount {
			t.Errorf("%v failed,scenario is ignored. expected %v, got %v", t.Name(), ignoredCount, ignored)
		}

		if result[scenario[0].Scenario] != scenarioResult {
			t.Errorf("%v failed,scenario result expected %v, got %v", t.Name(), scenarioResult, result[scenario[0].Scenario])
		}

	})

	t.Run("should make HTTP call and return parsing error", func(t *testing.T) {
		errorExpected := true
		scenario := []model.TestScenario{
			{
				Scenario:           "mock scenario",
				URL:                "https://mock.typicode.com/todos/11",
				Method:             "GET",
				Ignore:             false,
				ExpectedStatusCode: 200,
				ExpectedResult: map[string]interface{}{
					"userId":    "1",
					"id":        "11",
					"title":     "vero rerum temporibus dolor",
					"completed": true,
				},
			},
		}

		mock := &mockClient{
			DoReturn: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       ioutil.NopCloser(strings.NewReader(`{"userId":"1","id":"11","title":"vero rerum temporibus dolor","completed": true,}`)),
			},
		}

		InitHTTPClient(mock)

		_, _, err := MakeHTTPCall(scenario)
		if err != nil && !errorExpected {
			t.Errorf("%v failed,scenario returned error %v", t.Name(), err)
		}
	})

	t.Run("should make HTTP call and fail for invalid results", func(t *testing.T) {
		ignoredCount := 0
		scenarioResult := false
		scenario := []model.TestScenario{
			{
				Scenario:           "mock scenario",
				URL:                "https://mock.typicode.com/todos/11",
				Method:             "GET",
				Ignore:             false,
				ExpectedStatusCode: 200,
				ExpectedResult: map[string]interface{}{
					"userId":    "1",
					"id":        "11",
					"title":     "vero rerum temporibus dolor",
					"completed": true,
				},
			},
		}

		mock := &mockClient{
			DoReturn: &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader(`{"userId":"1","id":"12","title":"vero rerum temporibus dolr","completed": true}`)),
			},
		}

		InitHTTPClient(mock)

		result, ignored, _ := MakeHTTPCall(scenario)
		if ignored != ignoredCount {
			t.Errorf("%v failed,scenario is ignored. expected %v, got %v", t.Name(), ignoredCount, ignored)
		}

		if result[scenario[0].Scenario] != scenarioResult {
			t.Errorf("%v failed,scenario result expected %v, got %v", t.Name(), scenarioResult, result[scenario[0].Scenario])
		}

	})

	t.Run("should make HTTP call and compare result for POST Method", func(t *testing.T) {
		ignoredCount := 0
		scenarioResult := true
		scenario := []model.TestScenario{
			{
				Scenario:           "mock method POST scenario",
				URL:                "https://mock.typicode.com/todos/11",
				Method:             "POST",
				Type:               "GraphQL",
				Body:               "query{teacher{firstName}}",
				Ignore:             false,
				ExpectedStatusCode: 200,
				ExpectedResult: map[string]interface{}{
					"userId":    "1",
					"id":        "11",
					"title":     "vero rerum temporibus dolor",
					"completed": true,
				},
			},
		}

		mock := &mockClient{
			DoReturn: &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader(`{"userId":"1","id":"11","title":"vero rerum temporibus dolor","completed": true}`)),
			},
		}

		InitHTTPClient(mock)

		result, ignored, _ := MakeHTTPCall(scenario)
		if ignored != ignoredCount {
			t.Errorf("%v failed,scenario is ignored. expected %v, got %v", t.Name(), ignoredCount, ignored)
		}

		if result[scenario[0].Scenario] != scenarioResult {
			t.Errorf("%v failed,scenario result expected %v, got %v", t.Name(), scenarioResult, result[scenario[0].Scenario])
		}

	})

	t.Run("should make HTTP call and compare result for POST Method 2", func(t *testing.T) {
		ignoredCount := 0
		scenarioResult := true
		scenario := []model.TestScenario{
			{
				Scenario:           "mock method POST scenario",
				URL:                "https://mock.typicode.com/todos/11",
				Method:             "POST",
				Body:               "{\"userId\":\"1\",\"id\":\"11\"}",
				Ignore:             false,
				ExpectedStatusCode: 200,
				ExpectedResult: map[string]interface{}{
					"userId":    "1",
					"id":        "11",
					"title":     "vero rerum temporibus dolor",
					"completed": true,
				},
			},
		}

		mock := &mockClient{
			DoReturn: &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader(`{"userId":"1","id":"11","title":"vero rerum temporibus dolor","completed": true}`)),
			},
		}

		InitHTTPClient(mock)

		result, ignored, _ := MakeHTTPCall(scenario)
		if ignored != ignoredCount {
			t.Errorf("%v failed,scenario is ignored. expected %v, got %v", t.Name(), ignoredCount, ignored)
		}

		if result[scenario[0].Scenario] != scenarioResult {
			t.Errorf("%v failed,scenario result expected %v, got %v", t.Name(), scenarioResult, result[scenario[0].Scenario])
		}

	})

	t.Run("should make HTTP call and return parsing request body error", func(t *testing.T) {
		errorExpected := true
		scenario := []model.TestScenario{
			{
				Scenario:           "mock method POST scenario",
				URL:                "https://mock.typicode.com/todos/11",
				Method:             "POST",
				Body:               "{\"userId\":\"1\",\"id\":\"11\",}",
				Ignore:             false,
				ExpectedStatusCode: 200,
				ExpectedResult: map[string]interface{}{
					"userId":    "1",
					"id":        "11",
					"title":     "vero rerum temporibus dolor",
					"completed": true,
				},
			},
		}

		mock := &mockClient{
			DoReturn: &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader(`{"userId":"1","id":"11","title":"vero rerum temporibus dolor","completed": true,}`)),
			},
		}

		InitHTTPClient(mock)

		_, _, err := MakeHTTPCall(scenario)
		if err != nil && !errorExpected {
			t.Errorf("%v failed,scenario returned error %v", t.Name(), err)
		}
	})
}
