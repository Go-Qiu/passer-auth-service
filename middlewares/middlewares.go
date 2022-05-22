package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-qiu/passer-auth-service/jwt"
	"github.com/joho/godotenv"
)

// ValidateJWT is a middleware that will check for the presence of a 'Token' attribute in the request header.
// It will permit the request to continue its flow to the secureed api endpoint if the 'Token' is present and valid.
// A valid 'Token' must satisfy the following:
// - the signature segment of the 'Token' must be consistent when this middleware signs the content of the Header and Payload segments (of the 'Token') with the secret key;
// - the 'exp' attribute in the Payload (encoded in Base64 format) has not come to pass.
func ValidateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// declare custom loggers
		// infoLog := log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)
		errorLog := log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)

		// get .env values
		err := godotenv.Load()
		if err != nil {
			errString := "[JWT]: fail to load .env"
			errorLog.Println(errString)
			// http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			msg := fmt.Sprintf(`{
				"ok": false,
				"msg": "%s",
				"data": {}
			}`, errString)
			fmt.Fprintln(w, msg)
			return
		}
		JWT_SECRET_KEY := os.Getenv("JWT_SECRET_KEY")

		// get the jwt from the request header.
		authorization := r.Header.Get("Authorization")
		if strings.TrimSpace(authorization) == "" {

			errString := "[Middleware]: no token found"
			errorLog.Println(errString)

			// http.Error(w, msg, http.StatusForbidden)

			w.WriteHeader(http.StatusForbidden)
			w.Header().Set("Content-Type", "application/json")
			msg := fmt.Sprintf(`{
				"ok": false,
				"msg": "%s",
				"data": {}
			}`, errString)
			fmt.Fprintln(w, msg)
			return
		}

		token := strings.TrimPrefix(authorization, "Bearer ")
		if strings.TrimSpace(token) == "" {
			// empty token
			errString := "[Middleware]: no token found"
			errorLog.Println(errString)
			// http.Error(w, errString, http.StatusForbidden)

			w.WriteHeader(http.StatusForbidden)
			w.Header().Set("Content-Type", "application/json")
			msg := fmt.Sprintf(`{
				"ok": false,
				"msg": "%s",
				"data": {}
			}`, errString)
			fmt.Fprintln(w, msg)
			return
		}

		// ok.
		// jwt validation logic here.
		ok, err := jwt.Verify(token, JWT_SECRET_KEY)
		if err != nil {

			errorLog.Println(err.Error())
			// http.Error(w, err.Error(), http.StatusForbidden)

			w.WriteHeader(http.StatusForbidden)
			w.Header().Set("Content-Type", "application/json")
			msg := fmt.Sprintf(`{
				"ok": false,
				"msg": "%s",
				"data": {}
			}`, err)
			fmt.Fprintln(w, msg)
			return
		}

		if !ok {

			errString := "[JWT]: fail to validate token"
			// http.Error(w, errString, http.StatusForbidden)

			w.WriteHeader(http.StatusForbidden)
			w.Header().Set("Content-Type", "application/json")
			msg := fmt.Sprintf(`{
				"ok": false,
				"msg": "%s",
				"data": {}
			}`, errString)
			fmt.Fprintln(w, msg)
			return
		}

		// direct the request to the next handler.
		next.ServeHTTP(w, r)
	})
}
