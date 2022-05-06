package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-qiu/passer-auth-service/data"
	"github.com/go-qiu/passer-auth-service/helpers"
)

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

// type outUser struct {
// 	Id       string `json:"id"`
// 	Email    string `json:"email"`
// 	Name     name   `json:"name"`
// 	IsActive bool   `json:"isActive"`
// 	Roles    bool   `json:"roles"`
// }

var (
	ErrAuthFail                error = errors.New("[API-Users]: authentication failure")
	ErrNotAllowedRequestMethod error = errors.New("[API-Users]: requst method is not allowed for this endpoint")
	ErrUserExisted             error = errors.New("[API-Users]: user already existed")
)

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

// handler for all users data related handling
func Handler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		// 'GET' request
		get(&w, r)
		//
	} else {
		// not a 'GET' request

		// request content-type is "application/json"
		if r.Header.Get("Content-Type") == "application/json" {

			// get the json content in the request body
			body, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()

			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if r.Method == http.MethodPost {
				// 'POST' request --> add

				// parse the json content into a struct
				// for easier handling
				var paramsAdd paramsAdd
				err = json.Unmarshal(body, &paramsAdd)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				// ok. struct is ready.

				// check if the user already existed.
				existed := existed(paramsAdd.Email)
				if existed {
					// user email already existed
					fmt.Fprintf(w, `{"ok": false, "msg": "%s", "data": {}}`, ErrUserExisted)
					return
				} else {
					// user email is new
					new, err := add(paramsAdd)

					if err != nil {
						fmt.Fprintln(w, `{"ok": false, "msg": "fail to add user", "data": {}}`)
					}
					fmt.Fprintf(w, `{"ok": true, "msg": "user added successfully", "data": %s}`, new)
					return
				}

				//
			} else if r.Method == http.MethodPut {

				// 'PUT' request --> update
				var paramsUpdate paramsUpdate
				err = json.Unmarshal(body, &paramsUpdate)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				updated, err := update(paramsUpdate)
				if err != nil {
					rtn := `{
						"ok": false,
						"msg": "fail to find user data to update",
						"data": {}
					}`
					fmt.Fprintln(w, rtn)
					return
				}

				// update successfully.
				fmt.Fprintf(w, `{
					"ok": true,
					"msg": "successfully updated user data",
					"data": %s
				}
				`, updated)
				return
				//
			} else if r.Method == http.MethodDelete {

				// 'DELETE' request --> remove
				var paramsRemove paramsRemove
				err = json.Unmarshal(body, &paramsRemove)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				err = remove(paramsRemove.Email)
				if err != nil {
					rtn := `{
						"ok" : false,
						"msg" : "user not found",
						"data" : {}	
					}`
					fmt.Fprintln(w, rtn)
					return
				}
				rtn := `{
					"ok" : true,
					"msg" : "user removed successfully",
					"data" : {}
				}`
				fmt.Fprintln(w, rtn)
				return
				//
			} else {
				// not any of the above methods.
				msg := fmt.Sprintf("Request method, '%s' is not allowed for this api endpoint.\n", r.Method)
				http.Error(w, msg, http.StatusForbidden)
				return
			}
			//
		}
		//
	}

}

// handler to get a specific user
func get(w *http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {

		// not a 'GET' request
		msg := fmt.Sprintf("Request method, '%s' is not allowed for this api endpoint.", r.Method)
		http.Error(*w, msg, http.StatusForbidden)
		return
	}

	// ok. it is a 'GET' request.
	// get the params passed in via the url
	params := r.URL.Query()

	if len(params) == 0 {
		// no parameters were passed in via the url.
		// list all users.
		getAll(w, r)
	}

	if len(params) > 0 && len(strings.TrimSpace(params.Get("id"))) != 0 {
		// id was passed in via the url
		(*w).Header().Set("Content-Type", "application/json")
		// get the user data point that matches the id
		found, err := ds.Find(params.Get("id"))
		if err != nil {
			log.Println(err)
			http.Error(*w, err.Error(), http.StatusInternalServerError)
			return
		}

		// ok.
		user := found.GetItem()
		rtn, _ := user.ToJson(true)
		fmt.Fprintln(*w, rtn)
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
