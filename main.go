package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
	"sync"
)

var errChan = make(chan error, 0)

func main() {
	fmt.Println("hello")

	//var testFiles featuresByFiles

	files, err := ioutil.ReadDir("./tests")
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	for _, file := range files {
		fmt.Println(file.Name())
		wg.Add(1)
		go func(fileNmae string) {
			defer wg.Done()
			data, err := ioutil.ReadFile("./tests/" + fileNmae)
			if err != nil {
				log.Fatal(err)
			}
			_ = handleTests(data)

		}(file.Name())
	}
	wg.Wait()
}

func handleTests(data []byte) (err error) {
	var scenarios testModel
	err = json.Unmarshal(data, &scenarios)
	if err != nil {
		return err
	}
	makeHttpCall(scenarios.Test)
	return nil
}

func makeHttpCall(scenarios testScenario) {
	for _, test := range scenarios {
		var (
			reqBody     []byte
			err         error
			requestBody io.Reader
		)

		if len(test.Body) > 0 {
			body := make(map[string]interface{}, 0)
			for _, v := range test.Body {
				body[v.Key] = v.Value
			}
			reqBody, err = json.Marshal(body)
			if err != nil {
				log.Println(err)
			}
			requestBody = strings.NewReader(string(reqBody))
		} else {
			requestBody = nil
		}

		req, err := http.NewRequest(test.Method, test.URL, requestBody)
		if err != nil {
			log.Println(err)
		}

		if test.Header.Authorization != "" {
			req.Header.Add("authorization", test.Header.Authorization)
		}

		bodyBytes, statusCode := doCall(req)
		result := make(map[string]interface{}, 0)
		err = json.Unmarshal(bodyBytes, &result)
		if err != nil {
			log.Println(err)
		}
		if statusCode != test.ExpectedStatusCode || !reflect.DeepEqual(result, test.ExpectedResult) {
			log.Println("test failed")
		}
	}
}

func doCall(req *http.Request) ([]byte, int) {

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return bodyBytes, resp.StatusCode
}

// featuresByFiles - read each files in directory
// decode them to mentioned type
type featuresByFiles map[string]interface{}

type testModel struct {
	Test testScenario `json:"tests"`
}
type testScenario []struct {
	Scenario string `json:"scenario"`
	Type     string `json:"type"`
	URL      string `json:"url"`
	Method   string `json:"method"`
	Header   struct {
		Authorization string `json:"authorization"`
		ContentType   string `json:"content-type"`
	} `json:"header"`
	Body []struct {
		Key   string      `json:"key"`
		Value interface{} `json:"value"`
	} `json:"body"`
	ExpectedStatusCode int         `json:"expectedStatusCode"`
	ExpectedResult     interface{} `json:"expectedResult"`
}
