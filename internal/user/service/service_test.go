package service_test

import (
	"errors"
	"fmt"
	"genesis-currency-api/internal/user/model"
	"genesis-currency-api/internal/user/service"
	"genesis-currency-api/mocks"
	"genesis-currency-api/pkg/dto"
	myerrors "genesis-currency-api/pkg/errors"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ServiceSuite struct {
	suite.Suite
	sut                *service.Service
	userRepositoryMock *mocks.Repository
}

func TestServiceImplSuite(t *testing.T) {
	suite.Run(t, &ServiceSuite{})
}

func (suite *ServiceSuite) SetupTest() {
	suite.userRepositoryMock = new(mocks.Repository)

	suite.sut = service.NewService(suite.userRepositoryMock)
}

func (suite *ServiceSuite) TestSave_checkResult() {
	// SETUP
	saveRequestDto := dto.UserSaveRequestDTO{
		Email: "test@example.com",
	}
	userModel := model.User{
		Email: saveRequestDto.Email,
	}
	savedUser := model.User{
		Email: saveRequestDto.Email,
		ID:    1,
	}
	suite.userRepositoryMock.On("ExistsByEmail", userModel.Email).Return(false)
	suite.userRepositoryMock.On("Create", userModel).Return(&savedUser, nil)

	// ACT
	returnedUser, err := suite.sut.Save(saveRequestDto)
	suite.Require().Nil(err)

	// VERIFY
	suite.Equal(returnedUser.Email, saveRequestDto.Email)
	suite.Equal(int64(1), returnedUser.ID)
	suite.userRepositoryMock.AssertExpectations(suite.T())
}

func (suite *ServiceSuite) TestSave_whenUserAlreadyExists() {
	// SETUP
	saveRequestDto := dto.UserSaveRequestDTO{
		Email: "test@example.com",
	}
	userModel := model.User{
		Email: saveRequestDto.Email,
	}
	var userWithEmailExistsError *myerrors.UserWithEmailExistsError

	suite.userRepositoryMock.On("ExistsByEmail", userModel.Email).Return(true)

	// ACT
	returnedUser, err := suite.sut.Save(saveRequestDto)
	suite.Require().Nil(returnedUser)

	// VERIFY
	suite.NotNil(err)
	suite.True(errors.As(err, &userWithEmailExistsError))
	suite.userRepositoryMock.AssertExpectations(suite.T())
}

func (suite *ServiceSuite) TestSave_whenErrorWhileCreating() {
	// SETUP
	saveRequestDto := dto.UserSaveRequestDTO{
		Email: "test@example.com",
	}
	userModel := model.User{
		Email: saveRequestDto.Email,
	}

	suite.userRepositoryMock.On("ExistsByEmail", userModel.Email).Return(false)
	suite.userRepositoryMock.On("Create", userModel).
		Return(nil, fmt.Errorf("test error"))

	// ACT
	returnedUser, err := suite.sut.Save(saveRequestDto)
	suite.Require().Nil(returnedUser)

	// VERIFY
	suite.NotNil(err)
	suite.Contains(err.Error(), "test error")
	suite.userRepositoryMock.AssertExpectations(suite.T())
}

func (suite *ServiceSuite) TestGetAll_dbContainsUsers_checkResult() {
	// SETUP
	usersFromDB := []model.User{
		{Email: "user1@example.com"},
		{Email: "user2@example.com"},
		{Email: "user3@example.com"},
	}

	suite.userRepositoryMock.On("GetAll").Return(&usersFromDB, nil)

	// ACT
	usersDTO, err := suite.sut.GetAll()
	suite.Require().Nil(err)

	// VERIFY
	suite.NotNil(usersDTO)
	suite.Equal(3, len(usersDTO))
	suite.userRepositoryMock.AssertExpectations(suite.T())
}

func (suite *ServiceSuite) TestGetAll_dbIsEmpty_checkResult() {
	// SETUP
	usersFromDB := make([]model.User, 0)

	suite.userRepositoryMock.On("GetAll").Return(&usersFromDB, nil)

	// ACT
	usersDTO, err := suite.sut.GetAll()
	suite.Require().Nil(err)

	// VERIFY
	suite.NotNil(usersDTO)
	suite.Equal(0, len(usersDTO))
	suite.userRepositoryMock.AssertExpectations(suite.T())
}

func (suite *ServiceSuite) TestGetAll_whenError() {
	// SETUP
	suite.userRepositoryMock.On("GetAll").Return(nil, fmt.Errorf("test error"))

	// ACT
	usersDTO, err := suite.sut.GetAll()
	suite.Require().Nil(usersDTO)

	// VERIFY
	suite.NotNil(err)
	suite.Contains(err.Error(), "test error")
	suite.userRepositoryMock.AssertExpectations(suite.T())
}
