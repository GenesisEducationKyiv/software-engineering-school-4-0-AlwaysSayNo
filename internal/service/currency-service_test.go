package service_test

import (
	"errors"
	"genesis-currency-api/internal/service"
	"genesis-currency-api/pkg/config"
	myerrors "genesis-currency-api/pkg/errors"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type CurrencyServiceImplSuite struct {
	suite.Suite
	sut service.CurrencyService
}

func TestCurrencyServiceImplSuite(t *testing.T) {
	suite.Run(t, &CurrencyServiceImplSuite{})
}

func (csis *CurrencyServiceImplSuite) SetupTest() {

}

func (csis *CurrencyServiceImplSuite) TestGetCurrencyInfo_checkResult() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/external-api" {
			csis.Failf("Expected to request '/fixedvalue', got: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(`[{"ccy":"USD","base_ccy":"UAH","buy":"39.95","sale":"40.87"}]`))
		if err != nil {
			return
		}
	}))

	defer server.Close()
	csis.sut = service.NewCurrencyServiceImpl(config.CurrencyServiceConfig{
		ThirdPartyAPI: server.URL + "/external-api",
	})

	currencyInfo := csis.sut.GetCurrencyInfo()

	csis.Equal(currencyInfo.FromCcy, "USD")
	csis.Equal(currencyInfo.ToCcy, "UAH")
	csis.Equal(currencyInfo.BuyRate, 39.95)
	csis.Equal(currencyInfo.SaleRate, 40.87)
	csis.NotNil(currencyInfo.UpdateDate)
}

func (csis *CurrencyServiceImplSuite) TestGetCurrencyRate_checkResult() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/external-api" {
			csis.Failf("Expected to request '/fixedvalue', got: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(`[{"ccy":"USD","base_ccy":"UAH","buy":"39.95","sale":"40.87"}]`))
		if err != nil {
			return
		}
	}))

	defer server.Close()
	csis.sut = service.NewCurrencyServiceImpl(config.CurrencyServiceConfig{
		ThirdPartyAPI: server.URL + "/external-api",
	})

	currencyRate := csis.sut.GetCurrencyRate()

	csis.Equal(currencyRate.Number, 40.87)
}

func (csis *CurrencyServiceImplSuite) TestUpdateCurrencyRates_errWhileGet() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/external-api" {
			csis.Failf("Expected to request '/fixedvalue', got: %s", r.URL.Path)
		}
		http.Error(w, "simulated error", http.StatusInternalServerError)
	}))

	defer server.Close()
	csis.sut = service.NewCurrencyServiceImpl(config.CurrencyServiceConfig{
		ThirdPartyAPI: server.URL + "/external-api",
	})

	var apiErr *myerrors.APIError

	err := csis.sut.UpdateCurrencyRates()

	csis.NotNil(err)
	csis.True(errors.As(err, &apiErr))
}

func (csis *CurrencyServiceImplSuite) TestUpdateCurrencyRates_nonOKStatusCode() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/external-api" {
			csis.Failf("Expected to request '/fixedvalue', got: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{`))
		if err != nil {
			return
		}
	}))

	defer server.Close()
	csis.sut = service.NewCurrencyServiceImpl(config.CurrencyServiceConfig{
		ThirdPartyAPI: server.URL + "/external-api",
	})

	var apiErr *myerrors.InvalidStateError

	err := csis.sut.UpdateCurrencyRates()

	csis.NotNil(err)
	csis.True(errors.As(err, &apiErr))
}

func (csis *CurrencyServiceImplSuite) TestUpdateCurrencyRates_invalidURL() {
	csis.sut = service.NewCurrencyServiceImpl(config.CurrencyServiceConfig{
		ThirdPartyAPI: "invalid-url",
	})

	var apiErr *myerrors.APIError

	err := csis.sut.UpdateCurrencyRates()

	csis.NotNil(err)
	csis.True(errors.As(err, &apiErr))
}

func (csis *CurrencyServiceImplSuite) TestGetCurrencyInfo_noCurrencyUSD() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/external-api" {
			csis.Failf("Expected to request '/fixedvalue', got: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(`[{"ccy":"EUR","base_ccy":"UAH","buy":"42.52","sale":"43.24"}]`))
		if err != nil {
			return
		}
	}))

	defer server.Close()
	csis.sut = service.NewCurrencyServiceImpl(config.CurrencyServiceConfig{
		ThirdPartyAPI: server.URL + "/external-api",
	})

	var apiErr *myerrors.APIError

	err := csis.sut.UpdateCurrencyRates()

	csis.NotNil(err)
	csis.True(errors.As(err, &apiErr))
}
