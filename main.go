package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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
	r.HandleFunc("/V1/SHOUT", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		decoder := json.NewDecoder(r.Body)
		var shout ShoutRequest
		err := decoder.Decode(&shout)
		if err != nil {
			fmt.Printf("Error json decoding: %v", r.Body)
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

	}).Methods("POST")
	http.Handle("/", r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
