package main

import (
	"testing"
)

// Call sendMail with an unusable email address
func TestSendMailInvalidEmail(t *testing.T) {
	var mail = Mail{
		Contact: "jean@com_pany.com",
		Subject: "Job offer for you",
		Body:    "Hello my name is Jean, I message you to pay you a lot of money",
	}
	_, code := sendMail(mail)
	if code != 400 {
		t.Fatalf(`Invalid email address are allowed`)
	}
}

// Call sendMail with an empty string as parameters
func TestSendMailEmpty(t *testing.T) {
	var mail = Mail{
		Contact: "jean@company.com",
		Subject: "Job offer for you",
		Body:    "",
	}
	_, code := sendMail(mail)
	if code != 400 {
		t.Fatalf(`Empty string are allowed on defaults params`)
	}
}
