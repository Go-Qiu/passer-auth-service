package main

import (
	"log"

	"github.com/go-qiu/passer-auth-service/data"
)

// JWTPayload struct is for holding the data used in generating the second segment (i.e. payload) of the JWT string.
type JWTPayload struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Roles    []string `json:"roles"`
	IsActive bool     `json:"isActive"`
	Iss      string   `json:"iss"`
	Exp      int64    `json:"exp"`
}

// JWTHeader struct is for holding the data used in generating the first segment of the JWT string.
type JWTHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

// application struct is for facilitating the implementation of the dependencies injection model.
type application struct {
	errorLog  *log.Logger
	infoLog   *log.Logger
	dataStore *data.DataStore
}
