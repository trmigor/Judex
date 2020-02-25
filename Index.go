package main

import (
	"log"       // Logs
	"net"       // Server logic
	"net/http"  // Server logic
	"time"      // Timing

	// Database
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

// Index handles GET request for start page
func Index(w http.ResponseWriter, r *http.Request) {
	// 404 error handle
	if r.URL.Path != "/" {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}

	// Checking credentials
	// If user is already logged in, redirect to /home
	ip, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		log.Println("/: Cannot discover user's IP")
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	userCredential := Credential{
		UserIP: net.ParseIP(ip),
	}

	if userCredential.UserIP == nil {
		log.Println("/: Cannot discover user's IP")
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

	// Executing template
	if err := templates.ExecuteTemplate(w, "index.html", ""); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
