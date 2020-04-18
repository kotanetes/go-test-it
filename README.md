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
  ```bash
  go-test-it
  ```
  * Need help?
  ```bash
  go-test-it -help
  ```
  
