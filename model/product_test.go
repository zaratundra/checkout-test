package model

import (
	"github.com/alfcope/checkouttest/errors"
	"testing"
)

var productCases = []struct {
	product     Product
	fieldErrors map[string]string
}{
	{ // Code is mandatory
		product:     Product{Code: "", Name: "Product 1", Price: 1000},
		fieldErrors: map[string]string{"code": "Invalid product code"},
	}, { // Name is not mandatory
		product:     Product{Code: "P1", Name: "", Price: 1000},
		fieldErrors: nil,
	}, { // Price equals Zero
		product:     Product{Code: "P1", Name: "", Price: 0},
		fieldErrors: map[string]string{"price": "Invalid product price"},
	}, { // Price negative
		product:     Product{Code: "P1", Name: "", Price: -1},
		fieldErrors: map[string]string{"price": "Invalid product price"},
	}, { // Multiple errors
		product: Product{Code: "", Name: "", Price: -1},
		fieldErrors: map[string]string{"code": "Invalid product code",
			"price": "Invalid product price"},
	},
}

func TestProducts(t *testing.T) {
	for _, tc := range productCases {

		err := tc.product.Validate()

		if err == nil {
			if tc.fieldErrors != nil {
				t.Errorf("There should have been validation errors: %v ", tc.fieldErrors)
			}
			return
		}

		if tc.fieldErrors == nil {
			t.Errorf("There have been unexpected validation errors: %v", err.Error())
		} else {
			validationError, ok := err.(*errors.ValidationError)
			if ok {
				if len(tc.fieldErrors) != len(validationError.Errors) {
					t.Errorf("Expected number of error %v - Actual number of errors: %v", len(tc.fieldErrors), len(validationError.Errors))
				}

				for _, field := range validationError.Errors {
					message, keyExists := tc.fieldErrors[field.Field]

					if !keyExists {
						t.Errorf("Error for field %v not found", field.Field)
					}

					if message != field.Message {
						t.Errorf("Message expected for field %v : %v but was %v", field.Field, field.Message, message)
					}
				}
			} else {
				t.Errorf("There has been a different error: %v", err.Error())
			}
		}
	}
}
