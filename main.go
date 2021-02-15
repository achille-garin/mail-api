package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-gomail/gomail"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
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
	w.Header().Set("Access-Control-Allow-Origin", "https://achille.garin.xyz")

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

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("MAIL_TO"))
	m.SetHeader("To", os.Getenv("MAIL_TO"))
	m.SetHeader("Subject", "Contact from website -> "+mail.Subject)
	m.SetBody("text/plain", "Email : "+mail.Contact+"\r\n"+mail.Body)

	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))

	d := gomail.NewDialer(os.Getenv("MAIL_SERVER"), port, os.Getenv("MAIL_USERNAME"), os.Getenv("MAIL_PASSWORD"))

	if err := d.DialAndSend(m); err != nil {
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
	_, err := os.Stat("/etc/letsencrypt/live/achille.garin.xyz/fullchain.pem")
	if os.IsNotExist(err) {
		log.Fatal(http.ListenAndServe(":8000", nil))
	} else {
		log.Fatal(http.ListenAndServeTLS(":8000", "/etc/letsencrypt/live/achille.garin.xyz/fullchain.pem", "/etc/letsencrypt/live/achille.garin.xyz/privkey.pem", nil))
	}
}

func main() {
	log.Print("Server is running")
	handleRequests()
}
