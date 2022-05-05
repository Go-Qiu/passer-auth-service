package main

import (
	"log"
	"net/http"

	"github.com/go-qiu/passer-auth-service/users"
)

type user struct {
	Username string `json:"username"`
	PwHash   string `json:"pwhash"`
	First    string `json:"first"`
	Last     string `json:"last"`
}

// var tpl *template.Template = template.Must(template.ParseGlob("templates/*"))
var mapUsers = map[string]user{}

// var mapSessions = map[string]string{}

func main() {

	addr := "localhost:8081"
	http.HandleFunc("/signup", signUp)
	http.HandleFunc("/hash", handleHash)
	http.HandleFunc("/users", users.Get)
	http.HandleFunc("/auth", users.Auth)

	log.Printf("HTTP Server is started and listening at %s ...\n", addr)
	log.Fatalln(http.ListenAndServe(addr, nil))
}
