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

// Home handles GET request for home page
func Home(w http.ResponseWriter, r *http.Request) {
	// Checking credentials
	// If user is not logged in, redirect to start page
	ip, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		log.Println("/home: Cannot discover user's IP")
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	userCredential := Credential {
		UserIP: net.ParseIP(ip),
	}

	if userCredential.UserIP == nil {
		log.Println("/home: Cannot discover user's IP")
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

	// Executing template
	if err := templates.ExecuteTemplate(w, "home.html", ""); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
