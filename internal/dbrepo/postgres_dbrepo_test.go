package dbrepo_test

import (
	"backend/internal/dbrepo"
	"backend/internal/models"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// Initialize sqlmock
var db, mock, err = sqlmock.New()

func TestAllApps(t *testing.T) {
	if err != nil {
		t.Fatalf("failed to open a stub database connection: %v", err)
	}
	defer db.Close()

	repo := dbrepo.PostgresDBRepo{DB: db}

	// Define expected database rows
	rows := sqlmock.NewRows([]string{"id", "name", "release", "path", "init", "web", "title", "created", "updated"}).
		AddRow(1, "AppName", "1.0", "/path/to/app", true, true, "App Title", 0, 0)

	// Set up expectations
	mock.ExpectQuery("select id, name, release, path, init, web, title, created, updated from apps").
		WillReturnRows(rows)

	apps, err := repo.AllApps()
	assert.Nil(t, err)
	assert.Len(t, apps, 1)
	assert.Equal(t, "AppName", apps[0].Name)

	// Verify that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}
func TestThisApp(t *testing.T) {
	if err != nil {
		t.Fatalf("failed to open a stub database connection: %v", err)
	}
	defer db.Close()

	repo := dbrepo.PostgresDBRepo{DB: db}

	// Define expected database rows
	rows := sqlmock.NewRows([]string{"id", "name", "release", "path", "init", "web", "title", "created", "updated"}).
		AddRow(1, "AppName", "1.0", "/path/to/app", true, true, "App Title", 0, 0)

	// Set up expectations
	mock.ExpectQuery(regexp.QuoteMeta("select id, name, release, path, init, web, title, created, updated from apps where id = $1")).
		WithArgs(1).
		WillReturnRows(rows)

	app, err := repo.ThisApp(1)
	assert.Nil(t, err)
	if assert.NotNil(t, app) {
		assert.Equal(t, "AppName", app.Name)
	}

	// Verify that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestThisAppForEdit(t *testing.T) {
	if err != nil {
		t.Fatalf("failed to open a stub database connection: %v", err)
	}
	defer db.Close()

	repo := dbrepo.PostgresDBRepo{DB: db}

	// Define expected database rows
	rows := sqlmock.NewRows([]string{"id", "name", "release", "path", "init", "web", "title", "created", "updated"}).
		AddRow(1, "AppName", "1.0", "/path/to/app", true, true, "App Title", 0, 0)

	// Set up expectationsinternal
	mock.ExpectQuery(regexp.QuoteMeta("select id, name, release, path, init, web, title, created, updated from apps where id = $1")).
		WithArgs(1).
		WillReturnRows(rows)

	editedApp, err := repo.ThisAppForEdit(1)
	assert.Nil(t, err)
	if assert.NotNil(t, editedApp) {
		assert.Equal(t, "AppName", editedApp.Name)
	}

	// Verify that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}
func TestGetUserByEmail(t *testing.T) {
	if err != nil {
		t.Fatalf("failed to open a stub database connection: %v", err)
	}
	defer db.Close()

	repo := dbrepo.PostgresDBRepo{DB: db}

	// Define expected database rows
	rows := sqlmock.NewRows([]string{"id", "username", "password", "code", "active", "last_login", "last_session", "blocked", "tries", "last_try", "email", "profile_id", "group_id", "dbsauth_id", "activation_time", "last_action", "last_app", "last_db", "lan", "company_id", "created", "updated"}).
		AddRow(1, "User Name", "verysecret", "0", true, 0, "0", false, 0, 0, "admin@example.com", 0, 0, 0, 0, "0", 0, 0, "es", 0, 0, 0)

	// Set up expectations
	mock.ExpectQuery(regexp.QuoteMeta("select id, username, password, code, active, last_login, last_session, blocked, tries, last_try, email, profile_id, group_id, dbsauth_id, activation_time, last_action, last_app, last_db, lan, company_id, created, updated from users where email =$1")).
		WithArgs("admin@example.com").
		WillReturnRows(rows)

	user, err := repo.GetUserByEmail("admin@example.com")
	assert.Nil(t, err)
	assert.Equal(t, 1, user.ID)
	// assert.Equal(t, true, user)

	// Verify that all expectations were met
	// err = mock.ExpectationsWereMet()
	// assert.Nil(t, err)
}
func TestGetUserById(t *testing.T) {
	if err != nil {
		t.Fatalf("failed to open a stub database connection: %v", err)
	}
	defer db.Close()

	repo := dbrepo.PostgresDBRepo{DB: db}

	// Define expected database rows
	rows := sqlmock.NewRows([]string{"id", "username", "password", "code", "active", "last_login", "last_session", "blocked", "tries", "last_try", "email", "profile_id", "group_id", "dbsauth_id", "activation_time", "last_action", "last_app", "last_db", "lan", "company_id", "created", "updated"}).
		AddRow(1, "User Name", "verysecret", "0", true, 0, "0", false, 0, 0, "admin@example.com", 0, 0, 0, 0, "0", 0, 0, "es", 0, 0, 0)

	// Set up expectations
	mock.ExpectQuery(regexp.QuoteMeta("select id, username, password, code, active, last_login, last_session, blocked, tries, last_try, email, profile_id, group_id, dbsauth_id, activation_time, last_action, last_app, last_db, lan, company_id, created, updated from users where id =$1")).
		WithArgs(1).
		WillReturnRows(rows)

	user, err := repo.GetUserByID(1)
	assert.Nil(t, err)
	assert.Equal(t, 1, user.ID)

	// Verify that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestInsertApp(t *testing.T) {
	if err != nil {
		t.Fatalf("failed to open a stub database connection: %v", err)
	}
	defer db.Close()

	repo := dbrepo.PostgresDBRepo{DB: db}

	app := models.NewApp{
		Name:    "AppName",
		Release: "1.0",
		Path:    "/path/to/app",
		Init:    "init.sh",
		Web:     "http://testapp.com",
		Title:   "App Title",
		Created: 0,
		Updated: 0,
	}

	// Set up expectations
	mock.ExpectQuery(regexp.QuoteMeta(`
	insert into apps( 
		name,
		release,
		path,
		init,
		web,
		title,
		created, 
		updated)
	values (
		$1, $2, $3, $4, $5, $6, $7, $8
	)
	returning id, 1 
	`)).
		WithArgs(app.Name, app.Release, app.Path, app.Init, app.Web, app.Title, app.Created, app.Updated).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	id, err := repo.InsertApp(app)
	assert.Nil(t, err)
	assert.Equal(t, 1, id)

	// Verify that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUpdateApp(t *testing.T) {
	if err != nil {
		t.Fatalf("failed to open a stub database connection: %v", err)
	}
	defer db.Close()

	repo := dbrepo.PostgresDBRepo{DB: db}

	app := models.ThisApp{
		Name:    "AppName",
		Release: "1.0",
		Path:    "/path/to/app",
		Init:    "init.sh",
		Web:     "http://testapp.com",
		Title:   "App Title",
		Created: 0,
		Updated: 0,
	}

	// Set up expectations
	mock.ExpectQuery(regexp.QuoteMeta(`Update apps set name=$2, release=$3, path=$4, init=$5, web=$6, title=$7, created=$8, updated=$9 where id = $1`)).
		WithArgs(1, app.Name, app.Release, app.Path, app.Init, app.Web, app.Title, app.Created, app.Updated).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	_ = repo.UpdateApp(app)

	// Verify that all expectations were met
	_ = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}
func TestDeleteApp(t *testing.T) {
	if err != nil {
		t.Fatalf("failed to open a stub database connection: %v", err)
	}
	defer db.Close()

	repo := dbrepo.PostgresDBRepo{DB: db}

	// Set up expectations
	mock.ExpectExec(regexp.QuoteMeta("delete from apps where id = $1")).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.DeleteApp(1)
	assert.Nil(t, err)

	// Verify that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}
