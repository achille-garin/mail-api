package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
)

type Mail struct {
	Contact string `json:"Contact"`
	Subject string `json:"Subject"`
	Body    string `json:"Body"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func sendMailRoute(w http.ResponseWriter, r *http.Request) {

	resp := make(map[string]string)
	w.Header().Set("content-type", "application/json")

	if r.Method == "POST" {
		reqBody, _ := ioutil.ReadAll(r.Body)
		var mail Mail
		json.Unmarshal(reqBody, &mail)

		message, code := sendMail(mail)
		resp["message"] = message
		w.WriteHeader(code)

	} else {
		resp["message"] = fmt.Sprintf("This URL does not accept %q requests", r.Method)
		w.WriteHeader(400)
	}

	json.NewEncoder(w).Encode(resp)
}

func sendMail(mail Mail) (string, int) {
	auth := smtp.PlainAuth("", os.Getenv("MAIL_USERNAME"), os.Getenv("MAIL_PASSWORD"), os.Getenv("MAIL_SERVER"))

	to := []string{os.Getenv("MAIL_TO")}
	msg := []byte("To: " + os.Getenv("MAIL_TO") + "\r\n" +
		"From: " + mail.Contact + "\r\n" +
		"Subject: Contact from website -> " + mail.Subject + "\r\n" +
		"\r\n" +
		mail.Body + "\r\n")
	err := smtp.SendMail(os.Getenv("MAIL_SERVER")+":"+os.Getenv("MAIL_PORT"), auth, mail.Contact, to, msg)
	if err != nil {
		return "Something went wrong when sending the mail", 500
	}
	return "Success", 200
}

func handleRequests() {
	http.HandleFunc("/", sendMailRoute)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func main() {
	handleRequests()
}
