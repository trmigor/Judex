package main

import (
	"log"				// Logs
	"net/smtp"			// Mail
	"path/filepath"		// Filepath join
)

// MailSender sends emails
func MailSender(pattern string, user *User) (error) {
	from := "j****t@gmail.com"
	password := "********"

	to := user.Email

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" + FileReader(filepath.Join(emailPatterns, pattern))

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, password, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Println("MailSender: Smtp:", err)
		return err
	}
	return nil
}