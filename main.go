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
	"path"
	"regexp"
)

var emailRegex string = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

type Mail struct {
	Contact string `json:"Contact"`
	Subject string `json:"Subject"`
	Body    string `json:"Body"`
}

func init() {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dir := path.Dir(ex)
	err = godotenv.Load(dir + "/.env")
	if err != nil {
		err = godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}

func sendMailRoute(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "achille.garin.xyz:443")

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

	log.Print(fmt.Sprintf("Request proccessing, return message : %q", resp["message"]))
	json.NewEncoder(w).Encode(resp)
}

func sendMail(mail Mail) (string, int) {
	if valid := checkEmptyString(mail); !valid {
		return "Empty string are not allowed", 400
	}
	if valid := checkEmailFormat(mail.Contact); !valid {
		return "Invalid email address", 400
	}

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

func checkEmptyString(mail Mail) bool {
	if mail.Contact == "" || mail.Subject == "" || mail.Body == "" {
		return false
	}
	return true
}

func checkEmailFormat(email string) bool {
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func handleRequests() {
	http.HandleFunc("/", sendMailRoute)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func main() {
	log.Print("Server is running")
	handleRequests()
}
