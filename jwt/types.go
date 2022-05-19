package jwt

// JWTPayload is the struct for holding the data used in generating the second segment of the JWT string.
type JWTPayload struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Roles    []string `json:"roles"`
	IsActive bool     `json:"isActive"`
	Iss      string   `json:"iss"`
	Exp      int64    `json:"exp"`
}

// JWTHeader is the struct for holding the data used in generating the first segment of the JWT string.
type JWTHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}
