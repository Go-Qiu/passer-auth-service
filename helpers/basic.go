package helpers

import (
	"github.com/go-qiu/passer-auth-service/data/models"
	"golang.org/x/crypto/bcrypt"
)

// Preload create the user data points for loading into the in-memory data store.
// This is facilitate development and testing.
// Only the admin account will be retained when the project is ready for deployment.
func Preload() ([]models.User, error) {

	pw := "Testing.12345"
	pwHash, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.MinCost)

	pwAdmin := "pA22er.54321"
	pwAdminHash, _ := bcrypt.GenerateFromPassword([]byte(pwAdmin), bcrypt.MinCost)

	users := []models.User{}

	uAdmin := models.User{
		Id:    "admin@passer.com",
		Email: "admin@passer.com", PwHash: string(pwAdminHash),
		Name:  models.Name{First: "Admin", Last: "PASSER"},
		Roles: []string{"ADMIN"}, IsActive: true}

	users = append(users, uAdmin)

	uMerchant01 := models.User{
		Id:    "xy.lim@bestbuy.com",
		Email: "xy.lim@bestbuy.com", PwHash: string(pwHash),
		Name:  models.Name{First: "X. Y.", Last: "Lim"},
		Roles: []string{"MERCHANT"}, IsActive: true}

	users = append(users, uMerchant01)

	uMerchant02 := models.User{
		Id:    "azi.abdu@bismi.com",
		Email: "azi.abdu@bismi.com", PwHash: string(pwHash),
		Name:  models.Name{First: "X. Y.", Last: "Lim"},
		Roles: []string{"MERCHANT"}, IsActive: true}

	users = append(users, uMerchant02)

	uUser01 := models.User{
		Id:    "jimmy.dean@gmail.com",
		Email: "jimmy.dean@gmail.com", PwHash: string(pwHash),
		Name:  models.Name{First: "Jimmy", Last: "Dean"},
		Roles: []string{"CONSUMER", "AGENT"}, IsActive: true}

	users = append(users, uUser01)

	uUser02 := models.User{
		Id:    "jolin.lim@gmail.com",
		Email: "jolin.lim@gmail.com", PwHash: string(pwHash),
		Name:  models.Name{First: "Jolin", Last: "Lim"},
		Roles: []string{"CONSUMER", "AGENT"}, IsActive: true}

	users = append(users, uUser02)

	uAgent01 := models.User{
		Id:    "joe.jet@gmail.com",
		Email: "joe.jet@gmail.com", PwHash: string(pwHash),
		Name:  models.Name{First: "Joe", Last: "Jet"},
		Roles: []string{"AGENT"}, IsActive: true}

	users = append(users, uAgent01)

	uAgent02 := models.User{
		Id:    "jacky.chuang@gmail.com",
		Email: "jacky.chuang@gmail.com", PwHash: string(pwHash),
		Name:  models.Name{First: "Jacky", Last: "Chuang"},
		Roles: []string{"AGENT"}, IsActive: true}

	users = append(users, uAgent02)

	uAgent03 := models.User{
		Id:    "alex.tao@gmail.com",
		Email: "alex.tao@gmail.com", PwHash: string(pwHash),
		Name:  models.Name{First: "Alex", Last: "Tao"},
		Roles: []string{"AGENT"}, IsActive: true}
	users = append(users, uAgent03)

	return users, nil
}
