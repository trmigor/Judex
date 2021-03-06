package main

import (
	"log"      // Logs
	"net"      // Server logic
	"net/http" // Server logic
	"time"     // Timing

	// Database
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

// Profile handles GET request for profile page
func Profile(w http.ResponseWriter, r *http.Request) {
	// Checking credentials
	// If user is not logged in, redirect to start page
	ip, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		log.Println("/profile: Cannot discover user's IP")
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	userCredential := Credential{
		UserIP: net.ParseIP(ip),
	}

	if userCredential.UserIP == nil {
		log.Println("/profile: Cannot discover user's IP")
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	credentialsCollection := client.Database("Judex").Collection("credentials")

	filter := bson.D{{Key: "userip", Value: userCredential.UserIP}}
	err = credentialsCollection.FindOne(context.TODO(), filter).Decode(&userCredential)

	if err == nil {
		if time.Now().After(userCredential.EndTime) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			credentialsCollection.DeleteMany(context.TODO(), filter)
			return
		}
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	// Looking for user information
	usersCollection := client.Database("Judex").Collection("users")

	var findResult User
	filter = bson.D{{Key: "username", Value: userCredential.Username}}
	err = usersCollection.FindOne(context.TODO(), filter).Decode(&findResult)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Executing template
	if err := templates.ExecuteTemplate(w, "profile.html", findResult); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}