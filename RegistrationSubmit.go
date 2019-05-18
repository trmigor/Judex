package main

import (
	"log"      // Logs
	"net"      // Server logic
	"net/http" // Server logic
	"time"     // Timing
	"fmt"	   // I/O formatting

	// Database
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
)

// RegSubmit handles POST request for registration submit
func RegSubmit(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		log.Println("/reg_submit: Cannot parse registration form")
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	// Checking if user wants to access /reg_submit without filling a form
	if len(r.Form) == 0 {
		log.Println("/reg_submit: Blank form")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	formResult := User{
		Username:   r.Form["username"][0],
		Email:      r.Form["email"][0],
		Password:   r.Form["password"][0],
		FirstName:  r.Form["f_name"][0],
		MiddleName: r.Form["m_name"][0],
		LastName:   r.Form["l_name"][0],
		Company:    r.Form["company"][0],
		Website:    r.Form["website"][0],
		Bio:        r.Form["bio"][0],
	}

	// Checking reCAPTCHA
	recaptchaResponse, responseFound := r.Form["g-recaptcha-response"]
	if responseFound {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Println("/sign_in_submit: Cannot discover user's IP")
			ErrorHandler(w, r, http.StatusInternalServerError)
			return
		}
		if !ProcessRequest(recaptchaResponse[0], ip) {
			log.Println("Wrong reCAPTCHA")
			fmt.Fprintln(w, "Wrong reCAPTCHA")
			return
		}
	}

	usersCollection := client.Database("Judex").Collection("users")

	// Checking the uniqueness of the username
	var findResult User
	filter := bson.D{{Key: "username", Value: formResult.Username}}
	err = usersCollection.FindOne(context.TODO(), filter).Decode(&findResult)

	if err == nil {
		log.Println("/reg_submit: Username", formResult.Username, "is not unique")
		ErrorHandler(w, r, DoubledUsername)
		return
	}

	// Checking the uniqueness of the email
	filter = bson.D{{Key: "email", Value: formResult.Email}}
	err = usersCollection.FindOne(context.TODO(), filter).Decode(&findResult)

	if err == nil {
		log.Println("/reg_submit: Email", formResult.Email, "is already used")
		ErrorHandler(w, r, DoubledEmail)
		return
	}

	// Saving received information
	_, err = usersCollection.InsertOne(context.TODO(), formResult)
	if err != nil {
		log.Println("/reg_submit: users: Cannot insert information")
		ErrorHandler(w, r, InsertError)
		return
	}

	log.Println("Registered. Username:", formResult.Username)

	MailSender("Welcome.txt", &formResult)

	// Creating a credential
	ip, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		log.Println("/reg_submit: Cannot discover user's IP")
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	userCredential := Credential{
		UserIP:    net.ParseIP(ip),
		Username:  formResult.Username,
		EnterTime: time.Now(),
		EndTime:   time.Now().Add(CredCur),
	}

	if userCredential.UserIP == nil {
		log.Println("/reg_submit: Cannot discover user's IP")
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	credentialsCollection := client.Database("Judex").Collection("credentials")

	filter = bson.D{{Key: "userip", Value: userCredential.UserIP}}
	credentialsCollection.DeleteMany(context.TODO(), filter)

	_, err = credentialsCollection.InsertOne(context.TODO(), userCredential)

	if err != nil {
		log.Println("/reg_submit: credentials: Cannot insert information")
		ErrorHandler(w, r, InsertError)
		return
	}

	log.Println("Logged in. Username:", formResult.Username)

	http.Redirect(w, r, "/home", http.StatusSeeOther)
	return
}
