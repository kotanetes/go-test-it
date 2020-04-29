package model2

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/sirupsen/logrus"
)

const (
	// GraphQL - constant used to check test type
	GraphQL = "graphql"
	Failed  = "FAILED"
	Passed  = "PASSED"
	Ignored = "IGNORED"
)

// FeaturesByFiles - read each files in directory
// decode them to mentioned type
type FeaturesByFiles map[string]interface{}

// TestModel - decode json file to this struct model
type TestModel struct {
	// IgnoreFile - attribute in JSON to make sure
	// if the file needs to be ignored or processed
	IgnoreFile bool `json:"ignoreFile"`
	// URL here gives flexibility to group all the scenarios
	// for example: If User wants to test different endpoints
	// running on same host/C-Name and domain. Provide One URL
	// mention endpoint in test scenario then Tool will take care of rest
	URL string `json:"url"`
	// Header added at TestModel level to identify the autorization type
	// based on AuthHeader details "go-test-it" will handle the authorization
	// if given type in AuthHeader is "oauth".
	// tools do a HTTP call to token endpoint
	// caches the token till object session is completed.
	Header AuthHeader `json:"header"`
	// TestScenaris holds the actual endpoint and
	// expected results that are to be compared against the response
	Tests  TestScenarios `json:"tests"`
	status string
	HTTPResult
}

// TestScenarios - slice of test scenarios
type TestScenarios []TestScenario

// TestScenario - slice of tests
type TestScenario struct {
	// Scenario is the name the scenario
	Scenario string `json:"scenario"`
	// Ignore - boolean indicator to ignore the test scenario or not
	Ignore bool `json:"ignore"`
	// Type - type of the API architecture
	// defaults to REST, Provide "graphql" in case of GraphQL
	Type string `json:"type"`
	// URL here gives flexibility to Test the scenario
	// for example: If User wants to test specific endpoints
	// This URL will over-ride the common URL string
	URL string `json:"url"`
	// EndPoint provides the path for the API endpoint to be called
	// EndPoint parameter will be concatinated with the Common URL or
	// Scenario level URL.
	EndPoint string `json:"endPoint"`
	// Method is allowed http method for the given request
	// for emaple:
	// "method": "GET"
	// "method": "POST"
	// "method": "PUT"
	//
	Method string `json:"method"`
	// Header model to handle authorization tokens
	// find further deatils in AuthHeader Struct
	Header AuthHeader `json:"header"`
	// Body - holds the payload of http request
	Body interface{} `json:"body"`
	// ExpectedStatusCode - status code that is expected to return
	ExpectedStatusCode int `json:"expectedStatusCode"`
	// ExpectedResult - data that is expected after making the HTTP call
	// also to compare the each and every field against the retunred http response.
	ExpectedResult     interface{} `json:"expectedResult"`
	ReturnedStatusCode int
	ReturnedResult     map[string]interface{}
}

// HTTPRequest generates the http.Request based on the scenario
func (s *TestScenario) HTTPRequest(m TestModel) (*http.Request, error) {
	logrus.Debugf("Excueting HTTPRequest")
	var (
		reqBody     []byte
		err         error
		requestBody io.Reader
		body        interface{}
	)

	if s.Body != nil {
		switch strings.ToLower(s.Type) {
		case GraphQL:
			bodyMap := make(map[string]interface{})
			bodyMap["query"] = s.Body
			body = bodyMap
		default:
			body = s.Body
		}
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}

		requestBody = strings.NewReader(string(reqBody))
	}

	return http.NewRequest(s.Method, formURL(m, *s), requestBody)
}

// SetHeader set headers
func (s *TestScenario) SetHeader(t TestModel, req *http.Request) {
	logrus.Debugf("Excueting SetHeader")
	if t.Header.Authorization != "" && strings.ToLower(t.Header.Type) == "" {
		req.Header.Add("authorization", t.Header.Authorization)
	}

	req.Header.Set("Content-Type", "application/json")

	if t.Header.CustomHeaders != nil {
		for k, v := range t.Header.CustomHeaders.(map[string]string) {
			req.Header.Set(k, v)
		}
	}
}

func formURL(m TestModel, s TestScenario) (url string) {
	logrus.Debugf("Excueting formURL")
	if m.URL != "" {
		url = m.URL
	}

	if s.URL != "" {
		url = s.URL
	}

	url = url + s.EndPoint

	return url
}

// AuthHeader model to handle authorization tokens
type AuthHeader struct {
	// Type is to identify the type of authorization
	Type string `json:"type"`
	// TokenURL for OAuth Authorization model
	// takes token endpoint and gets the authkey
	// from mentioned token API
	TokenURL string `json:"tokenUrl"`
	// TokenType is used for OAuth model
	// provide TokenType key in JSON file.
	// ex: "tokenType" : "Bearer"
	TokenType string `json:"tokenType"`
	// TokenKey - key parameter that has the auth value
	TokenKey string `json:"tokenKey"`
	// TokenData used for caching OAuth tokens
	// caches the auth token till the object's session
	// is completed or cleared by garbage collector
	TokenData map[string]interface{}
	// Authorization for Basic Auth model
	// Basic Auth is used to make HTTP Call
	// if provided ex: "authorization":"Basic awejdhfasd="
	Authorization string `json:"authorization"`
	// contentType sent as part of HTTP.Request headers
	// deafult contentType is set to "application/json"
	ContentType string `json:"content-type"`
	// TODO: For future USE
	XAuthorization string `json:"x-Authorization"`
	// CustomHeaders adds flexibility to send any additional headers
	// which are used APi specific.
	CustomHeaders interface{} `json:"customHeaders"`
}

// HTTPResult  struct to return
// results of each test scenarios
type HTTPResult struct {
	// FileName - name of the tested file
	FileName string
	// TestResults - key will have the scenario name
	// value is a string has status of the test
	TestResults struct {
		Ignored map[string]string
		Passed  map[string]string
		Failed  map[string]string
	}
	// Ignored - count of ignored tests
	Ignored int
	// Avg Reponse Time for all the scenarios
	// for future USE
	ResponseTime int
	// Error - list of errors for each test scenario
	Errors map[string][]Error
}

// Error - details of the errors like
// rootcause and trace of errors
// used to generate reports
type Error struct {
	Field     string
	Expected  interface{}
	Got       interface{}
	RootCause string
	Trace     string
}

// SetFileName - set file name
func (t *TestModel) SetFileName(name string) {
	logrus.Debugf("Excueting SetFileName")
	t.FileName = name
}

// IsFileIgnored verifies and return bool value
// if the file has been ignored or not
func (t *TestModel) IsFileIgnored() bool {
	logrus.Debugf("Excueting IsFileIgnored")
	return t.IgnoreFile
}

// ExcludeIgnoredScenarios - excludes scenarios that are ignored
// and return tests scenarios to be tested
func (t *TestModel) ExcludeIgnoredScenarios() TestScenarios {
	logrus.Debugf("Excueting ExcludeIgnoredScenarios")
	var toBeTested TestScenarios
	var ignoredTests = make(map[string]string, 0)
	for _, v := range t.Tests {
		switch v.Ignore {
		case false:
			toBeTested = append(toBeTested, v)
		case true:
			ignoredTests[v.Scenario] = Ignored
		}
	}

	if len(ignoredTests) > 0 {
		t.TestResults.Ignored = ignoredTests
	}

	return toBeTested
}

// CompareData compare results expected vs got
func (s *TestScenario) CompareData() (bool, []Error) {
	isValid := true
	var errors []Error
	if diff := cmp.Diff(s.ExpectedResult, s.ReturnedResult); diff != "" {
		isValid = false
		errors = append(errors, Error{RootCause: diff, Trace: "data mismatch (-)want, (+) got"})
	} else if s.ExpectedStatusCode != s.ReturnedStatusCode {
		isValid = false
		errors = append(errors, Error{RootCause: cmp.Diff(s.ExpectedStatusCode, s.ReturnedStatusCode), Trace: "data mismatch (-)want, (+) got"})
	}
	return isValid, errors
}

// ComputeResults - gives the final metrics
func (t *TestModel) ComputeResults() {

}
