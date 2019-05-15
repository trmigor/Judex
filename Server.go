package main

import (
	"fmt"			// I/O formatting
	"io"			// EOF
	"strings"		// Strings splitting
	"html/template" // Html usage
	"log"           // Logs
	"net"           // Server logic
	"net/http"      // Server logic
	"os"            // OS syscalls
	"time"          // Timing
	"os/signal"		// Signal handling

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

	// Server configure
	logPath, templatePath, staticPath, emailPatterns string
)

func main() {
	// Automatic configure
	config, err := os.OpenFile("server.config", os.O_RDONLY, 0400)

	if err != nil {
		log.Fatal(err)
	}

	for {
		var configInput string
		_, err := fmt.Fscanln(config, &configInput)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		setting := strings.Split(configInput, "=")
		value := strings.Split(setting[1], "\"")
		switch setting[0] {
		case "logPath":
			logPath = value[1]
		case "templatePath":
			templatePath = value[1]
		case "staticPath":
			staticPath = value[1]
		case "emailPatterns":
			emailPatterns = value[1]
		}
	}

	// Manual configure
	args := os.Args

	if len(args) >= 5 {
		emailPatterns = args[4]
	}
	if len(args) >= 4 {
		staticPath = args[3]
	}
	if len(args) >= 3 {
		templatePath = args[2]
	}
	if len(args) >= 2 {
		logPath = args[1]
	}

	// Creating of a log file
	logFile, err := os.OpenFile(logPath+"/logs.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(logFile)

	log.Println("STARTED")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Kill, os.Interrupt)
	go Stop(sigs, logFile)

	// Enabling the HTML templates
	// Slice of templates' names
	templateNames := []string{
		templatePath + "/index.html",
		templatePath + "/registration.html",
		templatePath + "/sign_in.html",
		templatePath + "/error404.html",
		templatePath + "/home.html",
		templatePath + "/profile.html",
		templatePath + "/problems.html",
	}

	templates = template.Must(template.ParseFiles(templateNames...))

	// Changing URL for "/static/"
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))

	// Defining handlers for requests
	http.HandleFunc("/", Index)
	http.HandleFunc("/registration", Registration)
	http.HandleFunc("/reg_submit", RegSubmit)
	http.HandleFunc("/sign_in", SignIn)
	http.HandleFunc("/sign_in_submit", SignInSubmit)
	http.HandleFunc("/home", Home)
	http.HandleFunc("/sign_out", SignOut)
	http.HandleFunc("/profile", Profile)
	http.HandleFunc("/change_submit", ChangeSubmit)
	http.HandleFunc("/problems", Problems)

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

	Stop(nil, logFile)
}
