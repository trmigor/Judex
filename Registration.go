package main

import (
	"log"      // Logs
	"net"      // Server logic
	"net/http" // Server logic
	"time"     // Timing

	// Database
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
)

// RegistrationPage is a structure for displaying sign up page
type RegistrationPage struct {
	UsernameError string
	EmailError string
}

// Registration handles GET request for registration page
func Registration(w http.ResponseWriter, r *http.Request) {
	// Checking credentials
	// If user is already logged in, redirect to /home
	ip, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		log.Println("/registration: Cannot discover user's IP")
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	userCredential := Credential{
		UserIP: net.ParseIP(ip),
	}

	if userCredential.UserIP == nil {
		log.Println("/registration: Cannot discover user's IP")
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	credentialsCollection := client.Database("Judex").Collection("credentials")

	filter := bson.D{{Key: "userip", Value: userCredential.UserIP}}
	err = credentialsCollection.FindOne(context.TODO(), filter).Decode(&userCredential)

	if err == nil {
		if time.Now().Before(userCredential.EndTime) {
			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		}
		credentialsCollection.DeleteMany(context.TODO(), filter)
	}

	page := RegistrationPage{
		UsernameError: "none",
		EmailError: "none",
	}

	// Executing template
	if err := templates.ExecuteTemplate(w, "registration.html", page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
