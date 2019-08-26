package requests

import (
	"encoding/json"
	"github.com/alfcope/checkouttest/model"
	"io"
)

type AddItemRequest struct {
	Code model.ProductCode `json:"code"`
}

func NewAddItemRequest(body io.Reader) (*AddItemRequest, error) {
	var addItemRequest AddItemRequest

	decoder := json.NewDecoder(body)

	if err := decoder.Decode(&addItemRequest); err != nil {
		return nil, err
	}

	return &addItemRequest, nil
}
