package main

import (
	"fmt"			// I/O formatting
	"io"			// EOF
	"strings"		// Strings split
	"html/template" // Html usage
	"log"           // Logs
	"net"           // Server logic
	"net/http"      // Server logic
	"os"            // OS syscalls
	"time"          // Timing
	"os/signal"		// Signal handling
	"path/filepath"	// Filepath join

	"github.com/dpapathanasiou/go-recaptcha"

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
	logPath, templatePath, staticPath, emailPatterns, problemsPath, solutionsPath string
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
		case "problemsPath":
			problemsPath = value[1]
		case "solutionsPath":
			solutionsPath = value[1]
		}
	}

	// Manual configure
	args := os.Args

	if len(args) >= 7 {
		solutionsPath = args[6]
	}
	if len(args) >= 6 {
		problemsPath = args[5]
	}
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
	logFile, err := os.OpenFile(filepath.Join(logPath, "logs.log"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(logFile)

	log.Println("STARTED")

	// Catching signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Kill, os.Interrupt)
	go Stop(sigs, logFile)

	// Initializing reCAPTCHA
	reCaptchaFile, err := os.OpenFile("reCAPTCHA_private", os.O_RDONLY, 0400)
	if err != nil {
		log.Fatal(err)
	}
	var privateKey string
	fmt.Fscanln(reCaptchaFile, &privateKey)
	recaptcha.Init(privateKey)

	// Enabling the HTML templates
	// Slice of templates' names
	templateNames := []string{
		filepath.Join(templatePath, "index.html"),
		filepath.Join(templatePath, "registration.html"),
		filepath.Join(templatePath, "sign_in.html"),
		filepath.Join(templatePath, "error404.html"),
		filepath.Join(templatePath, "home.html"),
		filepath.Join(templatePath, "profile.html"),
		filepath.Join(templatePath, "problems.html"),
		filepath.Join(templatePath, "problem.html"),
		filepath.Join(templatePath, "post.html"),
		filepath.Join(templatePath, "solution.html"),
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
	http.HandleFunc("/problem/", SingleProblem)
	http.HandleFunc("/post/", Post)

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
