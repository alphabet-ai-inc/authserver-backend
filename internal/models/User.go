package models

import (
	"errors"

	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

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

//	type User struct {
//		ID       int    `json:"id"`
//		Email    string `json:"email"`
//		Password string `json:"password"`
//	}
type PassMatch interface {
	PasswordMatches(plainText string) (bool, error)
}

type MockPassMatch struct {
	mock.Mock
}

func (m *MockPassMatch) PasswordMatches(plainText string) (bool, error) {
	args := m.Called(plainText)
	if args.Get(0) == nil {
		return false, args.Error(1)
	}
	return true, args.Error(1)
	// return args.Bool(0), args.Error(1) user, err := app.DB.GetUserByEmail(requestPayload.Email)

}

func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
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
