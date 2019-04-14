package main

import (
    "log"               // Logs
    "net/http"          // Server logic
    "html/template"     // Html usage
    "os"                // OS syscalls

    "context"
    _ "github.com/mongodb/mongo-go-driver/bson"
    "github.com/mongodb/mongo-go-driver/mongo"
    _ "github.com/mongodb/mongo-go-driver/mongo/options"
)


/*
    Structures
*/

// Structure for storing and using the user information
type User struct {
    Username string
    Email string
    Password string
    F_name string
    M_name string
    L_name string
    Company string
    Website string
    Bio string
}


/*
    Global variables
*/

// Html templates storage
var templates *template.Template
// MongoDB client
var client *mongo.Client


/*
    Request handlers
*/

// GET request handler for home page
func index(w http.ResponseWriter, r *http.Request) {
    if err := templates.ExecuteTemplate(w, "index.html", ""); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

// GET request handler for registration page
func registration(w http.ResponseWriter, r *http.Request) {
    if err := templates.ExecuteTemplate(w, "registration.html", ""); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

// POST request handler for registration submit
func reg_submit(w http.ResponseWriter, r *http.Request) {
    /*if err := templates.ExecuteTemplate(w, "index.html", ""); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }*/

    //collection := client.Database("Judex").Collection("users")


    r.ParseForm()
    log.Println(r.Form)
}


/*
    The core
*/

func main() {
    // Creating of a log file
    f, err := os.OpenFile("logs.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    defer f.Close()
    log.SetOutput(f)

    // Enabling the HTML templates
    templates = template.Must(template.ParseFiles("templates/index.html", "templates/registration.html"))

    // Changing URL for "/static/"
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    // Defining handlers for requests
    http.HandleFunc("/", index)
    http.HandleFunc("/registration", registration)
    http.HandleFunc("/reg_submit", reg_submit)

    // Connecting to the database
    client, err := mongo.Connect(context.TODO(), "mongodb://localhost:27017")
    err = client.Ping(context.TODO(), nil)
    if err != nil {
        log.Fatal(err)
    }
    log.Println("Connected to database")

    // Listening the 80/http port and serving clients
    log.Println("Listening")
    log.Fatal(http.ListenAndServe(":80", nil))

    // Database disconnection
    err = client.Disconnect(context.TODO())
    if err != nil {
        log.Fatal(err)
    }
    log.Println("Connection to MongoDB closed.")
}
