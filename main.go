package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gorilla/mux"
)

type ShoutRequest struct {
	Input  string `json:"INPUT"`
	Output string `json:"OUTPUT"`
}

func (s *ShoutRequest) Process() {
	s.Output = strings.ToUpper(s.Input)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/V1/SHOUT", ShoutBack).Methods("POST")
	r.HandleFunc("/V1/FUCK_OFF", FuckOff).Methods("POST")
	http.Handle("/", r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Printf("Listening on %v", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func ShoutBack(w http.ResponseWriter, r *http.Request) {
	log.Printf("POST /V1/SHOUT %v", r.RemoteAddr)

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var shout ShoutRequest
	err := decoder.Decode(&shout)
	if err != nil {
		log.Printf("Error json decoding: %v", r.Body)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	shout.Process()

	json, err := json.Marshal(shout)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func FuckOff(w http.ResponseWriter, r *http.Request) {
	// take path to right of /FUCK_OFF (query string?)
	pathSegments := strings.Split(r.URL.Path, "/")

	// TODO off by one?
	foArgs := path.Join(pathSegments[2:len(pathSegments)]...)
	host := "http://foaas.com"

	// create foaas request
	resp, err := http.Get(host + foArgs)
	defer resp.Body.Close()
	if err != nil {
		log.Printf("foaas upstream error: %v", err)
		http.Error(w, "Bad Upstream Service", http.StatusBadGateway)
		return
	}

	respBody, err := ioutil.ReadAll(r.Body)
	shoutyFuckOff := strings.ToUpper(string(respBody))

	// write response
	w.WriteHeader(resp.StatusCode)
	w.Write([]byte(shoutyFuckOff))
}
