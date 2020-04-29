package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	model "github.com/kotanetes/go-test-it/model2"
	"github.com/sirupsen/logrus"
)

const (
	// GraphQL - constant used to check test type
	GraphQL       = "graphql"
	expectedVsGot = "value for field %v: expected %v ,got %v"
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
func MakeHTTPCall(t model.TestModel) (model.TestModel, error) {
	logrus.Debugf("Excuting MakeHTTPCall")
	var (
		result                    map[string]interface{}
		finalResult, failedResult = make(map[string]string), make(map[string]string)
		testErrors                = make(map[string][]model.Error)
	)

	testScenarios := t.ExcludeIgnoredScenarios()

	for _, test := range testScenarios {
		req, err := test.HTTPRequest(t)
		if err != nil {
			logrus.Error(err)
			return t, err
		}

		test.SetHeader(t, req)

		bodyBytes, statusCode := doCall(req)
		err = json.Unmarshal(bodyBytes, &result)
		if err != nil {
			logrus.Error(err)
			return t, err
		}

		test.ReturnedStatusCode = statusCode
		test.ReturnedResult = result

		if ok, compareError := test.CompareData(); !ok {
			logrus.WithField("scenario", test.Scenario).Info("failed")
			failedResult[test.Scenario] = model.Failed
			testErrors[test.Scenario] = compareError
		} else {
			finalResult[test.Scenario] = model.Passed
		}

	}
	t.TestResults.Passed = finalResult
	t.TestResults.Failed = failedResult
	t.Errors = testErrors

	return t, nil
}
