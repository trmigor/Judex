package main

import (
	"fmt"		// Output formatting
	"net/http" 	// Server logic
)

const (
	// DoubledUsername is an error for a case, when user tries to create an account with already existing username
	DoubledUsername int = 1100
	// DoubledEmail is an error for a case, when user tries to create an account with already used email
	DoubledEmail int = 1101
	// InsertError is an error for a case, when database insert fails
	InsertError int = 1102

	// NoUsername is an error for a case, when user tries to sign in with non existing username
	NoUsername int = 1200
	// WrongPassword is an error for a case, when user tries to sign in with a wrong password
	WrongPassword int = 1201
)

/*
   Error handler

   Custom errors:
       11xx: Database Error:
		   1100 Doubled Username
		   1101 Doubled Email
           1102 Insert Error
       12xx: Client Error:
           1200 No Username
           1201 Wrong Password
*/

// ErrorHandler handles network and custom errors
func ErrorHandler(w http.ResponseWriter, r *http.Request, status int) {
	switch status {
	case http.StatusNotFound:
		if err := templates.ExecuteTemplate(w, "error404.html", ""); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case http.StatusInternalServerError:
		w.WriteHeader(status)
		fmt.Fprint(w, "Something is wrong with our server")
	case DoubledUsername:
		page := RegistrationPage{
			UsernameError: "block",
			EmailError: "none",
		}
	
		if err := templates.ExecuteTemplate(w, "registration.html", page); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case DoubledEmail:
		page := RegistrationPage{
			UsernameError: "none",
			EmailError: "block",
		}
	
		if err := templates.ExecuteTemplate(w, "registration.html", page); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case InsertError:
		fmt.Fprint(w, "We cannot save your information")
	case NoUsername:
		page := SignInPage{
			Username:      "",
			UsernameError: "block",
			PasswordError: "none",
		}
	
		if err := templates.ExecuteTemplate(w, "sign_in.html", page); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case WrongPassword:
		page := SignInPage{
			Username:      "",
			UsernameError: "none",
			PasswordError: "block",
		}
	
		if err := templates.ExecuteTemplate(w, "sign_in.html", page); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}