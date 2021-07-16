package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var s Server

func TestMain(m *testing.M) {
	s.Initialize()
	os.Exit(m.Run())
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func checkStats(t *testing.T, word string, expectedCount int64) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/stats/%s", word), nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var stats StatsResponse
	err := json.Unmarshal(response.Body.Bytes(), &stats)
	if err != nil {
		t.Error(err)
	} else if stats.Count != expectedCount {
		t.Errorf("Expected %d. Got %d", expectedCount, stats.Count)
	}
}

func Test_Text(t *testing.T) {
	req, _ := http.NewRequest("POST", "/count?input=text", strings.NewReader(`{"input":"Hi! My name is(what?), my name is(who?), my name is Slim Shady"}`))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusAccepted, response.Code)
	checkStats(t, "my", 3)

	req, _ = http.NewRequest("POST", "/count?input=text", strings.NewReader(`{"input":"Hi! My name is(what?), my name is(who?), my name is Slim Shady"}`))
	response = executeRequest(req)
	checkResponseCode(t, http.StatusAccepted, response.Code)
	checkStats(t, "my", 6)
}

func Test_File(t *testing.T) {
	req, _ := http.NewRequest("POST", "/count?input=file", strings.NewReader(`{"input":"tests/file1.txt"}`))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusAccepted, response.Code)
	checkStats(t, "hello", 32)
}

func Test_URL(t *testing.T) {
	req, _ := http.NewRequest("POST", "/count?input=url", strings.NewReader(`{"input":"https://raw.githubusercontent.com/YeruchamB/word-count-server/main/tests/file2.txt"}`))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusAccepted, response.Code)
	checkStats(t, "drift", 8)
}

func Test_Errors(t *testing.T) {
	// Bad input type
	req, _ := http.NewRequest("POST", "/count?input=bad", strings.NewReader(`{"input":"Hi! My name is(what?), my name is(who?), my name is Slim Shady"}`))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	// Bad body
	req, _ = http.NewRequest("POST", "/count?input=bad", strings.NewReader(`invalid json`))
	response = executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	// Missing word variable
	req, _ = http.NewRequest("GET", "/stats/", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
