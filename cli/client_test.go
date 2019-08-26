package cli

import (
	"fmt"
	"github.com/alfcope/checkouttest/api/responses"
	"github.com/alfcope/checkouttest/errors"
	"github.com/alfcope/checkouttest/internal/tests/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type CheckoutClientTestSuite struct {
	suite.Suite

	client *CheckoutClient
	server *mocks.CheckoutServerStub
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(CheckoutClientTestSuite))
}

func (suite *CheckoutClientTestSuite) SetupSuite() {
	suite.server = mocks.NewCheckServerStub("/api/v1")
	suite.client = NewCheckoutClient(suite.server.GetUrl(), 1)
}

func (suite *CheckoutClientTestSuite) TearDownTest() {
	suite.server.StubResponse(0, nil)
}

func (suite *CheckoutClientTestSuite) TestCreateBasketDuplicatedId() {
	// Given
	suite.server.StubResponse(responses.GetStatusByError(errors.NewPrimaryKeyError(uuid.New().String())), nil)

	// When
	b, err := suite.client.AddBasket()

	// Then
	suite.Equal("", b)
	suite.Equal(fmt.Sprintf("%d %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)), err.Error())
}

func (suite *CheckoutClientTestSuite) TestCreateBasket() {
	// Given
	basketId := uuid.New().String()
	suite.server.StubResponse(http.StatusCreated, responses.NewBasketResponse{Id: basketId})

	// When
	idResponse, err := suite.client.AddBasket()

	// Then
	suite.Nil(err)
	suite.Equal(basketId, idResponse)
}

func (suite *CheckoutClientTestSuite) TestAddItemBasketEmptyBasketId() {
	// Given
	basketId := "    "
	productCode := "TSHIRT"

	// When
	err := suite.client.AddItem(basketId, productCode)

	// Then
	suite.EqualError(err, "invalid request")
}

func (suite *CheckoutClientTestSuite) TestAddItemBasketEmptyProductCode() {
	// Given
	basketId := uuid.New().String()
	productCode := "    "

	// When
	err := suite.client.AddItem(basketId, productCode)

	// Then
	suite.EqualError(err, "invalid request")
}

func (suite *CheckoutClientTestSuite) TestAddItemBasketNotFoundError() {
	// Given
	basketId := uuid.New().String()
	productCode := "TSHIRT"
	suite.server.StubResponse(http.StatusNotFound, nil)

	// When
	err := suite.client.AddItem(basketId, productCode)

	// Then
	suite.EqualError(err, fmt.Sprintf("%d %s", http.StatusNotFound, http.StatusText(http.StatusNotFound)))
}

func (suite *CheckoutClientTestSuite) TestAddItemBasket() {
	// Given
	basketId := uuid.New().String()
	productCode := "TSHIRT"
	suite.server.StubResponse(http.StatusCreated, nil)

	// When
	err := suite.client.AddItem(basketId, productCode)

	// Then
	suite.Nil(err)
}

func (suite *CheckoutClientTestSuite) TestGetBasketPriceNotFoundError() {
	// Given
	suite.server.StubResponse(http.StatusNotFound, nil)

	// When
	price, err := suite.client.GetPrice(uuid.New().String())

	// Then
	suite.Equal(float64(-1), price)
	suite.NotNil(err)
	suite.EqualError(err, fmt.Sprintf("%d %s", http.StatusNotFound, http.StatusText(http.StatusNotFound)))
}

func (suite *CheckoutClientTestSuite) TestGetBasketPrice() {
	// Given
	suite.server.StubResponse(http.StatusOK, responses.PriceBasketResponse{Total: float64(6580) / 100})

	// When
	price, err := suite.client.GetPrice(uuid.New().String())

	// Then
	suite.Nil(err)
	suite.Equal(float64(6580)/100, price)
}

func (suite *CheckoutClientTestSuite) TestDeleteBasketNotFoundError() {
	// Given
	suite.server.StubResponse(http.StatusNotFound, nil)

	// When
	err := suite.client.DeleteBasket(uuid.New().String())

	// Then
	suite.EqualError(err, fmt.Sprintf("%d %s", http.StatusNotFound, http.StatusText(http.StatusNotFound)))
}

func (suite *CheckoutClientTestSuite) TestDeleteBasket() {
	// Given
	suite.server.StubResponse(http.StatusNoContent, nil)

	// When
	err := suite.client.DeleteBasket(uuid.New().String())

	// Then
	suite.Nil(err)
}
