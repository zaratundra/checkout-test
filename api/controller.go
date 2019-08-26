package api

import (
	"github.com/alfcope/checkouttest/api/requests"
	"github.com/alfcope/checkouttest/api/responses"
	"github.com/alfcope/checkouttest/pkg/logging"
	"github.com/gorilla/mux"
	"net/http"
)

type CheckoutController struct {
	checkoutService CheckoutService
}

func NewCheckoutController(router *mux.Router, service CheckoutService) *CheckoutController {
	controller := &CheckoutController{
		checkoutService: service,
	}

	controller.initializeRoutes(router)

	return controller
}

func (c *CheckoutController) initializeRoutes(router *mux.Router) {

	checkoutRouter := router.PathPrefix("/baskets").Subrouter()
	checkoutRouter.Use(logging.AccessLoggingMiddleware)

	// swagger:route POST / payments postPayment
	checkoutRouter.HandleFunc("/", c.CreateBasket()).Methods("POST").Headers("Accept", "application/json")
	// swagger:route GET /{id} payments getPayment
	checkoutRouter.HandleFunc("/{id}/items/", c.AddItem()).Methods("POST").Headers("Content-Type", "application/json")
	// swagger:route GET / payments getPaymentsPage
	checkoutRouter.HandleFunc("/{id}", c.GetPrice()).Methods("GET").Queries("price", "").Headers("Accept", "application/json")
	// swagger:route DELETE /{id} payments deletePayment
	checkoutRouter.HandleFunc("/{id}", c.DeleteBasket()).Methods("DELETE")
}

// PostPayment handles requests to add a payment into the system. The new payment
// will be linked to the organisation making the request.
// Http method: POST
// Path parameter: payment id
// Return: the new payment resource if successful or a http error code otherwise.
func (c *CheckoutController) CreateBasket() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logging.GetLoggerWithFields(r)

		basketId, err := c.checkoutService.CreateBasket()
		if err != nil {
			responses.ResponseError(w, logger, responses.GetStatusByError(err), err.Error())
			return
		}

		responses.Response(w, logger, http.StatusCreated, responses.NewBasketResponse{Id: basketId})
	}
}

// PostPayment handles requests to add a payment into the system. The new payment
// will be linked to the organisation making the request.
// Http method: POST
// Path parameter: payment id
// Return: the new payment resource if successful or a http error code otherwise.
func (c *CheckoutController) AddItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logging.GetLoggerWithFields(r)

		pathParameters := mux.Vars(r)
		basketId := pathParameters["id"]

		request, err := requests.NewAddItemRequest(r.Body)
		if err != nil {
			responses.ResponseError(w, logger, responses.GetStatusByError(err), err.Error())
			return
		}

		if request.Code == "" {
			responses.ResponseError(w, logger, http.StatusUnprocessableEntity, "Empty product code")
		}

		err = c.checkoutService.AddProduct(basketId, request.Code)
		if err != nil {
			responses.ResponseError(w, logger, responses.GetStatusByError(err), err.Error())
			return
		}

		responses.Response(w, logger, http.StatusCreated, nil)
	}
}

// PostPayment handles requests to add a payment into the system. The new payment
// will be linked to the organisation making the request.
// Http method: POST
// Path parameter: payment id
// Return: the new payment resource if successfull or a http error code otherwise.
func (c *CheckoutController) GetPrice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logging.GetLoggerWithFields(r)

		pathParameters := mux.Vars(r)
		basketId := pathParameters["id"]

		total, err := c.checkoutService.GetBasketPrice(basketId)
		if err != nil {
			responses.ResponseError(w, logger, responses.GetStatusByError(err), err.Error())
			return
		}
		responses.Response(w, logger, http.StatusOK, responses.PriceBasketResponse{Total: total})
	}
}

// PostPayment handles requests to add a payment into the system. The new payment
// will be linked to the organisation making the request.
// Http method: POST
// Path parameter: payment id
// Return: the new payment resource if successful or a http error code otherwise.
func (c *CheckoutController) DeleteBasket() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logging.GetLoggerWithFields(r)

		pathParameters := mux.Vars(r)
		basketId := pathParameters["id"]

		c.checkoutService.DeleteBasket(basketId)

		responses.Response(w, logger, http.StatusNoContent, nil)
	}
}
