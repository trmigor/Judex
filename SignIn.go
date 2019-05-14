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

// SignInPage is a structure for displaying sign in page
type SignInPage struct {
	Username      string
	UsernameError string
	PasswordError string
}

// SignIn handles GET request for sign in page
func SignIn(w http.ResponseWriter, r *http.Request) {
	// Checking credentials
	// If user is already logged in, redirect to /home
	ip, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		log.Println("/SignIn: Cannot discover user's IP")
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	userCredential := Credential{
		UserIP: net.ParseIP(ip),
	}

	if userCredential.UserIP == nil {
		log.Println("/SignIn: Cannot discover user's IP")
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

	page := SignInPage{
		Username:      "",
		UsernameError: "none",
		PasswordError: "none",
	}

	if err := templates.ExecuteTemplate(w, "sign_in.html", page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
