package main

import (
	"encoding/json"
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

func checkStatsResponse(t *testing.T, expectedCount int64, body []byte) {
	var stats StatsResponse
	err := json.Unmarshal(body, &stats)
	if err != nil {
		t.Error(err)
	} else if stats.Count != expectedCount {
		t.Errorf("Expected 2. Got %d", stats.Count)
	}
}

func Test_Text(t *testing.T) {
	req, _ := http.NewRequest("POST", "/count?input=text", strings.NewReader(`{"Input":"test test hello world"}`))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusAccepted, response.Code)

	req, _ = http.NewRequest("GET", "/stats/test", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	checkStatsResponse(t, 2, response.Body.Bytes())
}

func Test_File(t *testing.T) {
	req, _ := http.NewRequest("POST", "/count?input=file", strings.NewReader(`{"Input":"tests/file1.txt"}`))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusAccepted, response.Code)

	req, _ = http.NewRequest("GET", "/stats/hello", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	checkStatsResponse(t, 33, response.Body.Bytes())
}

//
//func Test_URL(t *testing.T) {
//	req, _ := http.NewRequest("POST", "/count?input=url", strings.NewReader(`{"Input":"tests/file1.txt"}`))
//	response := executeRequest(req)
//	checkResponseCode(t, http.StatusAccepted, response.Code)
//
//	req, _ = http.NewRequest("GET", "/stats/hello", nil)
//	response = executeRequest(req)
//	checkResponseCode(t, http.StatusOK, response.Code)
//	checkStatsResponse(t, 33, response.Body.Bytes())
//}
