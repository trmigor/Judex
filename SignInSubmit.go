package main

import (
	"log"      // Logs
	"net"      // Server logic
	"net/http" // Server logic
	"time"     // Timing
	"fmt"	   // I/O formatting

	// Database
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

// SignInCheck is a structure for checking sign in information
type SignInCheck struct {
	Username string
	Password string
}

// SignInSubmit handles POST request for sign in submit
func SignInSubmit(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Println("/sign_in_submit: Cannot parse sign in form")
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	// Checking if user wants to access /sign_in_submit without filling a form
	if len(r.Form) == 0 {
		log.Println("/sign_in_submit: Blank form")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	formResult := SignInCheck{
		Username: r.Form["username"][0],
		Password: r.Form["password"][0],
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

	// Checking the information
	usersCollection := client.Database("Judex").Collection("users")

	var findResult User
	filter := bson.D{{Key: "username", Value: formResult.Username}}
	err = usersCollection.FindOne(context.TODO(), filter).Decode(&findResult)

	if err != nil {
		log.Println("/sign_in_submit: No such username:", formResult.Username)
		ErrorHandler(w, r, NoUsername)
		return
	}

	if formResult.Password != findResult.Password {
		log.Println("/sign_in_submit: Password do not match")
		ErrorHandler(w, r, WrongPassword)
		return
	}

	// Creating a credential
	ip, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		log.Println("/sign_in_submit: Cannot discover user's IP")
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
		log.Println("/sign_in_submit: Cannot discover user's IP")
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	credentialsCollection := client.Database("Judex").Collection("credentials")

	filter = bson.D{{Key: "userip", Value: userCredential.UserIP}}
	credentialsCollection.DeleteMany(context.TODO(), filter)

	_, err = credentialsCollection.InsertOne(context.TODO(), userCredential)

	if err != nil {
		log.Println("/sign_in_submit: credentials: Cannot insert information")
		ErrorHandler(w, r, InsertError)
		return
	}

	log.Println("Logged in. Username:", formResult.Username)

	http.Redirect(w, r, "/home", http.StatusSeeOther)
	return
}
