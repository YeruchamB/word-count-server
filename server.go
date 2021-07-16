package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/vladimirvivien/automi/emitters"
	"log"
	"net/http"
	"os"
	"strings"
)

type Server struct {
	Router *mux.Router
}

type CountRequest struct {
	Input string `json:"input"`
}

type StatsResponse struct {
	Word  string `json:"word"`
	Count int64 `json:"count"`
}

func readFromBody(text string) error {
	return InputWordCount(strings.NewReader(text))
}

func readFromFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	return InputWordCount(emitters.Scanner(f, bufio.ScanLines))
}

func readFromUrl(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	return InputWordCount(emitters.Reader(resp.Body))
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	if payload != nil {
		response, _ := json.Marshal(payload)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

// Handler for adding input to the word count
func (s *Server) countHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the input from the body
	var countRequest CountRequest
	err := json.NewDecoder(r.Body).Decode(&countRequest)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Determine input type based on the 'input' path param
	switch r.URL.Query().Get("input") {
	case "text":
		err = readFromBody(countRequest.Input)
	case "file":
		err = readFromFile(countRequest.Input)
	case "url":
		err = readFromUrl(countRequest.Input)
	default:
		respondWithError(w, http.StatusBadRequest, "Invalid input type")
		return
	}

	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	} else {
		respondWithJSON(w, http.StatusAccepted, nil)
	}
}

// Handler for retrieving the stats for a given word
func (s *Server) statsHandler(w http.ResponseWriter, r *http.Request) {
	word := mux.Vars(r)["word"]

	// Retrieve the count from the cache
	count, err := GetCount(strings.ToLower(word))
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	respondWithJSON(w, http.StatusOK, StatsResponse{Word: word, Count: count})
}

func (s *Server) Initialize() {
	InitRedis()

	s.Router = mux.NewRouter()
	s.Router.HandleFunc("/count", s.countHandler).Methods("POST")
	s.Router.HandleFunc("/stats/{word:[a-zA-Z]+}", s.statsHandler).Methods("GET")
}

func (s *Server) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, s.Router))
}
