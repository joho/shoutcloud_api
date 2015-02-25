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
	r.PathPrefix("/V1/FUCK_OFF").HandlerFunc(FuckOff).Methods("GET")
	http.Handle("/", r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Printf("LISTENING ON %v", port)
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
		log.Printf("ERROR JSON DECODING: %v", r.Body)
		http.Error(w, "BAD REQUEST", http.StatusBadRequest)
		return
	}

	shout.Process()

	json, err := json.Marshal(shout)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("CONTENT-TYPE", "APPLICATION/JSON")
	w.Write(json)
}

func FuckOff(w http.ResponseWriter, r *http.Request) {
	// take path to right of /FUCK_OFF (TODO query string?)
	pathSegments := strings.Split(r.URL.Path, "/")

	foArgs := path.Join(pathSegments[3:len(pathSegments)]...)
	host := "http://foaas.com/"

	// create foaas request
	url := host + foArgs
	req, _ := http.NewRequest("GET", url, nil)
	if r.Header.Get("Accept") != "" {
		req.Header.Set("Accept", r.Header.Get("Accept"))
	}
	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Printf("FOAAS UPSTREAM ERROR: %v", err)
		http.Error(w, "BAD UPSTREAM SERVICE", http.StatusBadGateway)
		return
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	// TODO don't shout the case senstive things like asset paths
	shoutyFuckOff := strings.ToUpper(string(respBody))

	// write response
	w.WriteHeader(resp.StatusCode)
	w.Write([]byte(shoutyFuckOff))
}
