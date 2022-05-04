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
	}

	content := ""
	count := 1
	for accounts.GetSize() > 0 {
		user, _ := accounts.Pop()
		if count == 1 {
			// first user data point
			content += "{"
		} else {
			// more user data point
			content += ", {"
		}

		content += fmt.Sprintf("\"id\" : \"%s\", ", user.Id)
		content += fmt.Sprintf("\"email\" : \"%s\", ", user.Email)
		content += fmt.Sprintf("\"isActive\" : %v, ", user.IsActive)
		content += fmt.Sprintf("\"name\" : { \"first\" :\" %s\", \"last\" : \"%s\"}, ", user.Name.First, user.Name.Last)

		// handle roles attribure
		if len(user.Roles) == 0 {
			// empty roles attribute
			content += fmt.Sprintln("\"roles\" : []")
		} else {
			// non-empty roles attribute
			roles := ""
			count := 1
			for _, r := range user.Roles {
				if count == 1 {
					roles += fmt.Sprintf("\"%s\"", r)
				} else {
					roles += fmt.Sprintf(", \"%s\"", r)
				}
				count++
			}
			content += fmt.Sprintf("\"roles\" : [%s]", roles)
		}
		content += fmt.Sprintln("}")
		count++
		// fmt.Println(user)
		// fmt.Println()
	}
	rtn := fmt.Sprintf("[%s]", content)
	fmt.Fprintln(w, rtn)

}

// handler to get a specific user
func Get(w http.ResponseWriter, r *http.Request) {

}