package main

import (
	"net/http"

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

// Problems handles GET request for problems page
func Problems(w http.ResponseWriter, r *http.Request) {
	problemsCollection := client.Database("Judex").Collection("problems")

	filter := bson.D{}

	var problems []Problem

	cur, err := problemsCollection.Find(context.TODO(), filter)
	if err != nil {

	}

	for cur.Next(context.TODO()) {
		var elem Problem
		err := cur.Decode(&elem)

		if err != nil {

		}

		problems = append(problems, elem)
	}

	if err := cur.Err(); err != nil {

	}

	cur.Close(context.TODO())

	if err := templates.ExecuteTemplate(w, "problems.html", problems); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}