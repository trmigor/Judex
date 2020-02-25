package main

import (
	"log"      // Logs
	"net/http" // Server logic

	// Database
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

// ChangeSubmit handles POST request for profile change submit
func ChangeSubmit(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Println("/change_submit: Cannot parse update form")
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	// Checking if user wants to access /change_submit without filling a form
	if len(r.Form) == 0 {
		log.Println("/change_submit: Blank form")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	formResult := User{
		Username:   r.Form["username"][0],
		Email:      r.Form["email"][0],
		Password:   r.Form["password"][0],
		FirstName:  r.Form["f_name"][0],
		MiddleName: r.Form["m_name"][0],
		LastName:   r.Form["l_name"][0],
		Company:    r.Form["company"][0],
		Website:    r.Form["website"][0],
		Bio:        r.Form["bio"][0],
	}

	usersCollection := client.Database("Judex").Collection("users")

	// Checking the uniqueness of the email
	var findResult User
	filter := bson.D{{Key: "email", Value: formResult.Email}}
	err = usersCollection.FindOne(context.TODO(), filter).Decode(&findResult)

	if err == nil && findResult.Username != formResult.Username {
		log.Println("/change_submit: Email", formResult.Email, "is already used")
		ErrorHandler(w, r, DoubledEmail)
		return
	}

	// Saving received information
	filter = bson.D{{Key: "username", Value: formResult.Username}}
	update := bson.D {
		{Key: "$set", Value: bson.D {
			{Key: "email", Value: formResult.Email},
			{Key: "password", Value: formResult.Password},
			{Key: "firstname", Value: formResult.FirstName},
			{Key: "middlename", Value: formResult.MiddleName},
			{Key: "lastname", Value: formResult.LastName},
			{Key: "company", Value: formResult.Company},
			{Key: "website", Value: formResult.Website},
			{Key: "bio", Value: formResult.Bio},
		}},
	}
	_, err = usersCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Println("/change_submit: users: Cannot update information")
		ErrorHandler(w, r, UpdateError)
		return
	}

	log.Println("Updated. Username:", formResult.Username)

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
	return
}