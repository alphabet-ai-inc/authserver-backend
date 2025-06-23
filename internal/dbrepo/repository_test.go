package dbrepo_test

import (
	"backend/internal/dbrepo"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

// func TestAllApps(t *testing.T) {
// 	mockRepo := new(MockDBRepo)

// 	// Create a slice of expected apps
// 	expectedApps := []*models.ThisApp{
// 		{ID: 1, NewApp: models.NewApp{Name: "App 1"}},
// 		{ID: 1, NewApp: models.NewApp{Name: "App 2"}},
// 	}

// 	// Set up expectations
// 	mockRepo.On("AllApps").Return(expectedApps, nil)

// 	// Use the mock in a function that relies on DatabaseRepo
// 	apps, err := mockRepo.AllApps()

// 	// Log each expected app
// 	for _, app := range expectedApps {
// 		t.Log("expapps: ", app.String())
// 	}
// 	t.Log(err)

// 	// Assertions
// 	assert.Nil(t, err)
// 	assert.Equal(t, expectedApps, apps)
// 	mockRepo.AssertExpectations(t)
// }

// func TestThisAppForEdit(t *testing.T) {
// 	mockRepo := new(MockDBRepo)
// 	id := 1
// 	expectedApp := &models.ThisApp{ID: id}

// 	// Set up expectations
// 	mockRepo.On("ThisAppForEdit", id).Return(expectedApp, nil)

// 	// Use the mock in a function that relies on DatabaseRepo
// 	app, err := mockRepo.ThisAppForEdit(id)
// 	// Assertions
// 	assert.Nil(t, err)
// 	assert.Equal(t, expectedApp, app)
// 	mockRepo.AssertExpectations(t)
// }
// func TestThisApp(t *testing.T) {
// 	mockRepo := new(MockDBRepo)
// 	id := 1
// 	expectedApp := &models.ThisApp{ID: id}

// 	// Set up expectations
// 	mockRepo.On("ThisApp", id).Return(expectedApp, nil)

// 	// Use the mock in a function that relies on DatabaseRepo
// 	app, err := mockRepo.ThisApp(id)
// 	// Assertions
// 	assert.Nil(t, err)
// 	assert.Equal(t, expectedApp, app)
// 	mockRepo.AssertExpectations(t)
// }
// func TestInsertApp(t *testing.T) {
// 	mockRepo := new(MockDBRepo)
// 	newapp := models.NewApp{Name: "App 1"}
// 	id := 1

// 	// Set up expectations
// 	mockRepo.On("InsertApp", newapp).Return(id, nil)

// 	// Use the mock in a function that relies on DatabaseRepo
// 	newid, err := mockRepo.InsertApp(newapp)
// 	t.Log("newid: ", newid)

// 	// Assertions
// 	assert.Nil(t, err)
// 	assert.Equal(t, id, newid)
// 	mockRepo.AssertExpectations(t)
// }
// func TestUpdateApp(t *testing.T) {
// 	mockRepo := new(MockDBRepo)
// 	thisapp := models.ThisApp{ID: 1, NewApp: models.NewApp{Name: "App 1"}}

// 	// Set up expectations
// 	mockRepo.On("UpdateApp", thisapp).Return(nil)

// 	// Use the mock in a function that relies on DatabaseRepo
// 	err := mockRepo.UpdateApp(thisapp)

// 	// Assertions
// 	assert.Nil(t, err)
// 	mockRepo.AssertExpectations(t)
// }
// func TestDeleteApp(t *testing.T) {
// 	mockRepo := new(MockDBRepo)
// 	id := 1

// 	// Set up expectations
// 	mockRepo.On("DeleteApp", id).Return(nil)

// 	// Use the mock in a function that relies on DatabaseRepo
// 	err := mockRepo.DeleteApp(id)

// 	// Assertions
// 	assert.Nil(t, err)
// 	mockRepo.AssertExpectations(t)
// }
// func TestGetUserByEmail(t *testing.T) {
// 	mockRepo := new(MockDBRepo)
// 	email := "admin@example.com"
// 	expectedUser := &models.User{Email: email}

// 	// Set up expectations
// 	mockRepo.On("GetUserByEmail", email).Return(expectedUser, nil)

// 	// Use the mock in a function that relies on DatabaseRepo
// 	user, err := mockRepo.GetUserByEmail(email)
// 	t.Log("expuser: ", expectedUser)
// 	t.Log("user: ", user)

// 	// Assertions
// 	assert.Nil(t, err)
// 	assert.Equal(t, expectedUser, user)
// 	mockRepo.AssertExpectations(t)
// }
// func TestGetUserByID(t *testing.T) {
// 	mockRepo := new(MockDBRepo)
// 	id := 1
// 	expectedUser := &models.User{ID: id}

// 	// Set up expectations
// 	mockRepo.On("GetUserByID", id).Return(expectedUser, nil)

// 	// Use the mock in a function that relies on DatabaseRepo
// 	user, err := mockRepo.GetUserByID(id)
// 	t.Log("expuser: ", expectedUser)
// 	t.Log("user: ", user)

// 	// Assertions
// 	assert.Nil(t, err)
// 	assert.Equal(t, expectedUser, user)
// 	mockRepo.AssertExpectations(t)
// }
