package intergration_test

import (
	"context"
	"errors"
	"fmt"
	"genesis-currency-api/internal/model"
	"genesis-currency-api/internal/service"
	"genesis-currency-api/pkg/config"
	"genesis-currency-api/pkg/dto"
	myerrors "genesis-currency-api/pkg/errors"
	"genesis-currency-api/pkg/util"
	"github.com/docker/go-connections/nat"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // blank import needed for migration purposes
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"path/filepath"
	"testing"
)

type UserServiceImplSuite struct {
	suite.Suite
	DB        *gorm.DB
	tx        *gorm.DB
	service   service.UserService
	container testcontainers.Container
}

func TestUserServiceImplSuite(t *testing.T) {
	suite.Run(t, new(UserServiceImplSuite))
}

func (suite *UserServiceImplSuite) SetupSuite() {
	// Start container
	ctx := context.Background()

	cnf := config.DatabaseConfig{
		DBUser:     "root",
		DBPassword: "root",
		DBName:     "currency-api",
	}

	_, host, port, err := suite.createContainer(ctx, cnf)
	suite.Require().Nil(err)

	cnf.DBHost = host
	cnf.DBPort = port.Port()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cnf.DBUser, cnf.DBPassword, cnf.DBHost, cnf.DBPort, cnf.DBName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	suite.Require().Nil(err)

	suite.DB = db

	// Run migrations
	err = suite.runMigrations(dsn)
	suite.Require().Nil(err)
}

func (suite *UserServiceImplSuite) createContainer(ctx context.Context, cnf config.DatabaseConfig) (testcontainers.Container, string, nat.Port, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15.1",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     cnf.DBUser,
			"POSTGRES_PASSWORD": cnf.DBPassword,
			"POSTGRES_DB":       cnf.DBName,
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", "", err
	}

	suite.container = container

	// Connect to database
	host, err := container.Host(ctx)
	if err != nil {
		return nil, "", "", err
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, "", "", err
	}

	return container, host, port, nil
}

func (suite *UserServiceImplSuite) runMigrations(dsn string) error {
	migrationPath, err := getMigrationsPath()
	if err != nil {
		return err
	}

	fullPath := fmt.Sprintf("file:%s", migrationPath)

	m, err := migrate.New(fullPath, dsn)
	if err != nil {
		return err
	}
	defer func(m *migrate.Migrate) {
		err, _ := m.Close()
		suite.Require().Nil(err)
	}(m)

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	log.Println("migration done")

	return nil
}

func getMigrationsPath() (string, error) {
	rootPath, err := util.GetProjectRootPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(rootPath, "pkg", "common", "db", "migrations"), nil
}

func (suite *UserServiceImplSuite) TearDownSuite() {
	err := suite.container.Terminate(context.Background())
	suite.Require().Nil(err)
}

func (suite *UserServiceImplSuite) SetupTest() {
	// Start a new transaction
	tx := suite.DB.Begin()
	suite.Require().NotNil(tx)
	suite.tx = tx

	suite.service = service.NewUserServiceImpl(tx)
}

func (suite *UserServiceImplSuite) TearDownTest() {
	// Roll back the transaction
	err := suite.tx.Rollback().Error
	suite.Require().Nil(err)
}

func (suite *UserServiceImplSuite) TestSave_checkResult() {
	// SETUP
	saveRequestDto := dto.UserSaveRequestDTO{
		Email: "test@example.com",
	}

	// ACT
	user, err := suite.service.Save(saveRequestDto)

	// VERIFY
	suite.Nil(err)
	suite.Equal(saveRequestDto.Email, user.Email)
}

func (suite *UserServiceImplSuite) TestSave_whenUserAlreadyExists() {
	// SETUP
	saveRequestDto := dto.UserSaveRequestDTO{
		Email: "exists@example.com",
	}
	var userWithEmailExistsError *myerrors.UserWithEmailExistsError

	// ACT
	_, err := suite.service.Save(saveRequestDto)
	suite.Nil(err)

	user, err := suite.service.Save(saveRequestDto)

	// VERIFY
	suite.Equal(int64(0), user.ID)
	suite.Equal("", user.Email)
	suite.NotNil(err)

	suite.True(errors.As(err, &userWithEmailExistsError))
}

func (suite *UserServiceImplSuite) TestGetAll_dbIsEmpty_checkResult() {
	// ACT
	users, err := suite.service.GetAll()

	// VERIFY
	suite.Nil(err)
	suite.NotNil(users)
	suite.Equal(0, len(users))
}

func (suite *UserServiceImplSuite) TestGetAll_dbContainsUsers_checkResult() {
	// SETUP
	usersToSave := []model.User{
		{Email: "user1@example.com"},
		{Email: "user2@example.com"},
		{Email: "user3@example.com"},
	}

	for _, u := range usersToSave {
		r := suite.tx.Create(&u)
		suite.Nil(r.Error)
	}

	// ACT
	users, err := suite.service.GetAll()

	// VERIFY
	suite.Nil(err)
	suite.NotNil(users)
	suite.Equal(3, len(users))
}
