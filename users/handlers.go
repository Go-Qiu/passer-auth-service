package users

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-qiu/passer-auth-service/data"
	"github.com/go-qiu/passer-auth-service/helpers"
)

var ErrAuthFail error = errors.New("[API-Users]: authentication failure")

// var userList []models.User
var ds data.DataStore = *data.New()

func init() {

	userList, err := helpers.Preload()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	for _, u := range userList {
		ds.InsertNode(u)
	}

}

// handler to add a user
func Add(w http.ResponseWriter, r *http.Request) {

}

// handler to update a user
func Update(w http.ResponseWriter, r *http.Request) {

}

// handler to remove a user
func Remove(w http.ResponseWriter, r *http.Request) {

}

// handler to get a specific user
func Get(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {

		// not a 'GET' request
		msg := fmt.Sprintf("Request method, '%s' is not allowed for this api endpoint.", r.Method)
		http.Error(w, msg, http.StatusForbidden)
		return
	}

	// ok. it is a 'GET' request.
	// get the params passed in via the url
	params := r.URL.Query()

	if len(params) == 0 {
		// no parameters were passed in via the url.
		// list all users.
		getAll(&w, r)
	}

	if len(params) > 0 && len(strings.TrimSpace(params.Get("id"))) != 0 {
		// id was passed in via the url
		w.Header().Set("Content-Type", "application/json")
		// get the user data point that matches the id
		found, err := ds.Find(params.Get("id"))
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// ok.
		user := found.GetItem()
		rtn, _ := user.ToJson(true)
		fmt.Fprintln(w, rtn)
		return
	}

}

// handler to authenticate a user
func Auth(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		// not a 'POST' request
		msg := fmt.Sprintf("Request method, '%s' is not allowed for this api endpoint.", r.Method)
		http.Error(w, msg, http.StatusForbidden)
		return
	}

	// ok. it is a 'POST' request.
	if r.Header.Get("Content-Type") == "application/json" {
		// json data is in the request body.
		// get the json passed in.
		w.Header().Set("Content-Type", "application/json")
		err := execAuth(r)
		if err != nil {
			if err != ErrAuthFail {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else {
				// auth failure
				fmt.Fprintln(w, `{"ok": false}`)
				return
			}
		}

		// ok.
		fmt.Fprintln(w, `{"ok": true}`)
	}
}
