package models

type User struct {
	Id       string
	Email    string
	PwHash   string
	Name     Name
	IsActive bool
	Roles    []string
}

type Name struct {
	First string
	Last  string
}

func (u User) ParseToJson() string {
	return "{}"
}
