package users

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-qiu/passer-auth-service/data"
	"github.com/go-qiu/passer-auth-service/data/stack"
	"github.com/go-qiu/passer-auth-service/helpers"
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

// handler to add a user
func Add(w http.ResponseWriter, r *http.Request) {

}

// handler to update a user
func Update(w http.ResponseWriter, r *http.Request) {

}

// handler to remove a user
func Remove(w http.ResponseWriter, r *http.Request) {

}

// handler to list all the users
func GetAll(w http.ResponseWriter, r *http.Request) {

	// set the Header.Content-Type to "application/json" to
	// ensure the proper return of the outcome in json format
	// to the response
	w.Header().Set("Content-Type", "application/json")

	// instantiate a stack to cache the nodes
	accounts := stack.New()

	err := ds.ListAllNodes(&accounts, false)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	content := ""
	count := 1
	hasFailed := false

	for accounts.GetSize() > 0 {
		user, _ := accounts.Pop()
		if count == 1 {
			// first user data point
			// content += user.ParseToJson()
			c, err := user.ToJson(true)
			if err != nil {
				log.Println(err)
				hasFailed = true
				http.Error(w, err.Error(), http.StatusInternalServerError)
				break
			}
			content += c
		} else {
			// more user data point
			// content += fmt.Sprintf(", %s", user.ParseToJson())
			c, err := user.ToJson(true)
			if err != nil {
				log.Println(err)
				hasFailed = true
				http.Error(w, err.Error(), http.StatusInternalServerError)
				break
			}
			content += fmt.Sprintf(", %s", c)
		}
		count++
		// fmt.Println(user)
		// fmt.Println()
	}

	if hasFailed {
		return
	} else {
		rtn := fmt.Sprintf("[%s]", content)
		fmt.Fprintln(w, rtn)
	}
}

// handler to get a specific user
func Get(w http.ResponseWriter, r *http.Request) {

}
