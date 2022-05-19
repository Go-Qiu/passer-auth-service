package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-qiu/passer-auth-service/data/models"
	"github.com/go-qiu/passer-auth-service/jwt"
	"github.com/go-qiu/passer-auth-service/users"
	"github.com/joho/godotenv"
)

var (
	ErrAuthFail                error = errors.New("[API-Users]: authentication failure")
	ErrNotAllowedRequestMethod error = errors.New("[API-Users]: requst method is not allowed for this endpoint")
	ErrUserExisted             error = errors.New("[API-Users]: user already existed")
)

// Auth is a http handler for the 'POST' request to authenticate the user credentials, passed in via the request body.
func (a *application) Auth(w http.ResponseWriter, r *http.Request) {

	// get .env values
	err := godotenv.Load()
	if err != nil {
		errString := "[JWT]: fail to load .env"
		a.errorLog.Println(errString)
		a.clientError(w, http.StatusInternalServerError, errString)
		return
	}
	JWT_ISSUER := os.Getenv("JWT_ISSUER")

	JWT_EXP_MINUTES, err := strconv.Atoi(os.Getenv("JWT_EXP_MINUTES"))
	if err != nil {
		errString := "[JWT]: fail to load .env"
		a.errorLog.Println(errString)
		a.clientError(w, http.StatusInternalServerError, errString)
		return
	}

	// Only allow a 'POST' requst to continue.
	if r.Method != http.MethodPost {

		// not a 'POST' request
		msg := fmt.Sprintf("[AUTH]: request method, '%s' is not allowed for this api endpoint", r.Method)
		// http.Error(w, msg, http.StatusForbidden)
		a.clientError(w, http.StatusBadRequest, msg)
		return
	}

	// ok. it is a 'POST' request.
	if r.Header.Get("Content-Type") == "application/json" {

		// json data is in the request body.

		// let the requestor know that the response will be a JSON.
		w.Header().Set("Content-Type", "application/json")

		// execute the authentication.
		outcome, err := execAuth(a.dataStore, r)
		if err != nil {
			if err != ErrAuthFail {

				a.serverError(w, err)
				return

			} else {

				// auth failure
				msg := fmt.Sprintf(`{
					"ok": false,
					"msg": "[AUTH]: %s",
					"data": {}
				}`, ErrAuthFail.Error())

				a.clientError(w, http.StatusForbidden, msg)
				return

			}
		}

		// ok. authentication passed.

		var foundUser models.User
		err = json.Unmarshal([]byte(outcome), &foundUser)
		if err != nil {
			a.serverError(w, err)
			return

		}

		exp := time.Now().Add(time.Minute * time.Duration(JWT_EXP_MINUTES)).UnixMilli()
		pl := jwt.JWTPayload{
			Id:       foundUser.Email,
			Name:     foundUser.Name.First + " " + foundUser.Name.Last,
			Roles:    foundUser.Roles,
			IsActive: foundUser.IsActive,
			Iss:      JWT_ISSUER,
			Exp:      exp,
		}

		var token string
		token, err = generateJWT(pl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		msg := fmt.Sprintf(`{
			"ok" : true,
			"msg" : "[AUTH]: authentication successful",
			"data" : {
				"token" : "%s"
			}
		}`, token)

		w.Header().Set("Token", token)
		fmt.Fprintln(w, msg)
		return
	}
	//
}

// Users method to direct request to operate on users related data to
// the appropriate user data operations handler.
func (a *application) Users(w http.ResponseWriter, r *http.Request) {

	users.Handler(w, r, a.dataStore)
}
