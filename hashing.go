package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-qiu/passer-auth-service/helpers"
	"golang.org/x/crypto/bcrypt"
)

type contentToBeHashed struct {
	Input string `json:"input"`
}

type contentToBeVerified struct {
	Content string `json:"content"`
	Hash    string `json:"hash"`
}

type outcome struct {
	Ok   bool        `json:"ok"`
	Data interface{} `json:"data"`
}

func signUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("usrename")
		password := r.FormValue("password")
		confirmation := r.FormValue("confirmation")
		first := r.FormValue("first")
		last := r.FormValue("last")

		// exceptions handling
		if helpers.IsEmpty(username) {
			http.Error(w, "Username cannot be empty.", http.StatusForbidden)
			return
		}
		if helpers.IsEmpty(password) {
			http.Error(w, "Password cannot be empty.", http.StatusForbidden)
			return
		}
		if helpers.IsEmpty(confirmation) {
			http.Error(w, "Password cannot be empty.", http.StatusForbidden)
			return
		}
		if confirmation == password {
			http.Error(w, "2 password entries are not the same.", http.StatusForbidden)
			return
		}
		if helpers.IsEmpty(first) {
			http.Error(w, "First Name cannot be empty.", http.StatusForbidden)
			return
		}
		if helpers.IsEmpty(last) {
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

// function to generate a hash for the string content passed in
// via the request body.
func handleHash(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		// request method is 'POST'
		if r.Header.Get("Content-Type") == "application/json" {
			// json data is in the request body
			// get the json passed in via the request body
			handleHashingBodyData(&w, r)
		}

		if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
			handleHashingFormData(&w, r)
		}
	}
}

// test function
func verifyHash(w http.ResponseWriter, r *http.Request) {

	// get json passed in through the request body
	var tobeVerified contentToBeVerified
	err := json.NewDecoder(r.Body).Decode(&tobeVerified)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// verify the authenticity of thecontent against the hash
	err = bcrypt.CompareHashAndPassword([]byte(tobeVerified.Hash), []byte(tobeVerified.Content))
	w.Header().Set("Content-Type", "application/json")

	var rtn outcome

	if err == nil {
		// ok.  is matching.
		rtn = outcome{Ok: true, Data: nil}
	} else {
		// content and hash does not match
		rtn = outcome{Ok: false, Data: nil}
	}
	outp, err := json.Marshal(rtn)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(outp))
}

// function to handle json data send via the request body
func handleHashingBodyData(w *http.ResponseWriter, r *http.Request) {
	var tobeHashed contentToBeHashed
	err := json.NewDecoder(r.Body).Decode(&tobeHashed)
	if err != nil {
		// exceptions handling
		log.Println(err)
		http.Error(*w, err.Error(), http.StatusBadRequest)
	}

	p, err := bcrypt.GenerateFromPassword([]byte(tobeHashed.Input), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		http.Error(*w, err.Error(), http.StatusInternalServerError)
		return
	}

	rtn := outcome{Ok: true, Data: string(p)}
	outp, err := json.Marshal(rtn)
	if err != nil {
		log.Println(err)
		http.Error(*w, err.Error(), http.StatusInternalServerError)
	}
	// send the outocme json to the requesting device
	(*w).Header().Set("Content-Type", "application/json")
	fmt.Fprintln(*w, string(outp))
	//
}

// function to handle data send via a form post
// through a request
func handleHashingFormData(w *http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		// 'POST' request method
		// get the json passed in via the request body

		err := r.ParseForm()
		if err != nil {
			log.Println(err)
			http.Error(*w, err.Error(), http.StatusInternalServerError)
		}
		// p, err := hash(r.FormValue("input"))
		input := r.FormValue("input")
		p, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.MinCost)
		if err != nil {
			log.Println(err)
			http.Error(*w, err.Error(), http.StatusInternalServerError)
		}

		rtn := outcome{Ok: true, Data: string(p)}
		outp, err := json.Marshal(rtn)
		if err != nil {
			// exceptions handling
			rtn := outcome{Ok: false, Data: nil}
			outp, _ := json.Marshal(rtn)
			fmt.Fprintln(*w, string(outp))
			return
		}
		(*w).Header().Set("Content-Type", "application/json")
		fmt.Fprintln(*w, string(outp))
		return
	}
}
