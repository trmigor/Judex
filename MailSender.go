package main

import (
	"io"
	"os"
	"log"
	"fmt"
	"bufio"
	"net/smtp"
)

// MailSender sends emails
func MailSender(pattern string, user *User) (error) {
	from := "j****t@gmail.com"
	password := "********"

	to := user.Email

	file, err := os.OpenFile(emailPatterns + "/" + pattern, os.O_RDONLY, 0400)
	reader := bufio.NewReader(file)

	if err != nil {
		log.Println("Pattern file open:", err)
		return err
	}

	msg := "From: " + from + "\n" +
		"To: " + to + "\n"

	for {
		input, err := reader.ReadString('\n')

		if err == io.EOF {
			log.Println("EOF")
			break
		}

		if err != nil {
			log.Println("Reading:", err)
			return err
		}

		msg += fmt.Sprint(input)
	}

	err = smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, password, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Println("Smtp:", err)
		return err
	}
	return nil
}