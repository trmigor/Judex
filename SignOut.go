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

// SignOut handles request for signing out
func SignOut(w http.ResponseWriter, r *http.Request) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		log.Println("/SignIn_submit: Cannot discover user's IP")
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	userCredential := Credential{
		UserIP:    net.ParseIP(ip),
		EnterTime: time.Now(),
		EndTime:   time.Now().Add(CredCur),
	}

	if userCredential.UserIP == nil {
		log.Println("/SignIn_submit: Cannot discover user's IP")
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	credentialsCollection := client.Database("Judex").Collection("credentials")

	filter := bson.D{{"userip", userCredential.UserIP}}
	credentialsCollection.DeleteMany(context.TODO(), filter)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}