package main

import (
	"log"			// Logs
	"net"			// Server logic
	"net/http"		// Server logic
	"time"			// Timing
	"sort"			// Sorting

	// Database
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
)

// UsersRow holds information about user
type UsersRow struct {
	Number int
	User string
	Solved int
}

// Users handles GET request for users page
func Users(w http.ResponseWriter, r *http.Request) {
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

	var page []UsersRow

	usersCollection := client.Database("Judex").Collection("users")
	scoresCollection := client.Database("Judex").Collection("scores")

	filter = bson.D{}

	cur, err := usersCollection.Find(context.TODO(), filter)
	if err != nil && err.Error() != "mongo: no documents in result" {
		log.Println(err)
		return
	}

	for cur.Next(context.TODO()) {
		var elem User
		err := cur.Decode(&elem)

		if err != nil {
			log.Println(err)
			return
		}

		filter = bson.D{ {Key: "user", Value: elem.Username}, {Key: "score", Value: 100},}

		var solved []Score

		sc, err := scoresCollection.Find(context.TODO(), filter)
		if err != nil && err.Error() != "mongo: no documents in result" {
			log.Println(err)
			return
		}

		for sc.Next(context.TODO()) {
			var score Score
			err := cur.Decode(&score)

			if err != nil {
				log.Println(err)
				return
			}

			solved = append(solved, score)
		}

		if err := cur.Err(); err != nil {
			log.Println(err)
			return
		}
	
		cur.Close(context.TODO())

		page = append(page, UsersRow{ User: elem.Username, Solved: len(solved)})
	}

	if err := cur.Err(); err != nil {
		log.Println(err)
		return
	}

	cur.Close(context.TODO())

	sort.Slice(page, func(i, j int) bool { return page[i].Solved > page[j].Solved })

	for i := range page {
		page[i].Number = i + 1
	}

	// Executing template
	if err := templates.ExecuteTemplate(w, "users.html", page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
