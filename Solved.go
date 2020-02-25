package main

import (
	"log"			// Logs
	"net"			// Server logic
	"net/http"		// Server logic
	"time"			// Timing
	"strings"		// Strings split

	// Database
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

// SolvedPage holds information to show
type SolvedPage struct {
	Any bool
	Solved []ProblemsRow
	User string
}

// Solved handles GET request for page of solved problems
func Solved(w http.ResponseWriter, r *http.Request) {
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

	user := strings.Split(r.URL.Path, "/")[2]

	if user != "" {
		usersCollection := client.Database("Judex").Collection("users")

		var findResult User
		filter = bson.D{{Key: "username", Value: user}}
		err = usersCollection.FindOne(context.TODO(), filter).Decode(&findResult)

		if err != nil {
			ErrorHandler(w, r, http.StatusNotFound)
			return
		} 
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
	var page SolvedPage

	page.User = user
	
	if user == "" {
		page.User = userCredential.Username
	}

	// Looking for scores
	scoresCollection := client.Database("Judex").Collection("scores")

	for _, v := range problems {
		filter = bson.D{
			{Key: "problem", Value: v.Number},
			{Key: "user", Value: userCredential.Username},
		}

		var score Score

		err := scoresCollection.FindOne(context.TODO(), filter).Decode(&score)

		if err != nil && err.Error() != "mongo: no documents in result" {
			log.Println(err)
			return
		}

		if score.Score == 100 {
			page.Solved = append(page.Solved, ProblemsRow{v, score})
		}
	}

	page.Any = true

	if len(page.Solved) == 0 {
		page.Any = false
	}

	// Executing template
	if err := templates.ExecuteTemplate(w, "solved.html", page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
