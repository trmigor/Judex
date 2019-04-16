package main

import (
    "fmt"
    "log"               // Logs
    "net"               // Server logic
    "net/http"          // Server logic
    "html/template"     // Html usage
    "os"                // OS syscalls
    "time"              // Timing

    // Database
    "context"
    "github.com/mongodb/mongo-go-driver/bson"
    "github.com/mongodb/mongo-go-driver/mongo"
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

// Structure for storing credentials
type Credential struct {
    User_ip net.IP
    Username string
    Enter_time time.Time
    End_time time.Time
}

// Structure for checking sign in information
type Sign_in struct {
    Username string
    Password string
}


/*
    Constants
*/

// Credential duration
//const CRED_DUR time.Duration = 10800*1000*1000*1000
const CRED_DUR time.Duration = 30*1000*1000*1000


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
    // 404 error handle
    if r.URL.Path != "/" {
        errorHandler(w, r, http.StatusNotFound)
        return
    }

    // Checking credentials
    // If user is already logged in, redirect to /home
    ip, _, err := net.SplitHostPort(r.RemoteAddr)
    
    if err != nil {
        log.Println("/: Cannot discover user's IP")
        errorHandler(w, r, http.StatusInternalServerError)
        return
    }

    user_credential := Credential {
        User_ip:        net.ParseIP(ip),
    }

    if user_credential.User_ip == nil {
        log.Println("/: Cannot discover user's IP")
        errorHandler(w, r, http.StatusInternalServerError)
        return
    }

    credentials_collection := client.Database("Judex").Collection("credentials")

    filter := bson.D{{"user_ip", user_credential.User_ip}}
    err = credentials_collection.FindOne(context.TODO(), filter).Decode(&user_credential)

    if err == nil {
        if time.Now().Before(user_credential.End_time) {
            http.Redirect(w, r, "/home", http.StatusSeeOther)
            return
        } else {
            credentials_collection.DeleteMany(context.TODO(), filter)
        }
    }

    if err := templates.ExecuteTemplate(w, "index.html", ""); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

// GET request handler for registration page
func registration(w http.ResponseWriter, r *http.Request) {
    // Checking credentials
    // If user is already logged in, redirect to /home
    ip, _, err := net.SplitHostPort(r.RemoteAddr)
    
    if err != nil {
        log.Println("/registration: Cannot discover user's IP")
        errorHandler(w, r, http.StatusInternalServerError)
        return
    }

    user_credential := Credential {
        User_ip:    net.ParseIP(ip),
    }

    if user_credential.User_ip == nil {
        log.Println("/registration: Cannot discover user's IP")
        errorHandler(w, r, http.StatusInternalServerError)
        return
    }

    credentials_collection := client.Database("Judex").Collection("credentials")

    filter := bson.D{{"user_ip", user_credential.User_ip}}
    err = credentials_collection.FindOne(context.TODO(), filter).Decode(&user_credential)

    if err == nil {
        if time.Now().Before(user_credential.End_time) {
            http.Redirect(w, r, "/home", http.StatusSeeOther)
            return
        } else {
            credentials_collection.DeleteMany(context.TODO(), filter)
        }
    }

    if err := templates.ExecuteTemplate(w, "registration.html", ""); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

// GET request handler for sign in page
func sign_in(w http.ResponseWriter, r *http.Request) {
    // Checking credentials
    // If user is already logged in, redirect to /home
    ip, _, err := net.SplitHostPort(r.RemoteAddr)
    
    if err != nil {
        log.Println("/sign_in: Cannot discover user's IP")
        errorHandler(w, r, http.StatusInternalServerError)
        return
    }

    user_credential := Credential {
        User_ip:        net.ParseIP(ip),
    }

    if user_credential.User_ip == nil {
        log.Println("/sign_in: Cannot discover user's IP")
        errorHandler(w, r, http.StatusInternalServerError)
        return
    }

    credentials_collection := client.Database("Judex").Collection("credentials")

    filter := bson.D{{"user_ip", user_credential.User_ip}}
    err = credentials_collection.FindOne(context.TODO(), filter).Decode(&user_credential)

    if err == nil {
        if time.Now().Before(user_credential.End_time) {
            http.Redirect(w, r, "/home", http.StatusSeeOther)
            return
        } else {
            credentials_collection.DeleteMany(context.TODO(), filter)
        }
    }

    if err := templates.ExecuteTemplate(w, "sign_in.html", ""); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

// POST request handler for registration submit
func reg_submit(w http.ResponseWriter, r *http.Request) {

    err := r.ParseForm()
    
    if err != nil {
        log.Println("/reg_submit: Cannot parse registration form")
        errorHandler(w, r, http.StatusInternalServerError)
        return
    }

    // Checking if the user wants to access /reg_submit without filling a form
    if len(r.Form) == 0 {
        log.Println("/reg_submit: Blank form")
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    form_result := User {
        Username:   r.Form["username"][0],
        Email:      r.Form["email"][0],
        Password:   r.Form["password"][0],
        F_name:     r.Form["f_name"][0],
        M_name:     r.Form["m_name"][0],
        L_name:     r.Form["l_name"][0],
        Company:    r.Form["company"][0],
        Website:    r.Form["website"][0],
        Bio:        r.Form["bio"][0],
    }

    users_collection := client.Database("Judex").Collection("users")

    // Checking the uniqueness of the username
    var find_result User
    filter := bson.D{{"username", form_result.Username}}
    err = users_collection.FindOne(context.TODO(), filter).Decode(&find_result)
    
    if err == nil {
        log.Println("/reg_submit: Username", form_result.Username, "is not unique")
        errorHandler(w, r, 1100)
        return
    }

    // Saving received information
    _, err = users_collection.InsertOne(context.TODO(), form_result)
    if err != nil {
        log.Println("/reg_submit: users: Cannot insert information")
        errorHandler(w, r, 1101)
        return
    }

    log.Println("Registered. Username:", form_result.Username)

    // Creating a credential
    ip, _, err := net.SplitHostPort(r.RemoteAddr)
    
    if err != nil {
        log.Println("/reg_submit: Cannot discover user's IP")
        errorHandler(w, r, http.StatusInternalServerError)
        return
    }

    user_credential := Credential {
        User_ip:        net.ParseIP(ip),
        Username:       form_result.Username,
        Enter_time:     time.Now(),
        End_time:       time.Now().Add(CRED_DUR),
    }

    if user_credential.User_ip == nil {
        log.Println("/reg_submit: Cannot discover user's IP")
        errorHandler(w, r, http.StatusInternalServerError)
        return
    }

    credentials_collection := client.Database("Judex").Collection("credentials")

    filter = bson.D{{"user_ip", user_credential.User_ip}}
    credentials_collection.DeleteMany(context.TODO(), filter)

    _, err = credentials_collection.InsertOne(context.TODO(), user_credential)

    if err != nil {
        log.Println("/reg_submit: credentials: Cannot insert information")
        errorHandler(w, r, 1101)
        return
    }

    log.Println("Logged in. Username:", form_result.Username)

    http.Redirect(w, r, "/home", http.StatusSeeOther)
    return
}

func sign_in_submit(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    
    if err != nil {
        log.Println("/sign_in_submit: Cannot parse sign in form")
        errorHandler(w, r, http.StatusInternalServerError)
        return
    }

    // Checking if the user wants to access /sign_in_submit without filling a form
    if len(r.Form) == 0 {
        log.Println("/sign_in_submit: Blank form")
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    form_result := Sign_in {
        Username:   r.Form["username"][0],
        Password:   r.Form["password"][0],
    }

    // Checking the information
    users_collection := client.Database("Judex").Collection("users")

    var find_result User
    filter := bson.D{{"username", form_result.Username}}
    err = users_collection.FindOne(context.TODO(), filter).Decode(&find_result)

    if err != nil {
        log.Println("/sign_in_submit: No such username:", form_result.Username)
        errorHandler(w, r, 1200)
        return
    }

    if form_result.Password != find_result.Password {
        log.Println("/sign_in_submit: Password do not match")
        errorHandler(w, r, 1201)
        return
    }

    // Creating a credential
    ip, _, err := net.SplitHostPort(r.RemoteAddr)
    
    if err != nil {
        log.Println("/sign_in_submit: Cannot discover user's IP")
        errorHandler(w, r, http.StatusInternalServerError)
        return
    }

    user_credential := Credential {
        User_ip:        net.ParseIP(ip),
        Username:       form_result.Username,
        Enter_time:     time.Now(),
        End_time:       time.Now().Add(CRED_DUR),
    }

    if user_credential.User_ip == nil {
        log.Println("/sign_in_submit: Cannot discover user's IP")
        errorHandler(w, r, http.StatusInternalServerError)
        return
    }

    credentials_collection := client.Database("Judex").Collection("credentials")

    filter = bson.D{{"user_ip", user_credential.User_ip}}
    credentials_collection.DeleteMany(context.TODO(), filter)

    _, err = credentials_collection.InsertOne(context.TODO(), user_credential)

    if err != nil {
        log.Println("/sign_in_submit: credentials: Cannot insert information")
        errorHandler(w, r, 1101)
        return
    }

    log.Println("Logged in. Username:", form_result.Username)

    http.Redirect(w, r, "/home", http.StatusSeeOther)
    return
}

func home(w http.ResponseWriter, r *http.Request) {
    // Checking credentials
    // If user is not logged in, redirect to /
    ip, _, err := net.SplitHostPort(r.RemoteAddr)
    
    if err != nil {
        log.Println("/home: Cannot discover user's IP")
        errorHandler(w, r, http.StatusInternalServerError)
        return
    }

    user_credential := Credential {
        User_ip:        net.ParseIP(ip),
    }

    if user_credential.User_ip == nil {
        log.Println("/home: Cannot discover user's IP")
        errorHandler(w, r, http.StatusInternalServerError)
        return
    }

    credentials_collection := client.Database("Judex").Collection("credentials")

    filter := bson.D{{"user_ip", user_credential.User_ip}}
    err = credentials_collection.FindOne(context.TODO(), filter).Decode(&user_credential)

    if err == nil {
        if time.Now().After(user_credential.End_time) {
            http.Redirect(w, r, "/", http.StatusSeeOther)
            credentials_collection.DeleteMany(context.TODO(), filter)
            return
        }
    } else {
        http.Redirect(w, r, "/", http.StatusSeeOther)
    }
}

/*
    Error handler

    Custom errors:
        11xx: Database Error:
            1100 Doubled Username
            1101 Insert Error
        12xx: Client Error:
            1200 No Username
            1201 Wrong Password
*/

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
    switch status {
        case http.StatusNotFound:
            if err := templates.ExecuteTemplate(w, "error404.html", ""); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
            }
        case http.StatusInternalServerError:
            fmt.Fprint(w, "Something is wrong with our server")
        case 1100:
            fmt.Fprint(w, "Username is not unique")
        case 1101:
            fmt.Fprint(w, "We cannot save your information")
        case 1200:
            fmt.Fprint(w, "No such user")
        case 1201:
            fmt.Fprint(w, "Wrong password")
    }
}


/*
    Helping functions
*/

/*
    The core
*/

func main() {
    // Creating of a log file
    log_file, err := os.OpenFile("logs.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    defer log_file.Close()
    log.SetOutput(log_file)
    log.Println("STARTED");

    // Enabling the HTML templates
    templates = template.Must(template.ParseFiles("templates/index.html", "templates/registration.html", "templates/sign_in.html", "templates/error404.html"))

    // Changing URL for "/static/"
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    // Defining handlers for requests
    http.HandleFunc("/", index)
    http.HandleFunc("/registration", registration)
    http.HandleFunc("/reg_submit", reg_submit)
    http.HandleFunc("/sign_in", sign_in)
    http.HandleFunc("/sign_in_submit", sign_in_submit)
    http.HandleFunc("/home", home)

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
