package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/trmigor/Judex/testing_packages/compileandrun"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"net"
	"encoding/csv"
)

// SolutionForm holds information received from a request
type SolutionForm struct {
	Problem    int
	Compiler   string
	SourceText string
}

// SolutionSubmit handles POST request of sending a solution, checks it and forms a protocol
func SolutionSubmit(w http.ResponseWriter, r *http.Request) {
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

	err = r.ParseMultipartForm(32 << 20)

	if err != nil {
		log.Println("/reg_submit: Cannot parse registration form:", err)
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	// Checking if user wants to access /reg_submit without filling a form
	if len(r.Form) == 0 {
		log.Println("/reg_submit: Blank form")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	problemNumber, _ := strconv.Atoi(r.Form["problem"][0])

	formResult := SolutionForm{
		Problem: problemNumber,
	}

	solutionsCollection := client.Database("Judex").Collection("solutions")

	cur, err := solutionsCollection.Find(context.TODO(), bson.D{})

	if err != nil {
		log.Println(err)
		return
	}

	maxNumber := 0
	maxScore := 0

	for cur.Next(context.TODO()) {
		var elem Solution

		err := cur.Decode(&elem)

		if err != nil {
			log.Println(err)
			return
		}

		if elem.Number > maxNumber {
			maxNumber = elem.Number
		}

		if elem.Score > maxScore {
			maxScore = elem.Score
		}
	}

	maxNumber++

	err = os.Mkdir(filepath.Join(solutionsPath, strconv.Itoa(maxNumber)), 0777)

	if err != nil {
		log.Println(err)
		return
	}

	err = os.Mkdir(filepath.Join(tmpPath, strconv.Itoa(maxNumber)), 0777)

	if err != nil {
		log.Println(err)
		return
	}

	file, header, err := r.FormFile("sol_file")

	if err != nil {
		log.Println(err)
		return
	}

	defer file.Close()

	extension := strings.Split(header.Filename, ".")[1]

	filePath := filepath.Join(solutionsPath, strconv.Itoa(maxNumber), "sol_"+strconv.Itoa(maxNumber)+"."+extension)

	f1, err := os.OpenFile(filePath, os.O_RDWR | os.O_CREATE, 0666)

	if err != nil {
		log.Println(err)
		return
	}

	defer f1.Close()

	f2, err := os.OpenFile(filepath.Join(tmpPath, strconv.Itoa(maxNumber), "sol_" + strconv.Itoa(maxNumber)+ "." + extension), os.O_RDWR | os.O_CREATE, 0666)

	if err != nil {
		log.Println(err)
		return
	}

	defer f2.Close()

	msg, err := ioutil.ReadAll(file)

	f1.Write(msg)

	f2.Write(msg)

	files, err := ioutil.ReadDir(filepath.Join(problemsPath, strconv.Itoa(formResult.Problem), "tests"))

	if err != nil {
		log.Println(err)
		return
	}

	limits := FileReader(filepath.Join(problemsPath, strconv.Itoa(formResult.Problem), "limits.lim"))

	var TL time.Duration
	var ML int64

	if limits != "" {
		fmt.Sscan(strings.Split(limits, "=")[1], &TL)
		fmt.Sscan(strings.Split(limits, "=")[2], &ML)
	}

	c := compileandrun.Init{
		Solution:    maxNumber,
		Format:      "." + extension,
		Path:        filepath.Join(tmpPath, strconv.Itoa(maxNumber)),
		Compiler:    "gcc",
		TestsPath:   filepath.Join(problemsPath, strconv.Itoa(formResult.Problem), "tests"),
		TestsNumber: len(files) / 2,
		RunLimits: compileandrun.Limits{
			TL:  TL * time.Second,
			ML:  ML * 1024 * 1024,
			RTL: 10 * time.Second,
		},
	}
	p, err := c.Compile()

	if err != nil {
		log.Println(err)
		return
	}

	p.Wait()

	err = c.Run()

	if err != nil {
		log.Println(err)
		return
	}

	protocol, err := os.OpenFile(filepath.Join(tmpPath, strconv.Itoa(maxNumber), strconv.Itoa(maxNumber) + ".prot"), os.O_RDWR | os.O_CREATE, 0666)

	if err != nil {
		log.Println(err)
		return
	}

	defer protocol.Close()

	csvReader := csv.NewReader(protocol)

	results, err := csvReader.ReadAll()

	if err != nil {
		log.Println(err)
		return
	}

	var maxTime float64
	var maxMemory int64

	var ok int

	status := ""

	for _, v := range results {
		var testTime float64
		fmt.Sscan(v[2], &testTime)
		if testTime > maxTime {
			maxTime = testTime
		}

		var testMem int64
		fmt.Sscan(v[3], &testMem)
		if testMem > maxMemory {
			maxMemory = testMem
		}

		if v[4] != "OK" && status == "" {
			status = v[4]
		}

		if v[5] != "OK" && status == "" {
			status = "Partial solution"
		}

		if v[4] == "OK" && v[5] == "OK" {
			ok++
			status = "OK"
		}
	}

	postingTime := time.Now()

	score := ok * 100 / len(results)

	result := Solution {
		Number: maxNumber,
		Problem: formResult.Problem,
		User: userCredential.Username,
		Posting: postingTime,
		Time: maxTime,
		Memory: maxMemory,
		Status: status,
		Score: score,
	}

	solutionsCollection.InsertOne(context.TODO(), result)

	if score > maxScore {
		update := bson.D {
			{Key: "$set", Value: bson.D {
				{Key: "problem", Value: formResult.Problem},
				{Key: "user", Value: userCredential.Username},
				{Key: "score", Value: score},
			}},
		}

		filter = bson.D{
			{Key: "problem", Value: formResult.Problem},
			{Key: "user", Value: userCredential.Username},
		}

		scoresCollection := client.Database("Judex").Collection("scores")

		cur, err := scoresCollection.Find(context.TODO(), filter)
		scores := 0
		for cur.Next(context.TODO()) {
			var elem Problem
			err := cur.Decode(&elem)

			if err != nil {
				log.Println(err)
				return
			}
			scores++
		}
		if err := cur.Err(); err != nil {
			log.Println(err)
			return
		}
	
		cur.Close(context.TODO())
		
		if err != nil || scores == 0 {
			scoresCollection.InsertOne(context.TODO(), bson.D {
				{Key: "problem", Value: formResult.Problem},
				{Key: "user", Value: userCredential.Username},
				{Key: "score", Value: score},
			})
		}
		scoresCollection.UpdateOne(context.TODO(), filter, update)
	}


	newProtocol, err := os.OpenFile(filepath.Join(solutionsPath, strconv.Itoa(maxNumber), strconv.Itoa(maxNumber) + ".prot"), os.O_RDWR | os.O_CREATE, 0666)

	if err != nil {
		log.Println(err)
		return
	}

	defer newProtocol.Close()

	msg = []byte(FileReader(filepath.Join(tmpPath, strconv.Itoa(maxNumber), strconv.Itoa(maxNumber) + ".prot")))

	//fmt.Println(msg)

	newProtocol.Write(msg)
	
	http.Redirect(w, r, "/solutions/" + strconv.Itoa(maxNumber), http.StatusSeeOther)
}
