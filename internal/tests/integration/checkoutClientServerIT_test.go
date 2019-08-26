package integration

import (
	"fmt"
	"github.com/alfcope/checkouttest/cli"
	"github.com/stretchr/testify/suite"
	"net/http"
	"regexp"
	"testing"
)

func isUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

type CheckoutServiceClientITSuite struct {
	suite.Suite
	client *cli.CheckoutClient
}

func TestCheckoutServiceClientIT(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping ITs in short mode.")
	}

	suite.Run(t, new(CheckoutServiceClientITSuite))
}

func (suite *CheckoutServiceClientITSuite) SetupSuite() {
	suite.client = cli.NewCheckoutClient("http://localhost:7070", 1)
}

func (suite *CheckoutServiceClientITSuite) TestAddBasket() {
	id, err := suite.client.AddBasket()

	suite.Nil(err)
	suite.True(isUUID(id))
}

func (suite *CheckoutServiceClientITSuite) TestAddNonExistingProduct() {
	id, err := suite.client.AddBasket()
	if err != nil {
		suite.T().Errorf("error creating basket: %v", err.Error())
	}

	err = suite.client.AddItem(id, "FAKE")

	suite.EqualError(err, fmt.Sprintf("%d %s", http.StatusNotFound, http.StatusText(http.StatusNotFound)))
}

func (suite *CheckoutServiceClientITSuite) TestAddProductMultipleTimes() {
	id, err := suite.client.AddBasket()
	if err != nil {
		suite.T().Errorf("error creating basket: %v", err.Error())
	}

	for i := 0; i < 5; i++ {
		err = suite.client.AddItem(id, "VOUCHER")
		if err != nil {
			suite.T().Errorf("error adding product: %v", err.Error())
		}
	}
}

func (suite *CheckoutServiceClientITSuite) TestGetPrice() {
	id, err := suite.client.AddBasket()
	if err != nil {
		suite.T().Errorf("error creating basket: %v", err.Error())
	}

	products := []string{"VOUCHER", "TSHIRT", "MUG"}

	for _, product := range products {
		err = suite.client.AddItem(id, product)
		if err != nil {
			suite.T().Errorf("error adding product: %v", err.Error())
		}
	}

	price, err := suite.client.GetPrice(id)

	suite.Nil(err)
	suite.True(float64(3250)/100 == price)

	// With promotions
	products = []string{"VOUCHER", "VOUCHER", "TSHIRT", "TSHIRT"}

	for _, product := range products {
		err = suite.client.AddItem(id, product)
		if err != nil {
			suite.T().Errorf("error adding product: %v", err.Error())
		}
	}

	price, err = suite.client.GetPrice(id)

	suite.Nil(err)
	suite.True(float64(7450)/100 == price)
}

func (suite *CheckoutServiceClientITSuite) TestDeleteBasket() {
	id, err := suite.client.AddBasket()
	if err != nil {
		suite.T().Errorf("error creating basket: %v", err.Error())
	}

	err = suite.client.DeleteBasket(id)

	suite.Nil(err)

	err = suite.client.AddItem(id, "VOUCHER")
	suite.EqualError(err, fmt.Sprintf("%d %s", http.StatusNotFound, http.StatusText(http.StatusNotFound)))
}
