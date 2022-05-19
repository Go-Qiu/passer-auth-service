package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-qiu/passer-auth-service/jwt"
	"github.com/joho/godotenv"
)

func validateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// declare custom loggers
		// infoLog := log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)
		errorLog := log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)

		// get .env values
		err := godotenv.Load()
		if err != nil {
			errString := "[JWT]: fail to load .env"
			errorLog.Println(errString)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		JWT_SECRET_KEY := os.Getenv("JWT_SECRET_KEY")

		// get the jwt from the request header.
		authorization := r.Header.Get("Authorization")
		if strings.TrimSpace(authorization) == "" {

			errString := "[Middleware]: no token found"
			errorLog.Println(errString)

			msg := fmt.Sprintf(`{
				"ok": false,
				"msg": "%s",
				"data": {}
			}`, errString)
			http.Error(w, msg, http.StatusForbidden)
			return
		}

		token := strings.TrimPrefix(authorization, "Bearer ")
		if strings.TrimSpace(token) == "" {
			// empty token
			errString := "[Middleware]: no token found"
			errorLog.Println(errString)

			http.Error(w, errString, http.StatusForbidden)
			return
		}

		// ok.
		// jwt validation logic here.
		ok, err := jwt.Verify(token, JWT_SECRET_KEY)
		if err != nil {

			errorLog.Println(err.Error())
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		if !ok {

			errString := "[JWT]: fail to validate token"
			w.Header().Set("Content-Type", "application/json")

			http.Error(w, errString, http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
