package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Mail struct {
	Contact string `json:"Contact"`
	Subject string `json:"Subject"`
	Body    string `json:"Body"`
}

func sendMailRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		w.Header().Set("content-type", "application/json")

		reqBody, _ := ioutil.ReadAll(r.Body)
		var mail Mail
		json.Unmarshal(reqBody, &mail)

		json.NewEncoder(w).Encode(mail)

	} else {
		fmt.Fprintf(w, "This URL does not accept %q requests", r.Method)
	}
}

func handleRequests() {
	http.HandleFunc("/", sendMailRoute)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func main() {
	handleRequests()
}
