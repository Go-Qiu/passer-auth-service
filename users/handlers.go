package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-qiu/passer-auth-service/data"
	"golang.org/x/crypto/bcrypt"
)

var mapUsers = map[string]user{}

type user struct {
	Username string `json:"username"`
	PwHash   string `json:"pwhash"`
	First    string `json:"first"`
	Last     string `json:"last"`
}

type name struct {
	First string `json:"first"`
	Last  string `json:"last"`
}

type paramsAdd struct {
	Email    string   `json:"email"`
	Name     name     `json:"name"`
	Password string   `json:"password"`
	IsActive bool     `json:"isActive"`
	Roles    []string `json:"roles"`
}

type paramsRemove struct {
	Email string `json:"email"`
}

type updateableFields struct {
	IsActive bool     `json:"isActive"`
	Name     name     `json:"name"`
	Roles    []string `json:"roles"`
}
type paramsUpdate struct {
	Email   string           `json:"email"`
	Updates updateableFields `json:"updates"`
}

// JWTPayload is the struct for holding the data used in generating the second segment of the JWT string.
type JWTPayload struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Roles    []string `json:"roles"`
	IsActive bool     `json:"isActive"`
	Iss      string   `json:"iss"`
	Exp      int64    `json:"exp"`
}

// JWTHeader is the struct for holding the data used in generating the first segment of the JWT string.
type JWTHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

var (
	ErrAuthFail                error = errors.New("[API-Users]: authentication failure")
	ErrNotAllowedRequestMethod error = errors.New("[API-Users]: requst method is not allowed for this endpoint")
	ErrUserExisted             error = errors.New("[API-Users]: user already existed")
)

// var userList []models.User
// var ds data.DataStore = *data.New()

// func init() {

// 	userList, err := helpers.Preload()
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	for _, u := range userList {
// 		ds.InsertNode(u, u.Email)
// 	}
// }

// Handler handles all users data related data operations.
func Handler(w http.ResponseWriter, r *http.Request, ds *data.DataStore) {

	// set the response header, "Content-Type" to "application/json".
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		// 'GET' request
		handleGetRequest(&w, r, ds)
	} else {
		// not a 'GET' request

		// request content-type is "application/json"
		if r.Header.Get("Content-Type") == "application/json" {

			// get the json content in the request body
			body := getBody(&w, r)

			// ok. body content (in []byte format)
			// is ready for further handling of POST, PUT, DELETE
			// processing.
			if r.Method == http.MethodPost {

				// 'POST' request --> add
				handlePostRequest(&w, r, ds, body)
			} else if r.Method == http.MethodPut {

				// 'PUT' request --> update
				handlePutRequest(&w, r, ds, body)
			} else if r.Method == http.MethodDelete {

				// 'DELETE' request --> remove
				handleDeleteRequest(&w, r, ds, body)
			} else {
				// not any of the above methods.
				msg := fmt.Sprintf("Request method, '%s' is not allowed for this api endpoint.\n", r.Method)
				http.Error(w, msg, http.StatusForbidden)
				return
			}
		}
	}

}

// Signup handles request to add a user. WIP.
func SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("usrename")
		password := r.FormValue("password")
		confirmation := r.FormValue("confirmation")
		first := r.FormValue("first")
		last := r.FormValue("last")

		// exceptions handling
		if isEmptyString(username) {
			http.Error(w, "Username cannot be empty.", http.StatusForbidden)
			return
		}
		if isEmptyString(password) {
			http.Error(w, "Password cannot be empty.", http.StatusForbidden)
			return
		}
		if isEmptyString(confirmation) {
			http.Error(w, "Password cannot be empty.", http.StatusForbidden)
			return
		}
		if confirmation == password {
			http.Error(w, "2 password entries are not the same.", http.StatusForbidden)
			return
		}
		if isEmptyString(first) {
			http.Error(w, "First Name cannot be empty.", http.StatusForbidden)
			return
		}
		if isEmptyString(last) {
			http.Error(w, "Last Name cannot be empty.", http.StatusForbidden)
			return
		}
		if _, ok := mapUsers[username]; ok {
			http.Error(w, "Username is already taken.", http.StatusForbidden)
			return
		}

		// ok. ready.
		pwhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		if err != nil {
			http.Error(w, "Internal Server Error.", http.StatusInternalServerError)
			return
		}
		newUser := user{Username: username, PwHash: string(pwhash), First: first, Last: last}
		mapUsers[username] = newUser

		// redirect to post sign-up page
		// http.Redirect(w, r, "/", http.StatusSeeOther)
		w.Header().Set("Content-Type", "application/json")
		outcome, err := json.Marshal(newUser)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fmt.Fprintln(w, outcome)
		return
	}
	// tpl.ExecuteTemplate(w, "signup.html", user{})
}
