package main

import (
	"log"
	"net/http"

	"github.com/go-qiu/passer-auth-service/users"
)

// var tpl *template.Template = template.Must(template.ParseGlob("templates/*"))

func main() {

	addr := "localhost:8081"

	mux := http.NewServeMux()
	mux.HandleFunc("/signup", users.SignUp)
	mux.HandleFunc("/users", users.Handler)
	mux.HandleFunc("/auth", users.Auth)
	// mux.HandleFunc("/hash", handleHash)

	log.Printf("HTTPS Server is started and listening at %s ...\n", addr)
	// log.Fatalln(http.ListenAndServe(addr, mux))
	log.Fatalln(http.ListenAndServeTLS(addr, "./ssl/cert.pem", "./ssl/key.pem", mux))
}
