package dbrepo_test

import (
	"authserver-backend/internal/dbrepo"
	"authserver-backend/internal/models"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

const mockedfields = `
		name,
		release,
		path,
		init,
		web,
		title,
		created,
		updated
	`

// Helper function to set up a mock database and return a PostgresDBRepo instance
func setupMockDB(t *testing.T) (*dbrepo.PostgresDBRepo, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	repo := &dbrepo.PostgresDBRepo{DB: db}
	return repo, mock, func() { db.Close() }
}

// Test cases for PostgresDBRepo methods. These tests use sqlmock to simulate database interactions.
// The tests cover various scenarios including successful queries, errors, and edge cases.
// Same tests can be written for other methods in the PostgresDBRepo struct and for other database implementations.
func TestAllApps(t *testing.T) {
	repo, mock, closeFn := setupMockDB(t)
	defer closeFn()

	rows := sqlmock.NewRows([]string{
		"id", "name", "release", "path", "init", "web", "title", "created", "updated",
	}).AddRow(
		1, "App1", "v1.0", "/app1", "./app1.sh", "http://app1.example.com", "Title1", 16000000, 16000000,
	).AddRow(
		2, "App2", "v2.0", "/app2", "./app2.sh", "http://app2.example.com", "Title2", 16000000, 16000000,
	)
	mock.ExpectQuery(regexp.QuoteMeta(`
		select
			id,
			name,
			release,
			path,
			init,
			web,
			title,
			created,
			updated
		from
			apps
		order by
			name
	`)).WillReturnRows(rows)

	apps, err := repo.AllApps(mockedfields)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(apps) != 2 {
		t.Errorf("expected 2 apps, got %d", len(apps))
	}
}

func TestThisApp(t *testing.T) {
	repo, mock, closeFn := setupMockDB(t)
	defer closeFn()

	now := 160000000 // Simulated timestamp for testing
	row := sqlmock.NewRows([]string{
		"id", "name", "release", "path", "init", "web", "title", "created", "updated",
	}).AddRow(
		1, "App1", "v1.0", "/app1", "./app1.sh", "http://app1.example.com", "Title1", now, now,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`
	select
		id,
		name,
		release,
		path,
		init,
		web,
		title,
		created,
		updated
	from
		apps
	where id = $1
`)).WithArgs(1).WillReturnRows(row)

	app, err := repo.ThisApp(1, mockedfields)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if app.ID != 1 || app.Name != "App1" {
		t.Errorf("unexpected app data: %+v", app)
	}
}

func TestThisAppForEdit(t *testing.T) {
	repo, mock, closeFn := setupMockDB(t)
	defer closeFn()

	now := 160000000 // Simulated timestamp for testing
	row := sqlmock.NewRows([]string{
		"id", "name", "release", "path", "init", "web", "title", "created", "updated",
	}).AddRow(
		2, "App2", "v1.0", "/app2", "./app2.sh", "http://app2.example.com", "Title1", now, now,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`
	select
		id,
		name,
		release,
		path,
		init,
		web,
		title,
		created,
		updated
	from
		apps
	where id = $1
`)).WithArgs(2).WillReturnRows(row)

	app, err := repo.ThisAppForEdit(2, mockedfields)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if app.ID != 2 || app.Name != "App2" {
		t.Errorf("unexpected app data: %+v", app)
	}
}

func TestGetUserByEmail(t *testing.T) {
	repo, mock, closeFn := setupMockDB(t)
	defer closeFn()

	now := 160000000 // Simulated timestamp for testing
	row := sqlmock.NewRows([]string{
		"id", "username", "password", "code", "active", "last_login", "last_session", "blocked",
		"tries", "last_try", "email", "profile_id", "group_id", "dbsauth_id", "activation_time",
		"last_action", "last_app", "last_db", "lan", "company_id", "created", "updated",
	}).AddRow(
		1, "user1", "pass", "code", true, now, now, false,
		0, now, "user1@example.com", 1, 1, 1, now,
		now, 1, 1, "en", 1, now, now,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`select id, username, password, code, active, last_login, last_session, blocked,
	tries, last_try, email, profile_id, group_id, dbsauth_id, activation_time, last_action,
	last_app, last_db, lan, company_id, created, updated from users where email =$1`)).
		WithArgs("user1@example.com").WillReturnRows(row)

	user, err := repo.GetUserByEmail("user1@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Email != "user1@example.com" {
		t.Errorf("unexpected user data: %+v", user)
	}
}

func TestGetUserByEmail_EmptyEmail(t *testing.T) {
	repo, _, closeFn := setupMockDB(t)
	defer closeFn()

	_, err := repo.GetUserByEmail("")
	if err == nil {
		t.Error("expected error for empty email, got nil")
	}
}

func TestGetUserByID(t *testing.T) {
	repo, mock, closeFn := setupMockDB(t)
	defer closeFn()

	now := 160000000 // Simulated timestamp for testing
	row := sqlmock.NewRows([]string{
		"id", "username", "password", "code", "active", "last_login", "last_session", "blocked",
		"tries", "last_try", "email", "profile_id", "group_id", "dbsauth_id", "activation_time",
		"last_action", "last_app", "last_db", "lan", "company_id", "created", "updated",
	}).AddRow(
		2, "user2", "pass2", "code2", false, now, now, true,
		1, now, "user2@example.com", 2, 2, 2, now,
		now, 2, 2, "es", 2, now, now,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`select id, username, password, code, active, last_login, last_session, blocked,
	tries, last_try, email, profile_id, group_id, dbsauth_id, activation_time, last_action,
	last_app, last_db, lan, company_id, created, updated from users where id =$1`)).
		WithArgs(2).WillReturnRows(row)

	user, err := repo.GetUserByID(2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != 2 || user.UserName != "user2" {
		t.Errorf("unexpected user data: %+v", user)
	}
}

func TestInsertApp(t *testing.T) {
	repo, mock, closeFn := setupMockDB(t)
	defer closeFn()

	newApp := models.NewApp{
		Name:    "App3",
		Release: "v3.0",
		Path:    "/app3",
		Init:    "./init.sh",
		Web:     "http://app3.example.com",
		Title:   "Title3",
		Created: 160000000,
		Updated: 160000000,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`insert into apps(`+mockedfields+`) values ($2,$3,$4,$5,$6,$7,$8,$9) returning id`)).
		WithArgs(
			&newApp.Name,
			&newApp.Release,
			&newApp.Path,
			&newApp.Init,
			&newApp.Web,
			&newApp.Title,
			&newApp.Created,
			&newApp.Updated,
		).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(3))

	id, err := repo.InsertApp(newApp, mockedfields)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 3 {
		t.Errorf("expected id 3, got %d", id)
	}
}

func TestUpdateApp(t *testing.T) {
	repo, mock, closeFn := setupMockDB(t)
	defer closeFn()

	// Simulated app data for update
	app := models.ThisApp{
		ID: 1,
		NewApp: models.NewApp{
			Name:    "App1",
			Release: "v1.1",
			Path:    "/app1",
			Init:    "./updated_init.sh",
			Web:     "http://updated-app1.example.com",
			Title:   "Title1",
			Created: 160000000,
			Updated: 160000000,
		},
	}

	mock.ExpectExec(regexp.QuoteMeta(`update apps set name = $2, release = $3, path = $4, init = $5, web = $6, title = $7, created = $8, updated = $9 where id = $1`)).
		WithArgs(
			app.ID,
			app.Name,
			app.Release,
			app.Path,
			app.Init,
			app.Web,
			app.Title,
			app.Created,
			app.Updated,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

		// Call the UpdateApp method
	err := repo.UpdateApp(app, mockedfields)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteApp(t *testing.T) {
	repo, mock, closeFn := setupMockDB(t)
	defer closeFn()

	mock.ExpectExec(regexp.QuoteMeta(`delete from apps where id = $1`)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteApp(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAllApps_QueryError(t *testing.T) {
	repo, mock, closeFn := setupMockDB(t)
	defer closeFn()

	mock.ExpectQuery("select(.|\\s)*from\\s+apps").WillReturnError(errors.New("query error"))

	_, err := repo.AllApps(mockedfields)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestThisApp_ScanError(t *testing.T) {
	repo, mock, closeFn := setupMockDB(t)
	defer closeFn()

	row := sqlmock.NewRows([]string{
		"id", "name", "release", "path", "init", "web", "title", "created", "updated",
	}).AddRow(
		"bad_id", "App1", "v1.0", "/app1", "./scan.sh", "http://error.example.com", "Title1", 160000000, 160000000,
	)

	mock.ExpectQuery("select(.|\\s)*from\\s+apps(.|\\s)*where id = \\$1").
		WithArgs(1).WillReturnRows(row)

	_, err := repo.ThisApp(1, mockedfields)
	if err == nil {
		t.Error("expected scan error, got nil")
	}
}

func TestGetUserByID_NotFound(t *testing.T) {
	repo, mock, closeFn := setupMockDB(t)
	defer closeFn()

	mock.ExpectQuery("select(.|\\s)*from users where id = \\$1").
		WithArgs(99).WillReturnRows(sqlmock.NewRows([]string{
		"id", "username", "password", "code", "active", "last_login", "last_session", "blocked",
		"tries", "last_try", "email", "profile_id", "group_id", "dbsauth_id", "activation_time",
		"last_action", "last_app", "last_db", "lan", "company_id", "created", "updated",
	}))

	_, err := repo.GetUserByID(99)
	if err == nil {
		t.Error("expected error for not found, got nil")
	}
}
