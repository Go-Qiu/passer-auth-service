package models

import "fmt"

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

	content := "{\n\t"
	content += fmt.Sprintf("\"id\" : \"%s\", \n\t", u.Id)
	content += fmt.Sprintf("\"email\" : \"%s\", \n\t", u.Email)
	content += fmt.Sprintf("\"isActive\" : %v, \n\t", u.IsActive)
	content += fmt.Sprintf("\"name\" : { \n\t\"first\" :\" %s\", \n\t\"last\" : \"%s\"\n\t}, ", u.Name.First, u.Name.Last)

	// handle roles attribure
	if len(u.Roles) == 0 {
		// empty roles attribute
		content += fmt.Sprintln("\"roles\" : []")
	} else {
		// non-empty roles attribute
		roles := ""
		count := 1
		for _, r := range u.Roles {
			if count == 1 {
				roles += fmt.Sprintf("\"%s\"\n\t", r)
			} else {
				roles += fmt.Sprintf(", \"%s\"\n\t", r)
			}
			count++
		}
		content += fmt.Sprintf("\"roles\" : [%s]", roles)
	}
	content += fmt.Sprintln("\n\t}")
	//
	return content
}
