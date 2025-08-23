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

type PostgresDBRepo struct {
	DB *sql.DB
}

const dbTimeout = time.Second * 3

// NewDatabase initializes a new database connection
func (m *PostgresDBRepo) ConnectToDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logerror.LogError(err)
	}

	return db, err
}

func (m *PostgresDBRepo) Connection() (*sql.DB, error) {
	if m.DB == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}
	return m.DB, nil

}

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
	returning id, 1 
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
