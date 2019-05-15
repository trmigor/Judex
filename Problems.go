package main

import (
	"net/http"
	"log"
	"time"
	"net"

	// Database
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
)

// Problem represents the problem data
type Problem struct {
	Number int
	Name string
	Author string
}

// Score holds information about user's score for a problem
type Score struct {
	Problem int
	User string
	Score int
}

// Row holds information for one output table line
type Row struct {
	Problem Problem
	Score Score
}

// Problems handles GET request for problems page
func Problems(w http.ResponseWriter, r *http.Request) {
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

	problemsCollection := client.Database("Judex").Collection("problems")

	filter = bson.D{}

	var problems []Problem

	cur, err := problemsCollection.Find(context.TODO(), filter)
	if err != nil && err.Error() != "mongo: no documents in result" {
		log.Println(err)
		return
	}

	for cur.Next(context.TODO()) {
		var elem Problem
		err := cur.Decode(&elem)

		if err != nil {
			log.Println(err)
			return
		}

		problems = append(problems, elem)
	}

	if err := cur.Err(); err != nil {
		log.Println(err)
		return
	}

	cur.Close(context.TODO())

	var page []Row

	scoresCollection := client.Database("Judex").Collection("scores")
	
	for _, v := range problems {
		filter = bson.D {
			{Key: "problem", Value: v.Number},
			{Key: "user", Value: userCredential.Username},
		}

		var score Score

		err := scoresCollection.FindOne(context.TODO(), filter).Decode(&score)

		if err != nil && err.Error() != "mongo: no documents in result" {
			log.Println(err)
			return
		}

		page = append(page, Row {v, score})
	} 

	if err := templates.ExecuteTemplate(w, "problems.html", page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}