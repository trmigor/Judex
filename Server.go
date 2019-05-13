package main

import (
	"html/template" // Html usage
	"log"           // Logs
	"net"           // Server logic
	"net/http"      // Server logic
	"os"            // OS syscalls
	"time"          // Timing

	// Database
	"context"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// User is a structure for storing and using the user information
type User struct {
	Username   string
	Email      string
	Password   string
	FirstName  string
	MiddleName string
	LastName   string
	Company    string
	Website    string
	Bio        string
}

// Credential is a structure for storing credentials
type Credential struct {
	UserIP    net.IP
	Username  string
	EnterTime time.Time
	EndTime   time.Time
}

// CredCur is a credential duration
const CredCur time.Duration = 3 * time.Hour

var (
	// Html templates storage
	templates *template.Template

	// MongoDB client
	client *mongo.Client
)

func main() {
	args := os.Args
	var logPath, templatePath string

	if len(args) >= 3 {
		templatePath = args[2]
	} else {
		templatePath = "/var/www/judex.vdi.mipt.ru"
	}
	if len(args) >= 2 {
		logPath = args[1]
	} else {
		logPath = "/var/cache/Judex"
	}

	// Creating of a log file
	logFile, err := os.OpenFile(logPath+"/logs.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.Println("STARTED")

	// Enabling the HTML templates
	// Slice of templates' names
	templateNames := []string{
		templatePath + "/templates/index.html",
		templatePath + "/templates/registration.html",
		templatePath + "/templates/sign_in.html",
		templatePath + "/templates/error404.html",
		templatePath + "/templates/home.html",
		templatePath + "/templates/profile.html",
	}

	templates = template.Must(template.ParseFiles(templateNames...))

	// Changing URL for "/static/"
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Defining handlers for requests
	http.HandleFunc("/", Index)
	http.HandleFunc("/registration", Registration)
	http.HandleFunc("/reg_submit", RegSubmit)
	http.HandleFunc("/sign_in", SignIn)
	http.HandleFunc("/sign_in_submit", SignInSubmit)
	http.HandleFunc("/home", Home)
	http.HandleFunc("/sign_out", SignOut)
	http.HandleFunc("/profile", Profile)

	// Connecting to the database
	client, err = mongo.Connect(context.TODO(), "mongodb://localhost:27017")
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
	log.Println("STOPPED")
}
