package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/alfcope/checkouttest/api/requests"
	"github.com/alfcope/checkouttest/api/responses"
	"github.com/alfcope/checkouttest/datasource"
	"github.com/alfcope/checkouttest/errors"
	"github.com/alfcope/checkouttest/internal/tests/mocks"
	"github.com/alfcope/checkouttest/model"
	"github.com/alfcope/checkouttest/pkg/logging"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type CheckoutControllerTestSuite struct {
	suite.Suite

	checkoutController CheckoutController
	checkoutService    CheckoutService
	datasourceMock     datasource.Datasource
}

func TestControllerSuite(t *testing.T) {
	suite.Run(t, new(CheckoutControllerTestSuite))
}

func (suite *CheckoutControllerTestSuite) SetupSuite() {
	apiRoute := mux.NewRouter().PathPrefix("/api/v1").Subrouter().StrictSlash(true)

	suite.datasourceMock = datasource.Datasource(mocks.NewDatasourceMock())
	suite.checkoutService = NewCheckoutService(suite.datasourceMock)
	suite.checkoutController = *NewCheckoutController(apiRoute, suite.checkoutService)
}

func (suite *CheckoutControllerTestSuite) TearDownTest() {
	suite.datasourceMock.(*mocks.DatasourceMock).ExpectedCalls = nil
	suite.datasourceMock.(*mocks.DatasourceMock).Calls = nil
}

func (suite *CheckoutControllerTestSuite) TestCreateBasketDuplicatedId() {
	// Given
	basketId := uuid.New().String()
	suite.datasourceMock.(*mocks.DatasourceMock).On("AddBasket", mock.AnythingOfType("*model.Basket")).Return(errors.NewPrimaryKeyError(basketId))

	// When
	req, err := http.NewRequest("POST", "/baskets/", nil)
	if err != nil {
		suite.T().Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := logging.AccessLoggingMiddleware(suite.checkoutController.CreateBasket())

	handler.ServeHTTP(rr, req)

	// Then
	suite.Equal(http.StatusInternalServerError, rr.Code)
}

func (suite *CheckoutControllerTestSuite) TestCreateBasket() {
	// Given
	suite.datasourceMock.(*mocks.DatasourceMock).On("AddBasket", mock.AnythingOfType("*model.Basket")).Return(nil)

	// When
	req, err := http.NewRequest("POST", "/baskets/", nil)
	if err != nil {
		suite.T().Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := logging.AccessLoggingMiddleware(suite.checkoutController.CreateBasket())

	handler.ServeHTTP(rr, req)

	// Then
	suite.Equal(http.StatusCreated, rr.Code)

	var nbr = new(responses.NewBasketResponse)
	err = json.Unmarshal(rr.Body.Bytes(), &nbr)

	if err != nil {
		suite.T().Errorf("Error unmarshalling new basket response: %v", err)
	}

	suite.NotNil(nbr)
	suite.NotEqual("", nbr.Id)
}

func (suite *CheckoutControllerTestSuite) TestAddNonExistingProduct() {
	// Given
	productCode := "FAKE"
	suite.datasourceMock.(*mocks.DatasourceMock).On("GetProduct",
		mock.AnythingOfType("model.ProductCode")).Return(*new(model.Product), errors.NewProductNotFound(productCode))

	// When
	reqBodyBytes := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBytes).Encode(requests.AddItemRequest{Code: model.ProductCode(productCode)})
	if err != nil {
		suite.T().Errorf("Error encoding request: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("/baskets/%v", uuid.New().String()), bytes.NewBuffer(reqBodyBytes.Bytes()))
	if err != nil {
		suite.T().Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := logging.AccessLoggingMiddleware(suite.checkoutController.AddItem())

	handler.ServeHTTP(rr, req)

	// Then
	suite.Equal(http.StatusNotFound, rr.Code)

	//Checking there has not been any call to get the basket
	suite.datasourceMock.(*mocks.DatasourceMock).AssertNotCalled(suite.T(), "GetBasket", mock.AnythingOfType("string"))
}

func (suite *CheckoutControllerTestSuite) TestAddProductWrongPayload() {
	// Given
	productCode := "FAKE"
	suite.datasourceMock.(*mocks.DatasourceMock).On("GetProduct",
		mock.AnythingOfType("model.ProductCode")).Return(*new(model.Product), errors.NewProductNotFound(productCode))

	// When
	reqBodyBytes := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBytes).Encode(requests.AddItemRequest{Code: model.ProductCode("")})
	if err != nil {
		suite.T().Errorf("Error encoding request: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("/baskets/%v", uuid.New().String()), bytes.NewBuffer(reqBodyBytes.Bytes()))
	if err != nil {
		suite.T().Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := logging.AccessLoggingMiddleware(suite.checkoutController.AddItem())

	handler.ServeHTTP(rr, req)

	// Then
	suite.Equal(http.StatusUnprocessableEntity, rr.Code)

	//Checking there has not been any call to get the basket
	suite.datasourceMock.(*mocks.DatasourceMock).AssertNotCalled(suite.T(), "GetBasket", mock.AnythingOfType("string"))
}

func (suite *CheckoutControllerTestSuite) TestAddProductToNonExistingBasket() {
	// Given
	productCode := "FAKE"
	basketId := uuid.New().String()

	suite.datasourceMock.(*mocks.DatasourceMock).On("GetProduct",
		mock.AnythingOfType("model.ProductCode")).Return(*new(model.Product), nil)
	suite.datasourceMock.(*mocks.DatasourceMock).On("GetBasket",
		mock.AnythingOfType("string")).Return(new(model.Basket), errors.NewBasketNotFound(basketId))

	// When
	reqBodyBytes := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBytes).Encode(requests.AddItemRequest{Code: model.ProductCode(productCode)})
	if err != nil {
		suite.T().Errorf("Error encoding request: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("/baskets/%v", uuid.New().String()), bytes.NewBuffer(reqBodyBytes.Bytes()))
	if err != nil {
		suite.T().Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := logging.AccessLoggingMiddleware(suite.checkoutController.AddItem())

	handler.ServeHTTP(rr, req)

	// Then
	suite.Equal(http.StatusNotFound, rr.Code)
}

func (suite *CheckoutControllerTestSuite) TestAddProduct() {
	// Given
	basketId := uuid.New().String()
	var productCode model.ProductCode = "P1"
	product := model.Product{Code: productCode, Name: "Prod 1", Price: 1000}

	suite.datasourceMock.(*mocks.DatasourceMock).On("GetProduct",
		mock.AnythingOfType("model.ProductCode")).Return(product, nil)
	suite.datasourceMock.(*mocks.DatasourceMock).On("GetBasket",
		mock.AnythingOfType("string")).Return(model.NewBasket(basketId), nil)

	// When
	reqBodyBytes := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBytes).Encode(requests.AddItemRequest{Code: productCode})
	if err != nil {
		suite.T().Errorf("Error encoding request: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("/baskets/%v", uuid.New().String()), bytes.NewBuffer(reqBodyBytes.Bytes()))
	if err != nil {
		suite.T().Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := logging.AccessLoggingMiddleware(suite.checkoutController.AddItem())

	handler.ServeHTTP(rr, req)

	// Then
	suite.Equal(http.StatusCreated, rr.Code)
}

func (suite *CheckoutControllerTestSuite) TestGetPriceNonExistingBasket() {
	// Given
	basketId := uuid.New().String()

	suite.datasourceMock.(*mocks.DatasourceMock).On("GetBasket",
		mock.AnythingOfType("string")).Return(new(model.Basket), errors.NewBasketNotFound(basketId))

	// When
	req, err := http.NewRequest("GET", fmt.Sprintf("/baskets/%s?price", basketId), nil)
	if err != nil {
		suite.T().Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := logging.AccessLoggingMiddleware(suite.checkoutController.GetPrice())

	handler.ServeHTTP(rr, req)

	// Then
	suite.Equal(http.StatusNotFound, rr.Code)
}

func (suite *CheckoutControllerTestSuite) TestGetPriceEmptyBasket() {
	// Given
	basketId := uuid.New().String()
	promotions := []model.Promotion{model.NewBulkPromotion(map[model.ProductCode][]model.BulkOfferRule{"P1": {{Buy: 3, Price: 900}}}),
		model.NewFreeItemsPromotion(map[model.ProductCode][]model.FreeItemsOfferRule{"P2": {{Buy: 3, Free: 1}}})}

	suite.datasourceMock.(*mocks.DatasourceMock).On("GetBasket",
		mock.AnythingOfType("string")).Return(model.NewBasket(basketId), nil)
	suite.datasourceMock.(*mocks.DatasourceMock).On("GetPromotions").Return(promotions)

	// When
	req, err := http.NewRequest("GET", fmt.Sprintf("/baskets/%s?price", basketId), nil)
	if err != nil {
		suite.T().Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := logging.AccessLoggingMiddleware(suite.checkoutController.GetPrice())

	handler.ServeHTTP(rr, req)

	// Then
	suite.Equal(http.StatusOK, rr.Code)

	var pbr = new(responses.PriceBasketResponse)
	err = json.Unmarshal(rr.Body.Bytes(), &pbr)

	if err != nil {
		suite.T().Errorf("Error unmarshalling basket price response: %v", err)
	}

	suite.Equal(float64(0), pbr.Total)
}

func (suite *CheckoutControllerTestSuite) TestDeleteNonExistingBasket() {
	// Given
	basketId := uuid.New().String()

	suite.datasourceMock.(*mocks.DatasourceMock).On("DeleteBasket", mock.AnythingOfType("string"))

	// When
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/baskets/%s/", basketId), nil)
	if err != nil {
		suite.T().Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := logging.AccessLoggingMiddleware(suite.checkoutController.DeleteBasket())

	handler.ServeHTTP(rr, req)

	// Then
	suite.Equal(http.StatusNoContent, rr.Code)
}

func (suite *CheckoutControllerTestSuite) TestDeleteBasket() {
	// Given
	basketId := uuid.New().String()

	suite.datasourceMock.(*mocks.DatasourceMock).On("DeleteBasket", mock.AnythingOfType("string"))

	// When
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/baskets/%s", basketId), nil)
	if err != nil {
		suite.T().Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := logging.AccessLoggingMiddleware(suite.checkoutController.DeleteBasket())

	handler.ServeHTTP(rr, req)

	// Then
	suite.Equal(http.StatusNoContent, rr.Code)
}
