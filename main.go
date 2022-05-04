package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-qiu/passer-auth-service/data"
	"github.com/go-qiu/passer-auth-service/data/stack"
	"github.com/go-qiu/passer-auth-service/helpers"
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

func init() {

	users, err := helpers.Preload()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	ds := data.New()
	for _, u := range users {
		ds.InsertNode(u)
	}

	accounts := stack.New()
	err = ds.ListAllNodes(&accounts, false)
	if err != nil {
		log.Println(err)
	}

	for accounts.GetSize() > 0 {
		user, _ := accounts.Pop()
		fmt.Println(user)
		fmt.Println()
	}

}

func main() {

	addr := "localhost:8081"
	http.HandleFunc("/signup", signUp)
	http.HandleFunc("/hash", handleHash)
	http.HandleFunc("/verify", verifyHash)

	log.Printf("HTTP Server is started and listening at %s ...\n", addr)
	log.Fatalln(http.ListenAndServe(addr, nil))
}
