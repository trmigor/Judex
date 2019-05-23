package main

import (
	"net/http"		// Server logic
	"log"			// Logs
	"time"			// Timing
	"net"			// Server logic
	"strings"		// Strings split
	"strconv"		// String convertations
	"fmt"			// I/O formatting
	"path/filepath"	// Filepath join
	"io/ioutil"		// Directory scan

	// Database
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
)

// Solution represents data for a single solution
type Solution struct {
	Number int
	Problem int
	User string
	Posting time.Time
	Time float64
	Memory int64
	Status string
	Score int
}

// ProblemPage holds information for page generating
type ProblemPage struct {
	Problem Problem
	Text string
	Input string
	Output string
	TimeLimit float64
	MemoryLimit int
	Examples []IOExample
	Solved bool
	Solutions []Solution
}

// IOExample holds examples for problem solution
type IOExample struct {
	Number int
	Input string
	Output string
}

// SingleProblem handles GET request for problem page
func SingleProblem(w http.ResponseWriter, r *http.Request) {
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

	if problemNumber == "" {
		http.Redirect(w, r, "/problems", http.StatusSeeOther)
	}

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
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}

	// Looking for solutions
	solutionsCollection := client.Database("Judex").Collection("solutions")

	filter = bson.D {
		{Key: "problem", Value: problem.Number},
		{Key: "user", Value: userCredential.Username},
	}

	var solutions []Solution

	cur, err = solutionsCollection.Find(context.TODO(), filter)
	if err != nil && err.Error() != "mongo: no documents in result" {
		log.Println(err)
		return
	}

	for cur.Next(context.TODO()) {
		var elem Solution
		err := cur.Decode(&elem)

		if err != nil {
			log.Println(err)
			return
		}

		solutions = append(solutions, elem)
	}

	if err := cur.Err(); err != nil {
		log.Println(err)
		return
	}

	cur.Close(context.TODO())

	// Preparing the template
	var page ProblemPage

	if len(solutions) == 0 {
		page.Solved = false
	} else {
		page.Solved = true
	}

	page.Problem = problem
	page.Solutions = solutions
	page.Text = FileReader(filepath.Join(problemsPath, strconv.Itoa(problem.Number), "text.txt"))
	page.Input = FileReader(filepath.Join(problemsPath, strconv.Itoa(problem.Number), "input.txt"))
	page.Output = FileReader(filepath.Join(problemsPath, strconv.Itoa(problem.Number), "output.txt"))

	files, err := ioutil.ReadDir(filepath.Join(problemsPath, strconv.Itoa(problem.Number), "examples"))
    if err != nil {
        log.Println(err)
	}

	maxExample := 0

	for _, v := range files {
		number, _ := strconv.Atoi(strings.Split(v.Name(), ".")[0])
		if number > maxExample {
			maxExample = number
		}
	}

	for i := 1; i <= maxExample; i++ {
		var example IOExample

		input := FileReader(filepath.Join(problemsPath, strconv.Itoa(problem.Number), "examples", strconv.Itoa(i) + ".in"))

		output := FileReader(filepath.Join(problemsPath, strconv.Itoa(problem.Number), "examples", strconv.Itoa(i) + ".out"))

		if input == "" || output == "" {
			break
		}

		example.Input = input
		example.Output = output
		example.Number = i
		page.Examples = append(page.Examples, example)
	}

	limits := FileReader(filepath.Join(problemsPath, strconv.Itoa(problem.Number), "limits.lim"))

	if limits != "" {
		fmt.Sscan(strings.Split(limits, "=")[1], &page.TimeLimit)
		fmt.Sscan(strings.Split(limits, "=")[2], &page.MemoryLimit)
	}

	// Executing template
	if err := templates.ExecuteTemplate(w, "problem.html", page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}