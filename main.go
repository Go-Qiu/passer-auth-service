package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-qiu/passer-auth-service/data"
	"github.com/go-qiu/passer-auth-service/helpers"
)

var ds data.DataStore = *data.New()

func main() {

	addr := "localhost:8081"

	// Simulate a data pull of PASSER Locker Station
	// specific Parcel Job records from the HQ Data Center.
	// The records are inserted into the local data store.

	userList, err := helpers.Preload()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	for _, u := range userList {
		ds.InsertNode(u, u.Email)
	}

	// declare custom loggers
	infoLog := log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)

	// declare and instantiate a web application
	app := &application{
		errorLog:  errorLog,
		infoLog:   infoLog,
		dataStore: &ds,
	}

	// declare and instantiate a custom http server
	srv := &http.Server{
		Addr:     addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("HTTPS Server started and listening on https://%s ...", addr)
	err = srv.ListenAndServeTLS("./ssl/cert.pem", "./ssl/key.pem")
	if err != nil {
		errorLog.Fatal(err)
	}
}
