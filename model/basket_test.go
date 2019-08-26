package model

import (
	"fmt"
	"github.com/alfcope/checkouttest/errors"
	"github.com/google/uuid"
	"testing"
)

// Adding invalid product
func TestAddFirstProduct(t *testing.T) {
	basket := NewBasket(uuid.New().String())

	err := basket.AddProduct(Product{"P1", "Product 1", -10})
	if err != nil {
		if _, ok := err.(*errors.ValidationError); !ok {
			t.Errorf("Expected validation error but got %T", err)
		}
	} else {
		t.Errorf("Error expected but did not get one")
	}

	if len(basket.lines) > 0 {
		t.Errorf("There should not be any line")
	}
}

// Adding a  product
func TestAddProduct(t *testing.T) {
	basket := NewBasket(uuid.New().String())

	err := basket.AddProduct(Product{"P1", "Product 1", 800})
	if err != nil {
		t.Error("Unexpected error ", err.Error())
	}

	if len(basket.lines) != 1 {
		t.Error("There should be just one line")
	}

	if line, ok := basket.lines["P1"]; ok {
		if line.amount != 1 {
			t.Errorf("Got amount %v when wanted 1", line.amount)
		}
	} else {
		t.Error("Product line not found")
	}
}

// Adding same product multiple times
func TestAddProductMultipleTimes(t *testing.T) {
	basket := NewBasket(uuid.New().String())
	var times = 3

	for i := 0; i < times; i++ {
		err := basket.AddProduct(Product{"P1", "Product 1", 800})
		if err != nil {
			t.Error("Unexpected error ", err.Error())
		}
	}

	if len(basket.lines) != 1 {
		t.Error("There should be just one line")
	}

	if line, ok := basket.lines["P1"]; ok {
		if line.amount != times {
			t.Errorf("Got amount %v when wanted %v", line.amount, times)
		}
	} else {
		t.Error("Product line not found")
	}
}

// Adding multiple products
func TestAddMultipleProducts(t *testing.T) {
	basket := NewBasket(uuid.New().String())

	for i := 1; i < 4; i++ {
		err := basket.AddProduct(Product{ProductCode(fmt.Sprintf("P%d", i)),
			fmt.Sprintf("Product %d", i), 100 * i})
		if err != nil {
			t.Error("Unexpected error ", err.Error())
		}
	}

	if len(basket.lines) != 3 {
		t.Error("There should be 3 lines")
	}

	for i := 1; i < 4; i++ {
		if _, ok := basket.lines[ProductCode(fmt.Sprintf("P%d", i))]; !ok {
			t.Error("Product line not found")
		}
	}
}

var basketPriceCases = []struct {
	lines  map[ProductCode]Line
	offers []Promotion
	price  float64
}{
	{ // No active offers
		map[ProductCode]Line{"P1": {Product{"P1", "Prod name 1", 1000}, 3}},
		[]Promotion{},
		float64(1000*3) / 100,
	}, { // Empty basket
		map[ProductCode]Line{},
		[]Promotion{NewBulkPromotion(map[ProductCode][]BulkOfferRule{"P1": {{3, 820}}})},
		float64(0),
	}, { // Basket without any products in offer
		map[ProductCode]Line{"P1": {Product{"P1", "Prod name 1", 1000}, 3}},
		[]Promotion{NewBulkPromotion(map[ProductCode][]BulkOfferRule{"P2": {{3, 820}}})},
		float64(1000*3) / 100,
	}, { // Basket with all products matching an offer
		map[ProductCode]Line{"P1": {Product{"P1", "Prod name 1", 1000}, 3}},
		[]Promotion{NewBulkPromotion(map[ProductCode][]BulkOfferRule{"P1": {{3, 820}}})},
		float64(820*3) / 100,
	}, { // Basket with products matching an offer several times
		map[ProductCode]Line{"P1": {Product{"P1", "Prod name 1", 1000}, 9}},
		[]Promotion{NewBulkPromotion(map[ProductCode][]BulkOfferRule{"P1": {{3, 820}}})},
		float64(820*9) / 100,
	}, { // Basket with products matching an offer several times plus extra number
		map[ProductCode]Line{"P1": {Product{"P1", "Prod name 1", 1000}, 7}},
		[]Promotion{NewBulkPromotion(map[ProductCode][]BulkOfferRule{"P1": {{3, 820}}})},
		float64(820*7) / 100,
	}, { // Basket with same products matching different offers
		map[ProductCode]Line{"P1": {Product{"P1", "Prod name 1", 1000}, 5}},
		[]Promotion{NewBulkPromotion(map[ProductCode][]BulkOfferRule{"P1": {{3, 820}, {2, 930}}})},
		float64(820*5) / 100,
	}, { // Basket with different products matching different offers
		map[ProductCode]Line{"P1": {Product{"P1", "Prod name 1", 1030}, 3},
			"P2": {Product{"P2", "Prod name 2", 1545}, 3}},
		[]Promotion{NewBulkPromotion(map[ProductCode][]BulkOfferRule{"P1": {{3, 900}}}),
			NewFreeItemsPromotion(map[ProductCode][]FreeItemsOfferRule{"P2": {{3, 1}}})},
		float64(900*3+1545*2) / 100,
	}, { // Basket with different products matching same offer with rules for that products
		map[ProductCode]Line{"P1": {Product{"P1", "Prod name 1", 1030}, 3},
			"P2": {Product{"P2", "Prod name 2", 1545}, 4}},
		[]Promotion{NewBulkPromotion(map[ProductCode][]BulkOfferRule{"P1": {{3, 900}}, "P2": {{3, 1210}}})},
		float64(900*3+1210*4) / 100,
	}, { // Basket with different products matching same offer with rules for that products
		map[ProductCode]Line{"P1": {Product{"P1", "Prod name 1", 500}, 3},
			"P2": {Product{"P2", "Prod name 2", 2000}, 3},
			"P3": {Product{"P3", "Prod name 3", 750}, 1}},
		[]Promotion{NewBulkPromotion(map[ProductCode][]BulkOfferRule{"P2": {{3, 1900}}}),
			NewFreeItemsPromotion(map[ProductCode][]FreeItemsOfferRule{"P1": {{2, 1}}})},
		float64(500*2+1900*3+750) / 100,
	},
}

func TestBasketPrices(t *testing.T) {
	for _, tb := range basketPriceCases {

		basket := NewBasket(uuid.New().String())
		basket.lines = tb.lines

		p := basket.CalculatePrice(tb.offers)
		fmt.Printf(" ------------ Price: %v\n", p)
		if p != tb.price {
			t.Errorf("Wanted %v but got %v", tb.price, p)
		}
	}
}
