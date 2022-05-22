package main

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/go-qiu/passer-auth-service/middlewares"
)

// serverError will log server side errors and send a HTTP Internal Server Error to the requestor.
func (a *application) serverError(w http.ResponseWriter, err error) {
	// log the error on the server side
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	a.errorLog.Println(trace)

	// send an http error response to the requestor.
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError sends a specific http error, with status code and text, to the requestor.
func (a *application) clientError(w http.ResponseWriter, status int, msg ...string) {
	http.Error(w, http.StatusText(status), status)
}

// notFound sends a http error, indicating the Not Found error, to the requestor.
// func (a *application) notFound(w http.ResponseWriter) {
// 	a.clientError(w, http.StatusNotFound)
// }

// routes returns a server mux, containing all the path patterns to handlers mapping.
func (a *application) routes() *http.ServeMux {

	// it is recommended not to  use the default server mux implementation in the http package, in production.
	// recommended to declare a custom server mux, for use in instantiating a http server, in production.
	mux := http.NewServeMux()

	// fixed path patterns
	mux.HandleFunc("/auth", a.Auth)
	mux.Handle("/users", middlewares.ValidateJWT(http.HandlerFunc(a.Users)))
	mux.Handle("/verify", middlewares.ValidateJWT(http.HandlerFunc(a.Verify)))
	return mux
}
