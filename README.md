# go-test-it
Automated regression testing package written in Go.This tool supports both REST API and GraphQL API. All you need to do is write json files with all the possible scenarios. examples are available in /tests directory.

Installation instructions are listed below.

## Pre-Requsite
  Required go(>1.11) installed on machine.
   * if you are using MAC or Linux, make sure go/bin is added to PATH
 
## Installation
   ```bash
   go get github.com/kotanetes/go-test-it
   ```
## Usage
  * Create a JSON file similar to example file in /test/example.json in any directory
  
  * Need help?
  ```bash
  go-test-it -help
  ```
  * Create a direcotry and add json file, which have all the test scenarios refer to file in /test/example.json
  * By Default, `go-test-it` will read all .json files in the current directory and run the tests
     ```bash
      go-test-it
      ```
  * Tool has the ability to run tests in a specific file.
      ```bash
      go-test-it -file-path=./rest_services -file-name=service1.json
      ```     
  * Ignore Test Scenario
    * Add indicator `"ignore":true` to the test scenario and tool will skip the test scenario
    
 ## Results
  * Tool generates results as a HTMl Report,Pie Chart and also prints the results to console
    example:
    
    Screen Shot 2020-04-21 at 12.05.38 AM
  
