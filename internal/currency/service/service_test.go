package service_test

import (
	"fmt"
	"genesis-currency-api/internal/currency/service"
	"genesis-currency-api/mocks"
	"genesis-currency-api/pkg/dto"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ServiceSuite struct {
	suite.Suite
	sut                  *service.Service
	currencyProviderMock *mocks.Provider
}

func TestServiceImplSuite(t *testing.T) {
	suite.Run(t, &ServiceSuite{})
}

func (suite *ServiceSuite) SetupTest() {
	suite.currencyProviderMock = new(mocks.Provider)

	suite.sut = service.NewService(suite.currencyProviderMock)
}

func (suite *ServiceSuite) TestGetCurrencyRate_whenNoCachedValue_checkResult() {
	// SETUP
	number := 42.2
	suite.currencyProviderMock.On("GetCurrencyRate").Return(&dto.CurrencyResponseDTO{Number: number}, nil)

	// ACT
	responseDTO, err := suite.sut.GetCurrencyRate()
	suite.Require().Nil(err)

	// VERIFY
	suite.Equal(responseDTO.Number, number)
	suite.currencyProviderMock.AssertExpectations(suite.T())
}

func (suite *ServiceSuite) TestGetCurrencyRate_whenExistsCachedValue_checkResult() {
	// SETUP
	number := 42.2
	suite.currencyProviderMock.On("GetCurrencyRate").Return(&dto.CurrencyResponseDTO{Number: number}, nil)

	// ACT
	err := suite.sut.UpdateCurrencyRates()
	suite.Require().Nil(err)

	responseDTO, err := suite.sut.GetCurrencyRate()
	suite.Require().Nil(err)

	// VERIFY
	suite.Equal(responseDTO.Number, number)
	suite.currencyProviderMock.AssertExpectations(suite.T())
}

func (suite *ServiceSuite) TestGetCurrencyRate_whenError() {
	// SETUP
	suite.currencyProviderMock.On("GetCurrencyRate").Return(nil, fmt.Errorf("test error"))

	// ACT
	responseDTO, err := suite.sut.GetCurrencyRate()

	// VERIFY
	suite.NotNil(err)
	suite.Contains(err.Error(), "test error")
	suite.Equal(float64(0), responseDTO.Number)
	suite.currencyProviderMock.AssertExpectations(suite.T())
}

func (suite *ServiceSuite) TestGetCachedCurrency_whenExistsCachedValue_checkResult() {
	// SETUP
	number := 42.2
	suite.currencyProviderMock.On("GetCurrencyRate").Return(&dto.CurrencyResponseDTO{Number: number}, nil)

	// ACT
	cachedCurrency, err := suite.sut.GetCachedCurrency()
	suite.Require().Nil(err)

	// VERIFY
	suite.Equal(cachedCurrency.Number, number)
	suite.NotNil(cachedCurrency.UpdateDate)
	suite.currencyProviderMock.AssertExpectations(suite.T())
}

func (suite *ServiceSuite) TestGetCachedCurrency_whenNoCachedValue_checkResult() {
	// SETUP
	number := 42.2
	suite.currencyProviderMock.On("GetCurrencyRate").Return(&dto.CurrencyResponseDTO{Number: number}, nil)

	// ACT
	err := suite.sut.UpdateCurrencyRates()
	suite.Require().Nil(err)

	cachedCurrency, err := suite.sut.GetCachedCurrency()
	suite.Require().Nil(err)

	// VERIFY
	suite.Equal(cachedCurrency.Number, number)
	suite.NotNil(cachedCurrency.UpdateDate)
	suite.currencyProviderMock.AssertExpectations(suite.T())
}

func (suite *ServiceSuite) TestGetCachedCurrency_whenError() {
	// SETUP
	suite.currencyProviderMock.On("GetCurrencyRate").Return(nil, fmt.Errorf("test error"))

	// ACT
	cachedCurrency, err := suite.sut.GetCachedCurrency()

	// VERIFY
	suite.NotNil(err)
	suite.Contains(err.Error(), "test error")
	suite.Equal(float64(0), cachedCurrency.Number)
	suite.NotNil(cachedCurrency.UpdateDate)
	suite.currencyProviderMock.AssertExpectations(suite.T())
}

func (suite *ServiceSuite) TestUpdateCurrencyRates_checkResult() {
	// SETUP
	number := 42.2
	suite.currencyProviderMock.On("GetCurrencyRate").Return(&dto.CurrencyResponseDTO{Number: number}, nil)

	// ACT
	err := suite.sut.UpdateCurrencyRates()

	// VERIFY
	suite.Nil(err)
	suite.currencyProviderMock.AssertExpectations(suite.T())
}

func (suite *ServiceSuite) TestUpdateCurrencyRates_whenError() {
	// SETUP
	suite.currencyProviderMock.On("GetCurrencyRate").Return(nil, fmt.Errorf("test error"))

	// ACT
	err := suite.sut.UpdateCurrencyRates()

	// VERIFY
	suite.NotNil(err)
	suite.Contains(err.Error(), "test error")
	suite.currencyProviderMock.AssertExpectations(suite.T())
}
