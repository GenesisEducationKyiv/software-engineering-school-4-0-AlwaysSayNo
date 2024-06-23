package currency_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"genesis-currency-api/pkg/config"
	myerrors "genesis-currency-api/pkg/errors"
	"github.com/stretchr/testify/suite"
)

type CurrencyServiceSuite struct {
	suite.Suite
	sut CurrencyServiceInterface
}

func TestCurrencyServiceImplSuite(t *testing.T) {
	suite.Run(t, &CurrencyServiceSuite{})
}

func (suite *CurrencyServiceSuite) TestGetCurrencyInfo_checkResult() {
	// SETUP
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/external-api" {
			suite.Failf("Expected to request '/fixedvalue', got: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(`[{"ccy":"USD","base_ccy":"UAH","buy":"39.95","sale":"40.87"}]`))
		if err != nil {
			return
		}
	}))

	defer server.Close()
	suite.sut = NewCurrencyService(config.CurrencyRaterConfig{
		ThirdPartyAPIPrivateBank: server.URL + "/external-api",
	})

	// ACT
	currencyInfo, err := suite.sut.GetCurrencyInfo()
	suite.Require().Nil(err)

	// VERIFY
	suite.Equal(currencyInfo.FromCcy, "USD")
	suite.Equal(currencyInfo.ToCcy, "UAH")
	suite.Equal(currencyInfo.BuyRate, 39.95)
	suite.Equal(currencyInfo.SaleRate, 40.87)
	suite.NotNil(currencyInfo.UpdateDate)
}

func (suite *CurrencyServiceSuite) TestGetCurrencyRate_checkResult() {
	// SETUP
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/external-api" {
			suite.Failf("Expected to request '/fixedvalue', got: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(`[{"ccy":"USD","base_ccy":"UAH","buy":"39.95","sale":"40.87"}]`))
		if err != nil {
			return
		}
	}))

	defer server.Close()
	suite.sut = NewCurrencyService(config.CurrencyRaterConfig{
		ThirdPartyAPIPrivateBank: server.URL + "/external-api",
	})

	// ACT
	currencyRate, err := suite.sut.GetCurrencyRate()
	suite.Require().Nil(err)

	// VERIFY
	suite.Equal(currencyRate.Number, 40.87)
}

func (suite *CurrencyServiceSuite) TestUpdateCurrencyRates_errWhileGet() {
	// SETUP
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/external-api" {
			suite.Failf("Expected to request '/fixedvalue', got: %s", r.URL.Path)
		}
		http.Error(w, "simulated error", http.StatusInternalServerError)
	}))

	defer server.Close()
	suite.sut = NewCurrencyService(config.CurrencyRaterConfig{
		ThirdPartyAPIPrivateBank: server.URL + "/external-api",
	})

	var apiErr *myerrors.APIError

	// ACT
	err := suite.sut.UpdateCurrencyRates()

	// VERIFY
	suite.NotNil(err)
	suite.True(errors.As(err, &apiErr))
}

func (suite *CurrencyServiceSuite) TestUpdateCurrencyRates_nonOKStatusCode() {
	// SETUP
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/external-api" {
			suite.Failf("Expected to request '/fixedvalue', got: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{`))
		if err != nil {
			return
		}
	}))

	defer server.Close()
	suite.sut = NewCurrencyService(config.CurrencyRaterConfig{
		ThirdPartyAPIPrivateBank: server.URL + "/external-api",
	})

	var apiErr *myerrors.InvalidStateError

	// ACT
	err := suite.sut.GetCurrencyInfo()

	// VERIFY
	suite.NotNil(err)
	suite.True(errors.As(err, &apiErr))
}

func (suite *CurrencyServiceSuite) TestUpdateCurrencyRates_invalidURL() {
	// SETUP
	suite.sut = NewCurrencyService(config.CurrencyRaterConfig{
		ThirdPartyAPIPrivateBank: "invalid-url",
	})

	var apiErr *myerrors.APIError

	// ACT
	err := suite.sut.GetCurrencyInfo()

	// VERIFY
	suite.NotNil(err)
	suite.True(errors.As(err, &apiErr))
}

func (suite *CurrencyServiceSuite) TestGetCurrencyInfo_noCurrencyUSD() {
	// SETUP
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/external-api" {
			suite.Failf("Expected to request '/fixedvalue', got: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(`[{"ccy":"EUR","base_ccy":"UAH","buy":"42.52","sale":"43.24"}]`))
		if err != nil {
			return
		}
	}))

	defer server.Close()
	suite.sut = NewCurrencyService(config.CurrencyRaterConfig{
		ThirdPartyAPIPrivateBank: server.URL + "/external-api",
	})

	var apiErr *myerrors.APIError

	// ACT
	err := suite.sut.GetCurrencyInfo()

	// VERIFY
	suite.NotNil(err)
	suite.True(errors.As(err, &apiErr))
}
