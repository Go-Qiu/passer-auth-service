package users

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-qiu/passer-auth-service/data/stack"
	"golang.org/x/crypto/bcrypt"
)

type paramsAuth struct {
	Email string `json:"email"`
	Pw    string `json:"pw"`
}

// function to list all the users (without the pwhash)
func getAll(w *http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		// not a 'GET' request
		msg := fmt.Sprintf("Request method, '%s' is not allowed for this api endpoint.", r.Method)
		http.Error(*w, msg, http.StatusForbidden)
		return
	}

	// ok.  is a 'GET' request.

	// set the Header.Content-Type to "application/json" to
	// ensure the proper return of the outcome in json format
	// to the response
	(*w).Header().Set("Content-Type", "application/json")

	// instantiate a stack to cache the nodes
	accounts := stack.New()

	err := ds.ListAllNodes(&accounts, false)
	if err != nil {
		log.Println(err)
		http.Error(*w, err.Error(), http.StatusInternalServerError)
	}

	content := ""
	count := 1
	hasFailed := false

	for accounts.GetSize() > 0 {
		user, _ := accounts.Pop()
		c, err := user.ToJson(false)
		if err != nil {
			log.Println(err)
			hasFailed = true
			http.Error(*w, err.Error(), http.StatusInternalServerError)
			break
		}
		if count == 1 {

			// is first user data point
			content += c
		} else {

			// is subsequent user data point
			content += fmt.Sprintf(", %s", c)
		}
		count++
	}

	if hasFailed {
		return
	} else {
		rtn := fmt.Sprintf("[%s]", content)
		fmt.Fprintln(*w, rtn)
	}
}

// function to execute the authentication check
func execAuth(r *http.Request) error {

	var params paramsAuth
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		log.Println(err)
		// http.Error(*w, err.Error(), http.StatusInternalServerError)
		return err
	}
	err = json.Unmarshal(b, &params)
	if err != nil {
		log.Println(err)
		// http.Error(*w, err.Error(), http.StatusInternalServerError)
		return err
	}

	// (*w).Header().Set("Content-Type", "application/json")
	found, err := ds.Find(params.Email)
	if err != nil {
		// fmt.Fprintln(*w, `{"ok": false}`)
		return err
	}

	// found.
	user := found.GetItem()
	err = bcrypt.CompareHashAndPassword([]byte(user.PwHash), []byte(params.Pw))
	if err != nil {
		// pwhash does not match.
		// fmt.Fprintln(*w, `{"ok": false}`)
		return ErrAuthFail
	} else {
		// pwhash matches
		// fmt.Fprintln(*w, `{"ok": true}`)
		return nil
	}
}
