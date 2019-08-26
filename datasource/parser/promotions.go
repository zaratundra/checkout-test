package parser

import (
	"fmt"
	"github.com/alfcope/checkouttest/errors"
	"github.com/alfcope/checkouttest/model"
)

func ParsePromotion(nodes map[string]interface{}) (model.Promotion, error) {

	if _, ok := nodes["code"]; !ok {
		return nil, errors.NewPromotionNotFound("")
	}

	switch nodes["code"].(string) {
	case "BULK":
		return parseBulkPromotion(nodes)

	case "FREE_ITEMS":
		return parseFreeItemsPromotion(nodes)

	default:
		return nil, errors.NewPromotionNotFound(nodes["code"].(string))
	}
}

func parseBulkPromotion(nodes map[string]interface{}) (*model.BulkPromotion, error) {
	var promos map[model.ProductCode][]model.BulkOfferRule

	rawPromos := nodes["promos"].([]interface{})
	promos = make(map[model.ProductCode][]model.BulkOfferRule, len(rawPromos))

	for _, rawPromo := range rawPromos {
		if _, ok := rawPromo.(map[string]interface{}); !ok {
			fmt.Printf("Invalid map: %v", rawPromo)
			continue
		}
		promo := rawPromo.(map[string]interface{})

		if _, ok := promo["product"].(string); !ok {
			fmt.Printf("Invalid product code: %v %T\n", promo["product"], promo["product"])
			continue
		}

		if _, ok := promo["rules"].([]interface{}); !ok {
			fmt.Printf("Invalid offer conditions: %v %T\n", promo["rules"], promo["rules"])
			continue
		}

		for _, rawRules := range promo["rules"].([]interface{}) {
			if _, ok := rawRules.(map[string]interface{}); !ok {
				fmt.Printf("Invalid map: %v", rawRules)
				continue
			}
			rule := rawRules.(map[string]interface{})

			if _, ok := rule["buy"].(float64); !ok {
				fmt.Printf("Invalid amount to buy: %v\n", rule["buy"])
				continue
			}
			if _, ok := rule["price"].(float64); !ok {
				fmt.Printf("Invalid price: %v\n", rule["price"])
				continue
			}

			promoRule := model.BulkOfferRule{
				Buy:   int(rule["buy"].(float64)),
				Price: int(rule["price"].(float64)),
			}

			if promosProduct, ok := promos[model.ProductCode(promo["product"].(string))]; ok {
				promos[model.ProductCode(promo["product"].(string))] = append(promosProduct, promoRule)
			} else {
				promos[model.ProductCode(promo["product"].(string))] = []model.BulkOfferRule{promoRule}
			}
		}
	}

	if len(promos) == 0 {
		return nil, errors.NewPromotionInvalid(nodes["code"].(string), "empty items list")
	}

	return model.NewBulkPromotion(promos), nil
}

func parseFreeItemsPromotion(nodes map[string]interface{}) (*model.FreeItemsPromotion, error) {
	var promos map[model.ProductCode][]model.FreeItemsOfferRule

	rawPromos := nodes["promos"].([]interface{})
	promos = make(map[model.ProductCode][]model.FreeItemsOfferRule, len(rawPromos))

	for _, rawPromo := range rawPromos {
		if _, ok := rawPromo.(map[string]interface{}); !ok {
			fmt.Printf("Invalid map: %v", rawPromo)
			continue
		}
		promo := rawPromo.(map[string]interface{})

		if _, ok := promo["product"].(string); !ok {
			fmt.Printf("Invalid product code: %v %T\n", promo["product"], promo["product"])
			continue
		}

		if _, ok := promo["rules"].([]interface{}); !ok {
			fmt.Printf("Invalid offer conditions: %v %T\n", promo["rules"], promo["rules"])
			continue
		}

		for _, rawRules := range promo["rules"].([]interface{}) {
			if _, ok := rawRules.(map[string]interface{}); !ok {
				fmt.Printf("Invalid map: %v", rawRules)
				continue
			}
			rule := rawRules.(map[string]interface{})

			if _, ok := rule["buy"].(float64); !ok {
				fmt.Printf("Invalid amount to buy: %v\n", rule["buy"])
				continue
			}
			if _, ok := rule["free"].(float64); !ok {
				fmt.Printf("Invalid amount: %v\n", rule["free"])
				continue
			}

			promoRule := model.FreeItemsOfferRule{
				Buy:  int(rule["buy"].(float64)),
				Free: int(rule["free"].(float64)),
			}

			if promosProduct, ok := promos[model.ProductCode(promo["product"].(string))]; ok {
				promos[model.ProductCode(promo["product"].(string))] = append(promosProduct, promoRule)
			} else {
				promos[model.ProductCode(promo["product"].(string))] = []model.FreeItemsOfferRule{promoRule}
			}
		}
	}

	if len(promos) == 0 {
		return nil, errors.NewPromotionInvalid(nodes["code"].(string), "empty items list")
	}

	return model.NewFreeItemsPromotion(promos), nil
}
