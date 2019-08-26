package model

import (
	"log"
)

type PromotionType string

type Promotion interface {
	GetType() PromotionType
	Resolve(map[ProductCode]Line, map[ProductCode]*[]int)
}

type BulkPromotion struct {
	//A map in case different bulk promotions are defined for different products
	//Key: ProductCode
	//Value: different possible conditions by product, for example => 3 - $19 | 5 - $15
	offers map[ProductCode][]BulkOfferRule
}

type BulkOfferRule struct {
	Buy   int
	Price int
}

func NewBulkPromotion(offers map[ProductCode][]BulkOfferRule) *BulkPromotion {
	return &BulkPromotion{
		offers: offers,
	}
}

func (b BulkPromotion) GetType() PromotionType {
	return "BULK"
}

func (b BulkPromotion) Resolve(lines map[ProductCode]Line, inOffer map[ProductCode]*[]int) {
	log.Printf("--- Bulk Promotion ---")
	for pCode, rules := range b.offers {
		log.Println("\tProduct: ", pCode)

		if line, ok := lines[pCode]; ok {
			log.Printf("\tFound %v in the basket\n", line.amount)

			for _, rule := range rules {
				amountAvailable := line.amount
				alreadyInOffer, ok := inOffer[pCode]

				if ok {
					amountAvailable = amountAvailable - len(*alreadyInOffer)
				}

				//promotions := amountAvailable / rule.Buy
				//log.Printf("\tEnough items for %v rule %v\n", promotions, rule)
				if amountAvailable >= rule.Buy {
					//elements := promotions * rule.Buy

					if !ok || alreadyInOffer == nil {
						log.Printf("\tCreating offer slice for product: %v\n", pCode)
						inOffer[pCode] = &[]int{}
					}

					for i := 0; i < amountAvailable; i++ {
						*inOffer[pCode] = append(*inOffer[pCode], rule.Price)
						log.Printf("\tinOffer length: %v\n", len(*inOffer[pCode]))

						log.Printf("\t%v elements remaining\n", amountAvailable-i-1)
					}
				}
			}
		}
	}
}

type FreeItemsPromotion struct {
	//A map in case different bulk promotions are defined for different products
	//Key: ProductCode
	//Value: slice with potentially different combinations of buy X get Y free
	offers map[ProductCode][]FreeItemsOfferRule
}

type FreeItemsOfferRule struct {
	Buy  int
	Free int
}

func NewFreeItemsPromotion(offers map[ProductCode][]FreeItemsOfferRule) *FreeItemsPromotion {
	return &FreeItemsPromotion{offers: offers}
}

func (f FreeItemsPromotion) GetType() PromotionType {
	return "FREE_ITEMS"
}

func (f FreeItemsPromotion) Resolve(lines map[ProductCode]Line, inOffer map[ProductCode]*[]int) {
	log.Printf("--- Free Items Promotion ---")
	for pCode, rules := range f.offers {
		log.Println("\tProduct: ", pCode)

		if line, ok := lines[pCode]; ok {
			log.Printf("\tFound %v in the basket\n", line.amount)

			for _, rule := range rules {
				amountAvailable := line.amount
				alreadyInOffer, ok := inOffer[pCode]

				if ok {
					amountAvailable = amountAvailable - len(*alreadyInOffer)
				}

				promotions := amountAvailable / rule.Buy
				log.Printf("\tEnough items for %v rule %v\n", promotions, rule)
				if promotions > 0 {
					elements := promotions * rule.Buy

					if !ok || alreadyInOffer == nil {
						log.Printf("\tCreating offer slice for product: %v\n", pCode)
						inOffer[pCode] = &[]int{}
					}

					for i := 0; i < elements; i++ {
						if i < rule.Free*promotions {
							*inOffer[pCode] = append(*inOffer[pCode], 0)
						} else {
							*inOffer[pCode] = append(*inOffer[pCode], line.Product.Price)
						}

						log.Printf("\tinOffer length: %v\n", len(*inOffer[pCode]))

						log.Printf("\t%v elements remaining\n", elements-i-1)
					}
				}
			}
		}
	}
}
