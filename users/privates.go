package users

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

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
		c, err := user.(models.User).ToJson(false)
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
func execAuth(r *http.Request) (string, error) {

	var params paramsAuth
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		log.Println(err)
		return "", err
	}
	err = json.Unmarshal(b, &params)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// (*w).Header().Set("Content-Type", "application/json")
	found, err := ds.Find(params.Email)
	if err != nil {
		return "", ErrAuthFail
	}

	// found.
	user := found.GetItem()
	err = bcrypt.CompareHashAndPassword([]byte(user.(models.User).PwHash), []byte(params.Pw))
	if err != nil {
		// pwhash does not match.
		return "", ErrAuthFail
	} else {
		// pwhash matches
		userJsonString, err := user.(models.User).ToJson(false)
		if err != nil {
			return "", err
		}

		return userJsonString, nil
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

	err = ds.InsertNode(u, u.Email)
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

// handler to get a specific user
func handleGetRequest(w *http.ResponseWriter, r *http.Request) {

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
		rtn, _ := user.(models.User).ToJson(true)
		fmt.Fprintln(*w, rtn)
		return
	}

}

// function to handle the post request
func handlePostRequest(w *http.ResponseWriter, r *http.Request, body []byte) {

	// parse the json content into a struct
	// for easier handling
	var paramsAdd paramsAdd
	err := json.Unmarshal(body, &paramsAdd)
	if err != nil {
		http.Error(*w, err.Error(), http.StatusInternalServerError)
	}

	// ok. struct is ready.

	// check if email value is empty
	if isEmptyString(paramsAdd.Email) {
		http.Error(*w, "email is a required attribute", http.StatusBadRequest)
		return
	}

	// check if first name value is empty
	if isEmptyString(paramsAdd.Name.First) {
		http.Error(*w, "name.first is a required attribute", http.StatusBadRequest)
		return
	}

	// check if last name value is empty
	if isEmptyString(paramsAdd.Name.Last) {
		http.Error(*w, "name.last is a required attribute", http.StatusBadRequest)
		return
	}

	// check if password value is empty
	if isEmptyString(paramsAdd.Password) {
		http.Error(*w, "password is a required attribute", http.StatusBadRequest)
		return
	}

	// check if email value is in a proper email format (e.g. joe.jet@motel168.com)
	if !isValidEmailFormat(paramsAdd.Email) {
		http.Error(*w, "email is not a valid format", http.StatusBadRequest)
		return
	}

	// check if roles value is nil (or empty)
	if isEmptyStringSlice(paramsAdd.Roles) {
		http.Error(*w, "roles is a required attribute and must not be empty", http.StatusBadRequest)
		return
	}

	if areValidRoles(paramsAdd.Roles) {
		http.Error(*w, "roles must contain valid values", http.StatusBadRequest)
		return
	}

	// check if the user already existed.
	existed := existed(paramsAdd.Email)
	if existed {
		// user email already existed
		fmt.Fprintf(*w, `{"ok": false, "msg": "%s", "data": {}}`, ErrUserExisted)
		return
	} else {
		// user email is new
		new, err := add(paramsAdd)

		if err != nil {
			fmt.Fprintln(*w, `{"ok": false, "msg": "fail to add user", "data": {}}`)
		}
		fmt.Fprintf(*w, `{"ok": true, "msg": "user added successfully", "data": %s}`, new)
		return
	}
}

// function to handle the put request
func handlePutRequest(w *http.ResponseWriter, r *http.Request, body []byte) {

	var paramsUpdate paramsUpdate
	err := json.Unmarshal(body, &paramsUpdate)
	if err != nil {
		http.Error(*w, err.Error(), http.StatusInternalServerError)
	}
	updated, err := update(paramsUpdate)
	if err != nil {
		rtn := `{
			"ok": false,
			"msg": "fail to find user data to update",
			"data": {}
		}`
		fmt.Fprintln(*w, rtn)
		return
	}

	// update successfully.
	fmt.Fprintf(*w, `{
		"ok": true,
		"msg": "successfully updated user data",
		"data": %s
	}
	`, updated)
	//
}

func handleDeleteRequest(w *http.ResponseWriter, r *http.Request, body []byte) {

	var paramsRemove paramsRemove
	err := json.Unmarshal(body, &paramsRemove)
	if err != nil {
		http.Error(*w, err.Error(), http.StatusInternalServerError)
	}
	err = remove(paramsRemove.Email)
	if err != nil {
		rtn := `{
			"ok" : false,
			"msg" : "user not found",
			"data" : {}	
		}`
		fmt.Fprintln(*w, rtn)
		return
	}
	rtn := `{
		"ok" : true,
		"msg" : "user removed successfully",
		"data" : {}
	}`
	fmt.Fprintln(*w, rtn)
}

// function to get the content of the request body.
func getBody(w *http.ResponseWriter, r *http.Request) []byte {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		log.Println(err)
		http.Error(*w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	// ok.
	return body
}
