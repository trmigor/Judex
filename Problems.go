package main

import (
	"log"			// Logs
	"net"			// Server logic
	"net/http"		// Server logic
	"time"			// Timing

	// Database
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
)

// Problem represents the problem data
type Problem struct {
	Number int
	Name   string
	Author string
}

// Score holds information about user's score for a problem
type Score struct {
	Problem int
	User    string
	Score   int
}

// ProblemsRow holds information for one output table line
type ProblemsRow struct {
	Problem Problem
	Score   int
}

// Problems handles GET request for problems page
func Problems(w http.ResponseWriter, r *http.Request) {
	// Checking credentials
	// If user is not logged in, redirect to start page
	ip, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		log.Println("/problems: Cannot discover user's IP")
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	userCredential := Credential{
		UserIP: net.ParseIP(ip),
	}

	if userCredential.UserIP == nil {
		log.Println("/problems: Cannot discover user's IP")
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

	// Looking for problems
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

	// Preparing template
	var page []ProblemsRow

	// Looking for scores
	solutionsCollection := client.Database("Judex").Collection("solutions")

	for _, v := range problems {
		filter = bson.D{
			{Key: "problem", Value: v.Number},
			{Key: "user", Value: userCredential.Username},
		}

		cur, err := solutionsCollection.Find(context.TODO(), filter)

		if err != nil && err.Error() != "mongo: no documents in result" {
			log.Println(err)
			return
		}

		maxScore := 0

		for cur.Next(context.TODO()) {
			var elem Solution
			err := cur.Decode(&elem)

			if elem.Score > maxScore {
				maxScore = elem.Score
			}

		}

		page = append(page, ProblemsRow{v, maxScore})
	}

	// Executing template
	if err := templates.ExecuteTemplate(w, "problems.html", page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
