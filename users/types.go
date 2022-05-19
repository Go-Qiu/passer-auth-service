package users

// paramsAuth struct is for holding the auth request body content.
type paramsAuth struct {
	Email string `json:"email"`
	Pw    string `json:"pw"`
}
