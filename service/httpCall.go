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

func doCall(req *http.Request) ([]byte, int) {

	client := &http.Client{}

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

// MakeHttpCall - performs an http call
func MakeHttpCall(scenarios model.TestScenario) {
	for _, test := range scenarios {
		var (
			reqBody     []byte
			err         error
			requestBody io.Reader
			body        interface{}
		)

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
		result := make(map[string]interface{}, 0)
		err = json.Unmarshal(bodyBytes, &result)
		if err != nil {
			logrus.Error(err)
		}

		switch {
		case statusCode != test.ExpectedStatusCode:
			logrus.Info(fmt.Sprintf("Test %v failed, expected status %v got %v", test.Scenario, test.ExpectedStatusCode, statusCode))
		case !reflect.DeepEqual(result, test.ExpectedResult):
			logrus.Info(fmt.Sprintf("est %v failed, expected result is not as return", test.Scenario))
		default:
			logrus.Info(fmt.Sprintf("Test %v Passed", test.Scenario))
		}
	}
}
