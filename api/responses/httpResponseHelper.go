package responses

import (
	"encoding/json"
	"github.com/alfcope/checkouttest/errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

type NewBasketResponse struct {
	Id string `json:"id"`
}

type PriceBasketResponse struct {
	Total float64 `json:"total"`
}

// Sends a response error
func ResponseError(w http.ResponseWriter, log *logrus.Entry, status int, msg string) {
	if log != nil && msg != "" {
		log.Error(msg)
	}

	w.WriteHeader(status)
}

func Response(w http.ResponseWriter, log *logrus.Entry, status int, payload interface{}) {
	w.WriteHeader(status)

	if payload != nil {
		w.Header().Set("Content-Type", "application/json")
		jsonEncoded, err := json.Marshal(payload)

		if err != nil {
			ResponseError(w, log, http.StatusInternalServerError, err.Error())
		}

		_, err = w.Write(jsonEncoded)
		if err != nil {
			ResponseError(w, log, http.StatusInternalServerError, err.Error())
		}
	}
}

func GetStatusByError(err error) int {
	switch err.(type) {
	case *errors.BasketNotFound, *errors.ProductNotFound, *errors.PromotionNotFound:
		return http.StatusNotFound
	}

	return http.StatusInternalServerError
}
