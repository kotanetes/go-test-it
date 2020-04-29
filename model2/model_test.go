package model2

import (
	"strings"
	"testing"
)

func Test_ExcludeIgnoredScenarios(t *testing.T) {
	t.Run("should return all excluded scenarios", func(t *testing.T) {
		exTestsCount := 3
		inTestCount := 2

		tests := &TestModel{
			Tests: []TestScenario{
				{
					Scenario: "Exclude this scenario 1",
					Ignore:   true,
				},
				{
					Scenario: "Exclude this scenario 2",
					Ignore:   true,
				},
				{
					Scenario: "Exclude this scenario 3",
					Ignore:   true,
				},
				{
					Scenario: "Include this scenario 1",
					Ignore:   false,
				},
				{
					Scenario: "Include this scenario 2",
					Ignore:   false,
				},
			},
		}

		tests.SetFileName("unit tests")

		res := tests.ExcludeIgnoredScenarios()

		if len(res) != inTestCount {
			t.Errorf("%v: expected %v, got %v", t.Name(), inTestCount, len(res))
		}
		if len(tests.TestResults.Ignored) != exTestsCount {
			t.Errorf("%v: expected %v, got %v", t.Name(), exTestsCount, len(tests.TestResults.Ignored))
		}

	})
}

func Test_IsFileIgnored(t *testing.T) {
	t.Run("should return bool value based on file indicator", func(t *testing.T) {
		tests := &TestModel{
			IgnoreFile: true,
		}

		res := tests.IsFileIgnored()

		if !res {
			t.Errorf("%v: expected %v, got %v", t.Name(), true, res)
		}
	})
}

func Test_formURL(t *testing.T) {
	t.Run("should return URL using url pattern and enpoint", func(t *testing.T) {
		expected := "https://unit-testing.mock.com/endpoint?id=1234"

		tm := TestModel{URL: "https://unit-testing.mock.com"}
		ts := TestScenario{EndPoint: "/endpoint?id=1234"}

		url := formURL(tm, ts)
		if !strings.EqualFold(url, expected) {
			t.Errorf("%v: expected %v, got %v", t.Name(), expected, url)
		}
	})

	t.Run("should return over rided URL using url in test scenario and enpoint", func(t *testing.T) {
		expected := "https://unit-testing.mock1.com/endpoint?id=1234"

		tm := TestModel{URL: "https://unit-testing.mock.com"}
		ts := TestScenario{URL: "https://unit-testing.mock1.com", EndPoint: "/endpoint?id=1234"}

		url := formURL(tm, ts)
		if !strings.EqualFold(url, expected) {
			t.Errorf("%v: expected %v, got %v", t.Name(), expected, url)
		}
	})
}

func Test_HttpRequest(t *testing.T) {
	t.Run("should return HTTP request for default with out error", func(t *testing.T) {
		ts := TestModel{
			Tests: TestScenarios{
				{
					Scenario: "mock unit test 1",
					URL:      "https://unit-testing.mock1.com",
					EndPoint: "/endpoint?id=1234",
					Method:   "POST",
					Body: map[string]interface{}{
						"id": 1234,
					},
				},
			},
		}
		_, err := ts.Tests[0].HTTPRequest(ts)
		if err != nil {
			t.Errorf("%v:%v", t.Name(), err)
		}
	})

	t.Run("should return HTTP request for graphql with out error", func(t *testing.T) {
		ts := TestModel{
			Tests: TestScenarios{
				{
					Scenario: "mock unit test 1",
					URL:      "https://unit-testing.mock1.com",
					EndPoint: "/endpoint?id=1234",
					Method:   "POST",
					Type:     GraphQL,
					Body:     "{query{mock(\"id\":1234){id name}}",
				},
			},
		}
		_, err := ts.Tests[0].HTTPRequest(ts)
		if err != nil {
			t.Errorf("%v:%v", t.Name(), err)
		}
	})
}

// func Test_SetHeaders(t *testing.T) {
// 	t.Run("should return ")
// }
