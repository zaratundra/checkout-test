package model

import (
	"github.com/alfcope/checkouttest/errors"
	"strings"
)

type ProductCode string

type Product struct {
	Code  ProductCode `json:"code"`
	Name  string      `json:"name"`
	Price int         `json:"price"`
}

func (p *Product) Validate() error {
	var validationErrorDescriptions []*errors.ValidationErrorDescription

	//TODO: error messages should not be hardcoded
	if len(strings.TrimSpace(string(p.Code))) == 0 {
		validationErrorDescriptions = append(validationErrorDescriptions, errors.NewValidationErrorDescription("code", "Invalid product code"))
	}

	if p.Price <= 0 {
		validationErrorDescriptions = append(validationErrorDescriptions, errors.NewValidationErrorDescription("price", "Invalid product price"))
	}

	if len(validationErrorDescriptions) > 0 {
		return errors.NewValidationError(validationErrorDescriptions)
	}

	return nil
}
