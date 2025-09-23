package models

import (
	"errors"

	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system with various attributes.
// JSON tags are included for serialization/deserialization.
// This is the main user model used throughout the application,
// including authentication for every app in the ecosystem
// and user management in the admin interface.
type User struct {
	ID             int    `json:"id"`
	UserName       string `json:"username"`
	Password       string `json:"password"`
	Code           string `json:"code"`
	Active         bool   `json:"active"`
	LastLogin      int    `json:"last_login"`
	LastSession    string `json:"last_session"`
	Blocked        bool   `json:"blocked"`
	Tries          int    `json:"tries"`
	LastTry        int64  `json:"last_try"`
	Email          string `json:"email"`
	ProfileId      int    `json:"profile_id"`
	GroupId        int    `json:"group_id"`
	DbsAuth        int    `json:"dbs_auth"`
	ActivationTime int64  `json:"activation_time"`
	LastAction     string `json:"last_action"`
	LastApp        int    `json:"last_app"`
	LastDb         int    `json:"last_db1"`
	Lan            string `json:"lan"`
	CompanyId      int    `json:"company_id"`
	Created        int64  `json:"created"`
	Updated        int64  `json:"updated"`
}

// PassMatch defines the interface for password matching.
type PassMatch interface {
	PasswordMatches(string) (bool, error)
}

// MockPWCheck is a mock implementation of the PassMatch interface for testing purposes.
type MockPWCheck struct {
	mock.Mock
	PassMatch
}

// PasswordMatches mocks the password matching function.
func (m *MockPWCheck) PasswordMatches(plainText string, uPasswd string) (bool, error) {
	args := m.Called(plainText)
	if args.Get(0) == nil {
		return false, args.Error(1)
	}
	return true, args.Error(1)
}

// PasswordMatches checks if the provided plain text password matches the stored hashed password.
func (u *User) PasswordMatches(plainText string, uPasswd string) (bool, error) {
	// Compare the plain text password with the hashed password
	err := bcrypt.CompareHashAndPassword([]byte(uPasswd), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

// This seems not to be the rigth place for these functions.
// Maybe a utils package would be better.
// Also could be moved to a service layer to the authservice package
// or to an authhandler or auth package inside the auth folder.
// It also will contain the autentication handlers and services now in api package.
