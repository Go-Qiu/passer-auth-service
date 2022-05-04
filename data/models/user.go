package models

import (
	"encoding/json"
	"fmt"
)

type User struct {
	Id       string   `json:"id"`
	Email    string   `json:"email"`
	PwHash   string   `json:"-"`
	Name     Name     `json:"name"`
	IsActive bool     `json:"isActive"`
	Roles    []string `json:"roles"`
}

type Name struct {
	First string `json:"first"`
	Last  string `json:"last"`
}

func (u User) parseAllToJson() string {

	content := "{\n\t"
	content += fmt.Sprintf("\"id\" : \"%s\", \n\t", u.Id)
	content += fmt.Sprintf("\"email\" : \"%s\", \n\t", u.Email)
	content += fmt.Sprintf("\"isActive\" : %v, \n\t", u.IsActive)
	content += fmt.Sprintf("\"pwhash\" : \"%s\", \n\t", u.PwHash)
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

func (u User) ToJson(all bool) (string, error) {
	if all {
		// all attributes
		rtn := u.parseAllToJson()
		return rtn, nil
	} else {
		// leave out pwhash attribute
		rtn, err := json.Marshal(u)
		if err != nil {
			return "", nil
		}
		return string(rtn), nil
	}

}
