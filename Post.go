package main

import (
	"net/http"		// Server logic
	"log"			// Logs
	"time"			// Timing
	"net"			// Server logic
	"strings"		// Strings split
	"strconv"		// String convertations

	// Database
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

// Post handles GET request for page to post a solution
func Post(w http.ResponseWriter, r *http.Request) {
	// Checking credentials
	// If user is not logged in, redirect to start page
	ip, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		log.Println("/problem/: Cannot discover user's IP")
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	userCredential := Credential {
		UserIP: net.ParseIP(ip),
	}

	if userCredential.UserIP == nil {
		log.Println("/problem/: Cannot discover user's IP")
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

	// Checking problem existence
	problemNumber := strings.Split(r.URL.Path, "/")[2]

	if problemNumber != "" {
		problemsCollection := client.Database("Judex").Collection("problems")

		filter = bson.D{}

		var problem Problem
		found := false

		cur, err := problemsCollection.Find(context.TODO(), filter)
		if err != nil && err.Error() != "mongo: no documents in result" {
			log.Println(err)
			return
		}

		for cur.Next(context.TODO()) {
			var elem Problem
			err := cur.Decode(&elem)

			if i, err := strconv.Atoi(problemNumber); err == nil {
				if i == elem.Number {
					problem = elem
					found = true
					break
				}
			} else {
				log.Println(err)
			}

			if err != nil {
				log.Println(err)
				return
			}
		}

		if err := cur.Err(); err != nil {
			log.Println(err)
			return
		}

		cur.Close(context.TODO())

		if !found {
			http.Redirect(w, r, "/post", http.StatusSeeOther)
		}

		if err := templates.ExecuteTemplate(w, "post.html", problem.Number); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := templates.ExecuteTemplate(w, "post.html", ""); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}