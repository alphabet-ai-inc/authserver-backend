package dbrepo

import (
	"authserver-backend/internal/models"
	"authserver-backend/logerror"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// PostgresDBRepo is a struct that holds the database connection pool.
// It implements the DatabaseRepo interface.
// This struct is used to interact with the PostgreSQL database.
// Changing the database type only requires changes in this file and conditionally use one or another
// in the main application code where the database connection is established.
// Or you can create a new file for the new database type and implement the same interface,
// then import the new file in the main application code and use it instead of this one.
type PostgresDBRepo struct {
	DB *sql.DB
}

const dbTimeout = time.Second * 3

// NewDatabase initializes a new database connection using the provided DSN (Data Source Name) and the PostgresDBRepo struct.
// It returns a pointer to the PostgresDBRepo struct and an error if any occurs during the connection process.

func (m *PostgresDBRepo) ConnectToDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logerror.LogError(err)
	}

	return db, err
}

// Connection returns the sql.DB connection pool.
// It returns an error if the database connection is not initialized.

func (m *PostgresDBRepo) Connection() (*sql.DB, error) {
	if m.DB == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}
	return m.DB, nil

}

// AllApps returns a slice of all applications from the database.
// This is for now the same for users and admins, but in the future it might be different.
// The usesr will see only the apps they have access to, while the admins will see all apps.
// The users will see only the fields relevant to linking to the app.
// Perhaps it is possible to use the same model for both, but selecting the fields to return in the query.
func (m *PostgresDBRepo) AllApps() ([]*models.ThisApp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	query := `
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
	`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []*models.ThisApp

	for rows.Next() {
		var thisapp models.ThisApp
		err := rows.Scan(
			&thisapp.ID,
			&thisapp.Name,
			&thisapp.Release,
			&thisapp.Path,
			&thisapp.Init,
			&thisapp.Web,
			&thisapp.Title,
			&thisapp.Created,
			&thisapp.Updated,
		)

		if err != nil {
			return nil, err
		}
		apps = append(apps, &thisapp)
	}

	return apps, nil
}

// ThisApp returns a single application by its ID from the database.
// It returns a pointer to the ThisApp struct and an error if any occurs during the query process.
func (m *PostgresDBRepo) ThisApp(id int) (*models.ThisApp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	query := `
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
`
	row := m.DB.QueryRowContext(ctx, query, id)

	var thisapp models.ThisApp

	err := row.Scan(
		&thisapp.ID,
		&thisapp.Name,
		&thisapp.Release,
		&thisapp.Path,
		&thisapp.Init,
		&thisapp.Web,
		&thisapp.Title,
		&thisapp.Created,
		&thisapp.Updated,
	)

	if err != nil {
		return nil, err
	}
	return &thisapp, err
}

// ThisAppForEdit returns a single application by its ID from the database for editing purposes.
func (m *PostgresDBRepo) ThisAppForEdit(id int) (*models.ThisApp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	query := `
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
`
	row := m.DB.QueryRowContext(ctx, query, id)

	var thisapp models.ThisApp

	err := row.Scan(
		&thisapp.ID,
		&thisapp.Name,
		&thisapp.Release,
		&thisapp.Path,
		&thisapp.Init,
		&thisapp.Web,
		&thisapp.Title,
		&thisapp.Created,
		&thisapp.Updated,
	)

	if err != nil {
		return nil, err
	}
	return &thisapp, err
}

// InsertApp inserts a new application into the database.
func (m *PostgresDBRepo) InsertApp(newapp models.NewApp) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	stmt := `
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
	returning id
	`
	var newID int

	err := m.DB.QueryRowContext(ctx, stmt,
		&newapp.Name,
		&newapp.Release,
		&newapp.Path,
		&newapp.Init,
		&newapp.Web,
		&newapp.Title,
		&newapp.Created,
		&newapp.Updated,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil

}

// UpdateApp updates an existing application in the database.
func (m *PostgresDBRepo) UpdateApp(thisapp models.ThisApp) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	stmt := `update apps set name=$2, release=$3, path=$4, init=$5, web=$6, title=$7, 
		created=$8, updated=$9
		where id = $1`

	_, err := m.DB.ExecContext(ctx, stmt,
		thisapp.ID,
		thisapp.Name,
		thisapp.Release,
		thisapp.Path,
		thisapp.Init,
		thisapp.Web,
		thisapp.Title,
		thisapp.Created,
		thisapp.Updated,
	)

	if err != nil {
		return err
	}

	return nil
}

// DeleteApp deletes an application from the database by its ID.
func (m *PostgresDBRepo) DeleteApp(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	stmt := `delete from apps where id = $1`

	_, err := m.DB.ExecContext(ctx, stmt, id)

	if err != nil {
		return err
	}

	return nil

}

// GetUserByEmail retrieves a user from the database by their email address.
func (m *PostgresDBRepo) GetUserByEmail(email string) (*models.User, error) {
	// Get user by email
	if _, err := m.Connection(); err != nil {
		log.Fatalf("Database not connected: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, username, password, code, active, last_login, last_session, blocked,
	tries, last_try, email, profile_id, group_id, dbsauth_id, activation_time, last_action,
	last_app, last_db, lan, company_id, created, updated from users where email =$1`

	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.UserName,
		&user.Password,
		&user.Code,
		&user.Active,
		&user.LastLogin,
		&user.LastSession,
		&user.Blocked,
		&user.Tries,
		&user.LastTry,
		&user.Email,
		&user.ProfileId,
		&user.GroupId,
		&user.DbsAuth,
		&user.ActivationTime,
		&user.LastAction,
		&user.LastApp,
		&user.LastDb,
		&user.Lan,
		&user.CompanyId,
		&user.Created,
		&user.Updated,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no user found with email: %s", email)
	}
	if err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}
	// return &user, err
	return &user, nil
}

// GetUserByID retrieves a user from the database by their ID.
func (m *PostgresDBRepo) GetUserByID(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	query := `select id, username, password, code, active, last_login, last_session, blocked, 
	tries, last_try, email, profile_id, group_id, dbsauth_id, activation_time, last_action, 
	last_app, last_db, lan, company_id, created, updated from users where id =$1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&user.ID,
		&user.UserName,
		&user.Password,
		&user.Code,
		&user.Active,
		&user.LastLogin,
		&user.LastSession,
		&user.Blocked,
		&user.Tries,
		&user.LastTry,
		&user.Email,
		&user.ProfileId,
		&user.GroupId,
		&user.DbsAuth,
		&user.ActivationTime,
		&user.LastAction,
		&user.LastApp,
		&user.LastDb,
		&user.Lan,
		&user.CompanyId,
		&user.Created,
		&user.Updated,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
