package dbrepo

// repository.go is the file that contains the repository struct and methods for database operations.
// It implements the DatabaseRepo interface. Any database motore (Postgres, MySQL, SQLite, etc.)
// can be used by creating a new struct that implements the DatabaseRepo interface.
import (
	"database/sql"

	"authserver-backend/internal/models"

	"github.com/stretchr/testify/mock"
)

// DatabaseRepo is the interface that wraps the basic methods for database operations. When testing,
// we use the mock implementation of this interface. When running the application, we use the
// actual implementation in postgresdb.go.
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

// MockDBRepo is a mock implementation of the DatabaseRepo interface for testing purposes.
// In the future, we can use a mocking library like testify to create a more robust mock.
// For now, we will implement the methods manually.
type MockDBRepo struct {
	mock.Mock
	DatabaseRepo
	Users models.User
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
