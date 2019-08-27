package model

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
	for pCode, rules := range b.offers {
		if line, ok := lines[pCode]; ok {

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
						inOffer[pCode] = &[]int{}
					}

					for i := 0; i < amountAvailable; i++ {
						*inOffer[pCode] = append(*inOffer[pCode], rule.Price)
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
	for pCode, rules := range f.offers {
		if line, ok := lines[pCode]; ok {

			for _, rule := range rules {
				amountAvailable := line.amount
				alreadyInOffer, ok := inOffer[pCode]

				if ok {
					amountAvailable = amountAvailable - len(*alreadyInOffer)
				}

				promotions := amountAvailable / rule.Buy
				if promotions > 0 {
					elements := promotions * rule.Buy

					if !ok || alreadyInOffer == nil {
						inOffer[pCode] = &[]int{}
					}

					for i := 0; i < elements; i++ {
						if i < rule.Free*promotions {
							*inOffer[pCode] = append(*inOffer[pCode], 0)
						} else {
							*inOffer[pCode] = append(*inOffer[pCode], line.Product.Price)
						}
					}
				}
			}
		}
	}
}
