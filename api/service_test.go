package api

import (
	"github.com/alfcope/checkouttest/datasource"
	"github.com/alfcope/checkouttest/errors"
	"github.com/alfcope/checkouttest/internal/tests/mocks"
	"github.com/alfcope/checkouttest/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type CheckoutServiceTestSuite struct {
	suite.Suite

	datasourceMock  datasource.Datasource
	checkoutService CheckoutService
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(CheckoutServiceTestSuite))
}

func (suite *CheckoutServiceTestSuite) SetupSuite() {
	suite.datasourceMock = datasource.Datasource(mocks.NewDatasourceMock())
	suite.checkoutService = NewCheckoutService(suite.datasourceMock)
}

func (suite *CheckoutServiceTestSuite) TearDownTest() {
	suite.datasourceMock.(*mocks.DatasourceMock).ExpectedCalls = nil
	suite.datasourceMock.(*mocks.DatasourceMock).Calls = nil
}

func (suite *CheckoutServiceTestSuite) TestCreateBasketDuplicatedId() {
	// Given
	basketId := uuid.New().String()
	suite.datasourceMock.(*mocks.DatasourceMock).On("AddBasket", mock.AnythingOfType("*model.Basket")).Return(errors.NewPrimaryKeyError(basketId))

	// When
	b, err := suite.checkoutService.CreateBasket()

	// Then
	suite.Equal("", b)
	if primaryKeyError, ok := err.(*errors.PrimaryKeyError); ok {
		suite.Equal(basketId, primaryKeyError.Id)
	} else {
		suite.T().Error("Error should be a primary key error")
	}
}

func (suite *CheckoutServiceTestSuite) TestCreateBasket() {
	// Given
	suite.datasourceMock.(*mocks.DatasourceMock).On("AddBasket", mock.AnythingOfType("*model.Basket")).Return(nil)

	// When
	b, err := suite.checkoutService.CreateBasket()

	// Then
	suite.NotEqual("", b)
	suite.Nil(err)
}

func (suite *CheckoutServiceTestSuite) TestAddNonExistingProduct() {
	// Given
	productCode := "FAKE"
	suite.datasourceMock.(*mocks.DatasourceMock).On("GetProduct",
		mock.AnythingOfType("model.ProductCode")).Return(*new(model.Product), errors.NewProductNotFound(productCode))

	// When
	err := suite.checkoutService.AddProduct(uuid.New().String(), model.ProductCode(productCode))

	// Then
	if productNotFound, ok := err.(*errors.ProductNotFound); ok {
		suite.Equal(productCode, productNotFound.Code)
	} else {
		suite.T().Error("Error should be a product not found error ")
	}

	//Checking there has not been any call to get the basket
	suite.datasourceMock.(*mocks.DatasourceMock).AssertNotCalled(suite.T(), "GetBasket", mock.AnythingOfType("string"))
}

func (suite *CheckoutServiceTestSuite) TestAddProductToNonExistingBasket() {
	// Given
	basketId := uuid.New().String()
	var productCode model.ProductCode = "P1"
	product := model.Product{Code: productCode, Name: "Prod 1", Price: 1000}

	suite.datasourceMock.(*mocks.DatasourceMock).On("GetProduct",
		mock.AnythingOfType("model.ProductCode")).Return(product, nil)
	suite.datasourceMock.(*mocks.DatasourceMock).On("GetBasket",
		mock.AnythingOfType("string")).Return(new(model.Basket), errors.NewBasketNotFound(basketId))

	// When
	err := suite.checkoutService.AddProduct(uuid.New().String(), productCode)

	// Then
	if basketNotFound, ok := err.(*errors.BasketNotFound); ok {
		suite.Equal(basketId, basketNotFound.Id)
	} else {
		suite.T().Error("Error should be a basket not found error ")
	}
}

func (suite *CheckoutServiceTestSuite) TestAddProduct() {
	// Given
	basketId := uuid.New().String()
	var productCode model.ProductCode = "P1"
	product := model.Product{Code: productCode, Name: "Prod 1", Price: 1000}

	suite.datasourceMock.(*mocks.DatasourceMock).On("GetProduct",
		mock.AnythingOfType("model.ProductCode")).Return(product, nil)
	suite.datasourceMock.(*mocks.DatasourceMock).On("GetBasket",
		mock.AnythingOfType("string")).Return(model.NewBasket(basketId), nil)

	// When
	err := suite.checkoutService.AddProduct(uuid.New().String(), productCode)

	// Then
	suite.Nil(err)
}

func (suite *CheckoutServiceTestSuite) TestGetPriceNonExistingBasket() {
	// Given
	basketId := uuid.New().String()

	suite.datasourceMock.(*mocks.DatasourceMock).On("GetBasket",
		mock.AnythingOfType("string")).Return(new(model.Basket), errors.NewBasketNotFound(basketId))

	// When
	price, err := suite.checkoutService.GetBasketPrice(uuid.New().String())

	// Then
	if basketNotFound, ok := err.(*errors.BasketNotFound); ok {
		suite.Equal(basketId, basketNotFound.Id)
	} else {
		suite.T().Error("Error should be a basket not found error ")
	}
	suite.Equal(float64(0), price)
}

func (suite *CheckoutServiceTestSuite) TestGetPriceEmptyBasket() {
	// Given
	basketId := uuid.New().String()
	promotions := []model.Promotion{model.NewBulkPromotion(map[model.ProductCode][]model.BulkOfferRule{"P1": {{Buy: 3, Price: 900}}}),
		model.NewFreeItemsPromotion(map[model.ProductCode][]model.FreeItemsOfferRule{"P2": {{Buy: 3, Free: 1}}})}

	suite.datasourceMock.(*mocks.DatasourceMock).On("GetBasket",
		mock.AnythingOfType("string")).Return(model.NewBasket(basketId), nil)
	suite.datasourceMock.(*mocks.DatasourceMock).On("GetPromotions").Return(promotions)

	// When
	price, err := suite.checkoutService.GetBasketPrice(uuid.New().String())

	// Then
	suite.Nil(err)
	suite.Equal(float64(0), price)
}
