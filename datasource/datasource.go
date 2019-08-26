package datasource

import (
	"encoding/json"
	"github.com/alfcope/checkouttest/config"
	"github.com/alfcope/checkouttest/datasource/parser"
	"github.com/alfcope/checkouttest/errors"
	"github.com/alfcope/checkouttest/model"
	"io/ioutil"
	"sync"
)

type Datasource interface {
	GetProduct(model.ProductCode) (model.Product, error)
	GetPromotions() []model.Promotion
	GetBasket(string) (*model.Basket, error)
	AddBasket(*model.Basket) error
	DeleteBasket(string)
}

type InMemoryDatasource struct {
	// products and promotions do not need mutex as they do not
	// change its state. Just once at startup
	products   map[model.ProductCode]model.Product
	promotions []model.Promotion

	baskets    map[string]*model.Basket
	basketsMux sync.RWMutex
}

func InitInMemoryDatasource(config config.DataConfig) (*InMemoryDatasource, error) {
	ds := InMemoryDatasource{
		products:   make(map[model.ProductCode]model.Product),
		promotions: make([]model.Promotion, 0),
		baskets:    make(map[string]*model.Basket),
		basketsMux: sync.RWMutex{},
	}

	err := ds.loadProducts(config.Products)
	if err != nil {
		return nil, err
	}

	err = ds.loadPromotions(config.Promotions)
	if err != nil {
		return nil, err
	}

	return &ds, nil
}

func (d *InMemoryDatasource) GetProduct(code model.ProductCode) (model.Product, error) {
	if product, ok := d.products[code]; ok {
		return product, nil
	}

	return *new(model.Product), errors.NewProductNotFound(string(code))
}

func (d *InMemoryDatasource) GetPromotions() []model.Promotion {
	return d.promotions[:]
}

func (d *InMemoryDatasource) GetBasket(id string) (*model.Basket, error) {
	if basket, ok := d.baskets[id]; ok {
		return basket, nil
	}

	return new(model.Basket), errors.NewBasketNotFound(id)
}

func (d *InMemoryDatasource) AddBasket(basket *model.Basket) error {
	d.basketsMux.Lock()
	defer d.basketsMux.Unlock()

	if _, ok := d.baskets[basket.Id]; !ok {
		d.baskets[basket.Id] = basket
		return nil
	}

	return errors.NewPrimaryKeyError(basket.Id)
}

func (d *InMemoryDatasource) DeleteBasket(basketId string) {
	d.basketsMux.Lock()
	defer d.basketsMux.Unlock()

	delete(d.baskets, basketId)
}

func (d *InMemoryDatasource) loadProducts(filePath string) error {
	var products []model.Product

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &products)
	if err != nil {
		return err
	}

	for _, p := range products {
		err := p.Validate()
		if err == nil {
			d.products[p.Code] = p
		}
	}

	return nil
}

func (d *InMemoryDatasource) loadPromotions(filePath string) error {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	var nodes []map[string]interface{}
	err = json.Unmarshal(file, &nodes)
	if err != nil {
		return err
	}

	for _, promotionNode := range nodes {
		promotion, err := parser.ParsePromotion(promotionNode)
		if err != nil {
			if _, ok := err.(*errors.PromotionNotFound); !ok {
				return err
			}
			continue
		}

		d.promotions = append(d.promotions, promotion)
	}

	return nil
}
