package main

import (
	"net/http"
	"net"
	"log"
	"go.mongodb.org/mongo-driver/bson"
	"context"
	"time"
	"strings"
	"strconv"
	"os"
	"path/filepath"
	"encoding/csv"
)

// SolutionPage holds information for page generating
type SolutionPage struct {
	Problem Problem
	PostTime time.Time
	Compiler string
	User string
	Tests []struct {
		Number string
		Verdict string
		Time string
		Memory string
		Passed bool
	}
	Passed int
	Score int
	MaxTime float64
	MaxMemory int64
	Code string
}

// SingleSolution handles GET request for solution page
func SingleSolution(w http.ResponseWriter, r *http.Request) {
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

	// Checking problem existence
	solutionURL := strings.Split(r.URL.Path, "/")[2]

	if solutionURL == "" {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}

	solutionsCollection := client.Database("Judex").Collection("solutions")

	solutionNumber, _ := strconv.Atoi(solutionURL)

	filter = bson.D {
		{Key: "number", Value: solutionNumber},
	}

	var result Solution

	err = solutionsCollection.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if result.User != userCredential.Username {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	problemsCollection := client.Database("Judex").Collection("problems")

	filter = bson.D {
		{Key: "number", Value: result.Problem},
	}

	var problem Problem

	problemsCollection.FindOne(context.TODO(), filter).Decode(&problem)

	page := SolutionPage {
		Problem: problem,
		PostTime: result.Posting,
		Compiler: "gcc",
		User: result.User,
		MaxMemory: result.Memory,
		MaxTime: result.Time,
		Score: result.Score,
	}

	file, err := os.OpenFile(filepath.Join(solutionsPath, solutionURL, solutionURL + ".prot"), os.O_RDWR, 0666)
	
	if err != nil {
		log.Println(err)
		return
	}

	reader := csv.NewReader(file)
	
	protocol, _ := reader.ReadAll()

	page.Tests = make([]struct {
		Number string
		Verdict string
		Time string
		Memory string
		Passed bool
	}, len(protocol))

	passed := 0

	for _, v := range protocol {
		testNumber, _ := strconv.Atoi(v[0])
		page.Tests[testNumber-1].Number = v[0]
		page.Tests[testNumber-1].Time = v[2]
		page.Tests[testNumber-1].Memory = v[3]
		if v[4] == "OK" {
			if v[5] != "OK" {
				page.Tests[testNumber-1].Verdict = "Wrong answer"
				continue
			}
			passed++
			page.Tests[testNumber-1].Passed = true
		}
		page.Tests[testNumber-1].Verdict = v[4]
	}

	page.Passed = passed

	code := FileReader(filepath.Join(solutionsPath, solutionURL, "sol_" + solutionURL + ".c"))

	page.Code = code

	// Executing template
	if err := templates.ExecuteTemplate(w, "solution.html", page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}