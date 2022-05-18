package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-qiu/passer-auth-service/data"
	"github.com/go-qiu/passer-auth-service/data/models"
	"github.com/go-qiu/passer-auth-service/jwt"
	"golang.org/x/crypto/bcrypt"
)

// paramsAuth type struct is used for unmarshalling
// the json send via the request body sent to the
// the endpoint, '/auth'.
type paramsAuth struct {
	Email string `json:"email"`
	Pw    string `json:"pw"`
}

// function to execute the authentication check
func execAuth(ds *data.DataStore, r *http.Request) (string, error) {

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

// generateJWT will generate a JWT using the header and payload passed in.
func generateJWT(payload JWTPayload) (string, error) {

	// secret key to use "P@ss3r.54321"
	header := `{
		"alg": "SHA512",
		"typ" : "JWT"
	}`

	// convert payload data to json string
	pl, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	token := jwt.Generate(header, string(pl), "P@ss3r.54321")

	return token, nil
}
