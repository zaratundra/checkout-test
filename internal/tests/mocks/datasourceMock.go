package mocks

import (
	"github.com/alfcope/checkouttest/model"
	"github.com/stretchr/testify/mock"
)

type DatasourceMock struct {
	mock.Mock
}

func NewDatasourceMock() *DatasourceMock {
	return &DatasourceMock{}
}

func (d *DatasourceMock) GetProduct(code model.ProductCode) (model.Product, error) {
	args := d.Called(code)

	var err error
	if args.Get(1) == nil {
		err = nil
	} else {
		err = args.Get(1).(error)
	}

	return args.Get(0).(model.Product), err
}

func (d *DatasourceMock) GetPromotions() []model.Promotion {
	args := d.Called()

	return args.Get(0).([]model.Promotion)
}

func (d *DatasourceMock) GetBasket(id string) (*model.Basket, error) {
	args := d.Called(id)

	var err error
	if args.Get(1) == nil {
		err = nil
	} else {
		err = args.Get(1).(error)
	}

	return args.Get(0).(*model.Basket), err
}

func (d *DatasourceMock) AddBasket(basket *model.Basket) error {
	args := d.Called(basket)

	var err error
	if args.Get(0) == nil {
		err = nil
	} else {
		err = args.Get(0).(error)
	}

	return err
}

func (d *DatasourceMock) DeleteBasket(basketId string) {
	d.Called(basketId)
}
