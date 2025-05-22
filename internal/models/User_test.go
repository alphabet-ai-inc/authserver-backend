package models_test

import (
	"backend/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordMatches(t *testing.T) {
	// Create a user with a hashed password
	plaintextPassword := "mysecretpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(plaintextPassword), bcrypt.DefaultCost)

	user := models.User{
		Password: string(hashedPassword),
	}

	// Test with the correct password
	match, err := user.PasswordMatches(plaintextPassword)
	assert.NoError(t, err)
	assert.True(t, match)

	// Test with an incorrect password
	incorrectPassword := "wrongpassword"
	match, err = user.PasswordMatches(incorrectPassword)
	assert.NoError(t, err)
	assert.False(t, match)

	// Test with an empty password
	match, err = user.PasswordMatches("")
	assert.NoError(t, err)
	assert.False(t, match)

	// Test with a nil or blank password stored in the struct if applicable
	user.Password = ""
	match, err = user.PasswordMatches("somepassword")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "hashedSecret too short to be a bcrypted password")
	assert.False(t, match)

	// Check that bcrypt returns a specific error for an invalid password hash
	user.Password = "invalidhash" // Set a clearly invalid hash
	match, err = user.PasswordMatches("mysecretpassword")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "hashedSecret too short to be a bcrypted password")
	assert.False(t, match)
}
