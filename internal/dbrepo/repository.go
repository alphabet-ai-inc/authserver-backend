package dbrepo

import (
	"database/sql"

	"backend/internal/models"

	"github.com/stretchr/testify/mock"
)

type DatabaseRepo interface {
	ConnectToDB(dsn string) (*sql.DB, error)
	Connection() (*sql.DB, error)
	AllApps() ([]*models.ThisApp, error)
	ThisApp(id int) (*models.ThisApp, error)
	ThisAppForEdit(id int) (*models.ThisApp, error)
	InsertApp(newapp models.NewApp) (int, error)
	UpdateApp(thisApp models.ThisApp) error
	DeleteApp(id int) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
}

type MockDBRepo struct {
	mock.Mock
	DatabaseRepo
	Users map[string]models.User
}

func (m *MockDBRepo) ConnectToDB(dsn string) (*sql.DB, error) {
	args := m.Called(dsn)
	return args.Get(0).(*sql.DB), args.Error(1)
}
func (m *MockDBRepo) Connection() (*sql.DB, error) {
	args := m.Called()
	return args.Get(0).(*sql.DB), args.Error(1)
}

func (m *MockDBRepo) AllApps() ([]*models.ThisApp, error) {
	args := m.Called()
	return args.Get(0).([]*models.ThisApp), args.Error(1)
}

func (m *MockDBRepo) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockDBRepo) GetUserByID(id int) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockDBRepo) ThisAppForEdit(id int) (*models.ThisApp, error) {
	args := m.Called(id)
	return args.Get(0).(*models.ThisApp), args.Error(1)
}

func (m *MockDBRepo) ThisApp(id int) (*models.ThisApp, error) {
	args := m.Called(id)
	return args.Get(0).(*models.ThisApp), args.Error(1)
}

func (m *MockDBRepo) InsertApp(newapp models.NewApp) (int, error) {
	args := m.Called(newapp)
	return args.Int(0), args.Error(1)
}

func (m *MockDBRepo) UpdateApp(thisapp models.ThisApp) error {
	args := m.Called(thisapp)
	return args.Error(0)
}

func (m *MockDBRepo) DeleteApp(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
