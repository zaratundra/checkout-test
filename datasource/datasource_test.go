package datasource

import (
	"github.com/alfcope/checkouttest/config"
	"github.com/alfcope/checkouttest/errors"
	"github.com/alfcope/checkouttest/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"testing"
)

type DatasourceTestSuite struct {
	suite.Suite
	inMemoryDatasource InMemoryDatasource
}

func TestDatasourceTestSuite(t *testing.T) {
	suite.Run(t, new(DatasourceTestSuite))
}

func (suite *DatasourceTestSuite) SetupTest() {
	suite.inMemoryDatasource = *suite.initializeDataSource()
}

func (suite *DatasourceTestSuite) initializeDataSource() *InMemoryDatasource {
	configuration, err := config.LoadConfiguration("../internal/tests/config", "service_config_test")
	if err != nil {
		suite.T().Errorf("Error loading configuration: %v", err.Error())
	}

	inMemoryDatasource, err := InitInMemoryDatasource(configuration.Data)
	if err != nil {
		suite.T().Errorf("Error initializing datasource: %s", err.Error())
	}
	suite.Equal(3, len(inMemoryDatasource.products))
	suite.Equal(2, len(inMemoryDatasource.promotions))

	return inMemoryDatasource
}

func (suite *DatasourceTestSuite) TestInMemoryDatasource_GetNonExistingProduct() {
	// Given
	var fakeProductCode model.ProductCode = "FAKE"

	// When
	_, err := suite.inMemoryDatasource.GetProduct(fakeProductCode)

	// Then
	suite.NotNil(err)
	if pnf, ok := err.(*errors.ProductNotFound); ok {
		suite.Equal(fakeProductCode, model.ProductCode(pnf.Code))
	} else {
		suite.T().Errorf("Wanted product not found error, got %T", err)
	}
}

func (suite *DatasourceTestSuite) TestInMemoryDatasource_GetProduct() {
	// Given
	var fakeProductCode model.ProductCode = "TSHIRT"

	// When
	p, err := suite.inMemoryDatasource.GetProduct(fakeProductCode)

	// Then
	suite.Nil(err)
	suite.Equal(fakeProductCode, model.ProductCode(p.Code))
	suite.Equal(2000, p.Price)
	suite.Equal("Cabify T-Shirt", p.Name)
}

func (suite *DatasourceTestSuite) TestInMemoryDatasource_GetPromotions() {
	// Given

	// When
	p := suite.inMemoryDatasource.GetPromotions()

	// Then
	suite.Equal(2, len(p))
	suite.Equal(model.PromotionType("BULK"), p[0].GetType())
	suite.Equal(model.PromotionType("FREE_ITEMS"), p[1].GetType())
}

func (suite *DatasourceTestSuite) TestInMemoryDatasource_GetNonExistingBasket() {
	// Given
	basketId := uuid.New().String()

	// When
	_, err := suite.inMemoryDatasource.GetBasket(basketId)

	// Then
	suite.NotNil(err)
	if bnf, ok := err.(*errors.BasketNotFound); ok {
		suite.Equal(basketId, bnf.Id)
	} else {
		suite.T().Errorf("Wanted basket not found error, got %T", err)
	}
}

func (suite *DatasourceTestSuite) TestInMemoryDatasource_GetBasket() {
	// Given
	// Not using the in-memory datasource from the suite to avoid concurrency errors
	inMemoryDatasource := suite.initializeDataSource()
	basket := model.NewBasket(uuid.New().String())
	inMemoryDatasource.baskets = map[string]*model.Basket{basket.Id: basket}

	// When
	b, err := inMemoryDatasource.GetBasket(basket.Id)

	// Then
	suite.Nil(err)
	suite.Equal(basket, b)
}

func (suite *DatasourceTestSuite) TestInMemoryDatasource_AddBasketDuplicated() {
	// Given
	// Not using the in-memory datasource from the suite to avoid concurrency errors
	inMemoryDatasource := suite.initializeDataSource()
	basket := model.NewBasket(uuid.New().String())
	inMemoryDatasource.baskets = map[string]*model.Basket{basket.Id: basket}

	// When
	err := inMemoryDatasource.AddBasket(basket)

	// Then
	suite.NotNil(err)
	if duplicatedPrimaryKey, ok := err.(*errors.PrimaryKeyError); ok {
		suite.Equal(basket.Id, duplicatedPrimaryKey.Id)
	} else {
		suite.T().Errorf("Wanted primary key error, got %T", err)
	}
	suite.Equal(1, len(inMemoryDatasource.baskets))
}

func (suite *DatasourceTestSuite) TestInMemoryDatasource_AddBasket() {
	// Given
	// Not using the in-memory datasource from the suite to avoid concurrency errors
	inMemoryDatasource := suite.initializeDataSource()
	basket := model.NewBasket(uuid.New().String())

	// When
	err := inMemoryDatasource.AddBasket(basket)

	// Then
	suite.Nil(err)
	suite.Equal(1, len(inMemoryDatasource.baskets))
}

func (suite *DatasourceTestSuite) TestInMemoryDatasource_DeleteNonExistingBasket() {
	// Given
	// Not using the in-memory datasource from the suite to avoid concurrency errors
	inMemoryDatasource := suite.initializeDataSource()
	basket := model.NewBasket(uuid.New().String())
	inMemoryDatasource.AddBasket(basket)

	// When
	inMemoryDatasource.DeleteBasket(uuid.New().String())

	// Then
	suite.Equal(1, len(inMemoryDatasource.baskets))
}

func (suite *DatasourceTestSuite) TestInMemoryDatasource_DeleteBasket() {
	// Given
	// Not using the in-memory datasource from the suite to avoid concurrency errors
	inMemoryDatasource := suite.initializeDataSource()
	basket := model.NewBasket(uuid.New().String())

	// When
	inMemoryDatasource.DeleteBasket(basket.Id)

	// Then
	suite.Equal(0, len(inMemoryDatasource.baskets))
}
