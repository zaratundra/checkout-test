package model

import (
	"testing"
)

var promotionCases = []struct {
	basketLines    map[ProductCode]Line // Items in the basket
	promo          Promotion            // Promotion to apply
	itWithoutPromo int                  // Number of items out of the promotion
	itWithPromo    int                  // Number of items eligible by the promotion
	free           int                  // Number of free items if applicable, -1 otherwise
}{
	// ----- BULK PROMOTION TESTS ------
	{ // Edge case: empty basket - without lines
		make(map[ProductCode]Line),
		NewBulkPromotion(map[ProductCode][]BulkOfferRule{"P2": {{Buy: 2, Price: 1000}}}),
		0,
		0,
		-1,
	}, { // Different products without matching any promotion
		map[ProductCode]Line{"P1": {Product: Product{Code: "P1", Name: "aaaa", Price: 1000}, amount: 1},
			"P2": {Product: Product{Code: "P2", Name: "bbbb", Price: 1200}, amount: 1},
			"P3": {Product: Product{Code: "P3", Name: "cccc", Price: 1500}, amount: 1}},
		NewBulkPromotion(map[ProductCode][]BulkOfferRule{"P2": {{Buy: 2, Price: 1000}}}),
		3,
		0,
		-1,
	}, { // Exact amount of items for a promotion
		map[ProductCode]Line{"P1": {Product: Product{Code: "P1", Name: "aaaa", Price: 1000}, amount: 3}},
		NewBulkPromotion(map[ProductCode][]BulkOfferRule{"P1": {{Buy: 3, Price: 850}}}),
		0,
		3,
		-1,
	}, { // Spare items
		map[ProductCode]Line{"P2": {Product: Product{Code: "P2", Name: "bbbb", Price: 1200}, amount: 3}},
		NewBulkPromotion(map[ProductCode][]BulkOfferRule{"P2": {{Buy: 2, Price: 1000}}}),
		0,
		3,
		-1,
	}, { // Exact amount of same items matching two different rules
		map[ProductCode]Line{"P1": {Product: Product{Code: "P1", Name: "aaaa", Price: 1000}, amount: 7}},
		NewBulkPromotion(map[ProductCode][]BulkOfferRule{"P1": {{Buy: 5, Price: 650}, {Buy: 2, Price: 850}}}),
		0,
		7,
		-1,
	}, {
		map[ProductCode]Line{"P3": {Product: Product{Code: "P3", Name: "cccc", Price: 1500}, amount: 15}},
		NewBulkPromotion(map[ProductCode][]BulkOfferRule{"P3": {{Buy: 4, Price: 1100}}}),
		0,
		15,
		-1,
	}, { // Exact amount of two different items matching two different rules
		map[ProductCode]Line{"P1": {Product: Product{Code: "P1", Name: "aaaa", Price: 1500}, amount: 3},
			"P2": {Product: Product{Code: "P2", Name: "bbbb", Price: 1200}, amount: 3}},
		NewBulkPromotion(map[ProductCode][]BulkOfferRule{"P1": {{Buy: 3, Price: 1300}},
			"P2": {{Buy: 3, Price: 1000}}}),
		0,
		6,
		-1,
	},
	// ----- FREE ITEMS PROMOTION TESTS ------
	{ // Edge case: empty basket - without lines
		make(map[ProductCode]Line),
		NewFreeItemsPromotion(map[ProductCode][]FreeItemsOfferRule{"P2": {{Buy: 2, Free: 1}}}),
		0,
		0,
		0,
	}, { // Different products without matching any promotion
		map[ProductCode]Line{"P1": {Product: Product{Code: "P1", Name: "aaaa", Price: 1000}, amount: 1},
			"P2": {Product: Product{Code: "P2", Name: "bbbb", Price: 1200}, amount: 1},
			"P3": {Product: Product{Code: "P3", Name: "cccc", Price: 1500}, amount: 1}},
		NewFreeItemsPromotion(map[ProductCode][]FreeItemsOfferRule{"P2": {{Buy: 2, Free: 1}}}),
		3,
		0,
		0,
	}, { // Exact amount of items for a promotion
		map[ProductCode]Line{"P1": {Product: Product{Code: "P1", Name: "aaaa", Price: 1000}, amount: 3}},
		NewFreeItemsPromotion(map[ProductCode][]FreeItemsOfferRule{"P1": {{Buy: 3, Free: 1}}}),
		0,
		3,
		1,
	}, {
		map[ProductCode]Line{"P2": {Product: Product{Code: "P2", Name: "bbbb", Price: 1200}, amount: 3}},
		NewFreeItemsPromotion(map[ProductCode][]FreeItemsOfferRule{"P2": {{Buy: 2, Free: 1}}}),
		1,
		2,
		1,
	}, { // Exact amount of same items matching two different rules
		map[ProductCode]Line{"P1": {Product: Product{Code: "P1", Name: "aaaa", Price: 1000}, amount: 7}},
		NewFreeItemsPromotion(map[ProductCode][]FreeItemsOfferRule{"P1": {{Buy: 5, Free: 3}, {Buy: 2, Free: 1}}}),
		0,
		7,
		4,
	}, {
		map[ProductCode]Line{"P3": {Product: Product{Code: "P3", Name: "cccc", Price: 1500}, amount: 15}},
		NewFreeItemsPromotion(map[ProductCode][]FreeItemsOfferRule{"P3": {{Buy: 4, Free: 1}}}),
		3,
		12,
		3,
	}, { // Exact amount of two different items matching two different rules for one promotion
		map[ProductCode]Line{"P1": {Product: Product{Code: "P1", Name: "aaaa", Price: 1500}, amount: 3},
			"P2": {Product: Product{Code: "P2", Name: "bbbb", Price: 1200}, amount: 3}},
		NewFreeItemsPromotion(map[ProductCode][]FreeItemsOfferRule{"P1": {{Buy: 3, Free: 1}},
			"P2": {{Buy: 3, Free: 1}}}),
		0,
		6,
		2,
	},
}

func TestPromotions(t *testing.T) {
	for _, tc := range promotionCases {
		inOffer := make(map[ProductCode]*[]int)

		tc.promo.Resolve(tc.basketLines, inOffer)

		cInOffer := 0
		cOutOffer := 0
		freeCounter := 0
		//For each product in the basket
		for pCode, line := range tc.basketLines {
			// Check how many are in/out/free offer
			if items, ok := inOffer[pCode]; ok {
				cInOffer += len(*items)
				cOutOffer += line.amount - len(*items)

				// If the test case expects free items
				if tc.free >= 0 {
					// Count number of items with price equals to zero
					for _, price := range *items {
						if price == 0 {
							freeCounter++
						}
					}
				}
			} else {
				cOutOffer += line.amount
			}
		}

		if cInOffer != tc.itWithPromo {
			t.Errorf("got %v items with promo, wanted %v", cInOffer, tc.itWithPromo)
		}
		if cOutOffer != tc.itWithoutPromo {
			t.Errorf("got %v items without promo, wanted %v", cOutOffer, tc.itWithoutPromo)
		}
		if tc.free >= 0 && freeCounter != tc.free {
			t.Errorf("got %v items free, wanted %v", freeCounter, tc.free)
		}

	}
}
