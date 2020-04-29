package service

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	model "github.com/kotanetes/go-test-it/model2"
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
		tm := model.NewTestModel()
		tm.Tests = scenario
		res, _ := MakeHTTPCall(tm)
		if len(res.TestResults.Ignored) != ignoredCount {
			t.Errorf("%v failed,scenario is ignored. expected %v, got %v", t.Name(), ignoredCount, len(res.TestResults.Ignored))
		}

	})

	t.Run("should make HTTP call and compare result", func(t *testing.T) {
		ignoredCount := 0
		scenarioResult := model.Passed
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

		tm := model.TestModel{Tests: scenario}
		result, _ := MakeHTTPCall(tm)
		if len(result.TestResults.Ignored) != ignoredCount {
			t.Errorf("%v failed,scenario is ignored. expected %v, got %v", t.Name(), ignoredCount, len(result.TestResults.Ignored))
		}

		if result.TestResults.Passed[scenario[0].Scenario] != scenarioResult {
			t.Errorf("%v failed,scenario result expected %v, got %v", t.Name(), scenarioResult, result.TestResults.Passed[scenario[0].Scenario])
		}

	})

	t.Run("should make HTTP call and compare result with Auth", func(t *testing.T) {
		ignoredCount := 0
		scenarioResult := model.Failed
		scenario := []model.TestScenario{
			{
				Scenario: "mock scenario",
				URL:      "https://mock.typicode.com/todos/11",
				Method:   "GET",
				Ignore:   false,
				Header: model.AuthHeader{
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

		tm := model.TestModel{Tests: scenario}
		result, _ := MakeHTTPCall(tm)
		if len(result.TestResults.Ignored) != ignoredCount {
			t.Errorf("%v failed,scenario is ignored. expected %v, got %v", t.Name(), ignoredCount, len(result.TestResults.Ignored))
		}

		if result.TestResults.Failed[scenario[0].Scenario] != scenarioResult {
			t.Errorf("%v failed,scenario result expected %v, got %v", t.Name(), scenarioResult, result.TestResults.Failed[scenario[0].Scenario])
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

		tm := model.TestModel{Tests: scenario}
		_, err := MakeHTTPCall(tm)
		if err != nil && !errorExpected {
			t.Errorf("%v failed,scenario returned error %v", t.Name(), err)
		}
	})

	t.Run("should make HTTP call and fail for invalid results", func(t *testing.T) {
		ignoredCount := 0
		scenarioResult := model.Failed
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

		tm := model.TestModel{Tests: scenario}
		result, _ := MakeHTTPCall(tm)
		if len(result.TestResults.Ignored) != ignoredCount {
			t.Errorf("%v failed,scenario is ignored. expected %v, got %v", t.Name(), ignoredCount, len(result.TestResults.Ignored))
		}

		if result.TestResults.Failed[scenario[0].Scenario] != scenarioResult {
			t.Errorf("%v failed,scenario result expected %v, got %v", t.Name(), scenarioResult, result.TestResults.Failed[scenario[0].Scenario])
		}

	})

	t.Run("should make HTTP call and compare result for POST Method", func(t *testing.T) {
		ignoredCount := 0
		scenarioResult := model.Passed
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

		tm := model.TestModel{Tests: scenario}
		result, _ := MakeHTTPCall(tm)
		if len(result.TestResults.Ignored) != ignoredCount {
			t.Errorf("%v failed,scenario is ignored. expected %v, got %v", t.Name(), ignoredCount, len(result.TestResults.Ignored))
		}

		if result.TestResults.Passed[scenario[0].Scenario] != scenarioResult {
			t.Errorf("%v failed,scenario result expected %v, got %v", t.Name(), scenarioResult, result.TestResults.Passed[scenario[0].Scenario])
		}

	})

	t.Run("should make HTTP call and compare result for POST Method 2", func(t *testing.T) {
		ignoredCount := 0
		scenarioResult := model.Passed
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

		tm := model.TestModel{Tests: scenario}
		result, _ := MakeHTTPCall(tm)
		if len(result.TestResults.Ignored) != ignoredCount {
			t.Errorf("%v failed,scenario is ignored. expected %v, got %v", t.Name(), ignoredCount, result.Ignored)
		}

		if result.TestResults.Passed[scenario[0].Scenario] != scenarioResult {
			t.Errorf("%v failed,scenario result expected %v, got %v", t.Name(), scenarioResult, result.TestResults.Passed[scenario[0].Scenario])
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

		tm := model.TestModel{Tests: scenario}
		_, err := MakeHTTPCall(tm)
		if err != nil && !errorExpected {
			t.Errorf("%v failed,scenario returned error %v", t.Name(), err)
		}
	})
}
