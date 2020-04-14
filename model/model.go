package model

// FeaturesByFiles - read each files in directory
// decode them to mentioned type
type FeaturesByFiles map[string]interface{}

// TestModel - decode json file to this struct model
type TestModel struct {
	Test []TestScenario `json:"tests"`
}

// TestScenario - slice of tests
type TestScenario struct {
	Scenario string `json:"scenario"`
	Ignore   bool   `json:"ignore"`
	Type     string `json:"type"`
	URL      string `json:"url"`
	Method   string `json:"method"`
	Header   struct {
		Authorization string `json:"authorization"`
		ContentType   string `json:"content-type"`
	} `json:"header"`
	Body               interface{} `json:"body"`
	ExpectedStatusCode int         `json:"expectedStatusCode"`
	ExpectedResult     interface{} `json:"expectedResult"`
}
