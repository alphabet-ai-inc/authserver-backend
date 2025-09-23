package dbrepo_test

import (
	"authserver-backend/internal/dbrepo"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

// this file contains mocking tests for the methods in the dbrepo package.
// In the future, we can use a mocking library like testify to create a more robust mock.
// For now, we will implement the methods manually.
func TestConnectToDB(t *testing.T) {
	mockRepo := dbrepo.MockDBRepo{}
	dsn := "user=postgres, password=postgres, dbname=postgres, sslmode=disable"
	// Set up expectations
	mockRepo.On("ConnectToDB", dsn).Return((*sql.DB)(nil), nil)
	// Use the mock in a function that relies on DatabaseRepo
	db, err := mockRepo.ConnectToDB(dsn)
	// Log the database connection
	t.Log("db: ", db)
	// Log the error
	t.Log("err: ", err)

	// Assertions
	assert.Nil(t, db)
	mockRepo.AssertExpectations(t)
}

func TestConnection(t *testing.T) {
	mockRepo := dbrepo.MockDBRepo{}
	// Set up expectations
	mockRepo.On("Connection").Return((*sql.DB)(nil), nil)

	// Use the mock in a function that relies on DatabaseRepo
	conn, err := mockRepo.Connection()
	// Log the database connection
	t.Log("conn: ", conn)
	// Log the error
	t.Log("err: ", err)

	// Assertions
	assert.Nil(t, conn)
	mockRepo.AssertExpectations(t)
}
