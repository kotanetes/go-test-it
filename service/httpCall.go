package service

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/kotanetes/go-test-it/model"
	"github.com/sirupsen/logrus"
)

const (
	// GraphQL - constant used to check test type
	GraphQL = "graphql"
)

var client remoteClient

type remoteClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// InitHTTPClient - intializes httpClient
// this can be used to mock the http calls for this project
func InitHTTPClient(c remoteClient) {
	client = c
}

// doCall does a http call to service.
// read the body from http.Response
// return bosy as byte array.
func doCall(req *http.Request) ([]byte, int) {

	resp, err := client.Do(req)
	if err != nil {
		logrus.Error(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
	}

	return bodyBytes, resp.StatusCode
}

// MakeHTTPCall receives the test scenarios and this function handles the scenarios
// 1. Identify the scenarios are that ignored and seperate rest of them.
// 2. Loop over the tests and based on the test type form the request body.
// 3. Create new http Request and make a call to service.
// 4. Read the http.Response and compare the results with the assertions in scenarios.
// 5. Generate the Test results and return thwm along with ignored scenarios count.
func MakeHTTPCall(scenarios []model.TestScenario) (map[string]bool, int) {
	var (
		reqBody     []byte
		err         error
		requestBody io.Reader
		body        interface{}
		ignored     int
	)

	finalResult := make(map[string]bool)
	testScenarios := make(map[string]model.TestScenario, 0)

	for _, scenario := range scenarios {
		if !scenario.Ignore {
			testScenarios[scenario.Scenario] = scenario
		} else {
			logrus.Info(fmt.Sprintf("Test %v Ignored", scenario.Scenario))
			ignored++
		}
	}

	for _, test := range testScenarios {
		result := make(map[string]interface{})
		if test.Body != nil {
			switch test.Type {
			case GraphQL:
				bodyMap := make(map[string]interface{})
				bodyMap["query"] = test.Body
				body = bodyMap
			default:
				body = test.Body
			}
			reqBody, err = json.Marshal(body)
			if err != nil {
				logrus.Error(err)
			}

			requestBody = strings.NewReader(string(reqBody))
		}

		req, err := http.NewRequest(test.Method, test.URL, requestBody)
		if err != nil {
			logrus.Error(err)
		}

		if test.Header.Authorization != "" {
			req.Header.Add("authorization", test.Header.Authorization)
		}
		req.Header.Set("Content-Type", "application/json")

		bodyBytes, statusCode := doCall(req)

		err = json.Unmarshal(bodyBytes, &result)
		if err != nil {
			logrus.Error(err)
		}

		switch {
		case statusCode != test.ExpectedStatusCode:
			logrus.Info(fmt.Sprintf("Test %v failed, expected status %v got %v", test.Scenario, test.ExpectedStatusCode, statusCode))
			finalResult[test.Scenario] = false
		case !reflect.DeepEqual(result, test.ExpectedResult):
			fmt.Printf("expected:%v,got: %v\n", test.ExpectedResult, result)
			logrus.Info(fmt.Sprintf("Test %v failed, retunred response is not as expected", test.Scenario))
			finalResult[test.Scenario] = false
		default:
			logrus.Info(fmt.Sprintf("Test %v Passed", test.Scenario))
			finalResult[test.Scenario] = true
		}
	}
	return finalResult, ignored
}
