# Integration Test Outcome Logger 

A small go program that stores the outcome of QA integration tests for auditing purposes.

### Features
 * Data can be retrieved as png.
 * Store different environments eg. qa and staging.

This just allows us to see when the tests were last run on a bamboo page.

## Getting Started

```
// Install the dependancies
$ go get github.com/fogleman/gg
$ go get github.com/gorilla/mux

// Build and run
$ go build src/main.go src/data.go
$ ./main
```

#### Different environments

```
// This is how you build for a different enviroment
$ env GOOS=linux GOARCH=amd64 go build src/main.go src/data.go
```

## Usage

This contains some example useage

```
// Create a qa ping for the test service.
curl -X POST localhost:8000/ping/qa/test -H 'Content-type: application/json' -d '{"passed": true}'

// Get all pings
curl localhost:8000/ping

// Get all pings for environment
curl localhost:8000/ping/qa

// Get all pings for environment and service
curl localhost:8000/ping/qa/test

// Get the image for service with last pass and fail
curl localhost:8000/ping/qa/test.png
```

## Example

![An example of the output on bamboo](bambooexample.png)
