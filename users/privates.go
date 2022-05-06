package users

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-qiu/passer-auth-service/data/models"
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
		return err
	}
	err = json.Unmarshal(b, &params)
	if err != nil {
		log.Println(err)
		return err
	}

	// (*w).Header().Set("Content-Type", "application/json")
	found, err := ds.Find(params.Email)
	if err != nil {
		return ErrAuthFail
	}

	// found.
	user := found.GetItem()
	err = bcrypt.CompareHashAndPassword([]byte(user.PwHash), []byte(params.Pw))
	if err != nil {
		// pwhash does not match.
		return ErrAuthFail
	} else {
		// pwhash matches
		return nil
	}
}

// function to add a user
func add(p paramsAdd) (string, error) {

	var u models.User

	u.Id = p.Email
	u.Email = p.Email
	u.Name.First = p.Name.First
	u.Name.Last = p.Name.Last
	u.IsActive = p.IsActive
	u.Roles = p.Roles

	pwhash, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	u.PwHash = string(pwhash)

	err = ds.InsertNode(u)
	if err != nil {
		return "", err
	}

	// get the new user added from the in-memory data store
	n, err := ds.Find(p.Email)
	if err != nil {
		return "", err
	}
	new := n.GetItem()

	rtn, err := json.Marshal(new)
	if err != nil {
		return "", nil
	}
	return string(rtn), nil
}

// function to update a user
func update(p paramsUpdate) (string, error) {

	updates := models.User{}
	updates.Id = p.Email
	updates.Email = p.Email
	updates.IsActive = p.Updates.IsActive
	updates.Name.First = p.Updates.Name.First
	updates.Name.Last = p.Updates.Name.Last
	updates.Roles = p.Updates.Roles

	updated, err := ds.Update(p.Email, updates)
	if err != nil {
		return "{}", err
	}

	u, err := json.Marshal(updated)
	if err != nil {
		return "{}", err
	}
	return string(u), nil
}

// function to remove a user
func remove(email string) error {

	err := ds.Remove(email)
	if err != nil {
		return err
	}

	return nil
}

// function to check (by email) if a data point (i.e. user)
// existed in the in-memory data store.
func existed(email string) bool {
	found, err := ds.Find(email)
	if err != nil && found == nil {
		return false
	}

	// found user data point
	return true
}
