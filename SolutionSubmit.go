package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/trmigor/Judex/testing_packages/compileandrun"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// SolutionForm holds information received from a request
type SolutionForm struct {
	Problem    int
	Compiler   string
	SourceText string
}

// SolutionSubmit handles POST request of sending a solution, checks it and forms a protocol
func SolutionSubmit(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)

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
		//Compiler: r.Form["compiler"][0],
	}

	solutionsCollection := client.Database("Judex").Collection("solutions")

	cur, err := solutionsCollection.Find(context.TODO(), bson.D{})

	if err != nil {
		log.Println(err)
		return
	}

	maxNumber := 0

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

	f1, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		log.Println(err)
		return
	}

	defer f1.Close()

	io.Copy(f1, file)

	f2, err := os.OpenFile(filepath.Join(tmpPath, strconv.Itoa(maxNumber), "sol_"+strconv.Itoa(maxNumber)+"."+extension), os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		log.Println(err)
		return
	}

	defer f2.Close()

	io.Copy(f2, file)

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
		TestsNumber: len(files),
		RunLimits: compileandrun.Limits{
			TL:  TL * time.Second,
			ML:  ML * 1024 * 1024,
			RTL: 10 * time.Second,
		},
	}
	c.Compile()

}
