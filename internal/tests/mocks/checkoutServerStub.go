package mocks

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
)

type CheckoutServerStub struct {
	url        string
	httpServer *httptest.Server
	context    StubContext
}

type StubContext struct {
	responseStatusCode int
	payload            interface{}
}

func NewCheckServerStub(path string) *CheckoutServerStub {
	server := &CheckoutServerStub{
		url:        path,
		httpServer: nil,
		context:    StubContext{},
	}

	server.httpServer = server.initServer(path)

	return server
}

func (c *CheckoutServerStub) initServer(urlPath string) *httptest.Server {

	httpServer := httptest.NewServer(c.initializeRoutes(urlPath))

	return httpServer
}

func (c *CheckoutServerStub) initializeRoutes(urlPath string) *mux.Router {
	r := mux.NewRouter()
	fmt.Printf("%v/baskets/\n", urlPath)
	r.HandleFunc(fmt.Sprintf("%v/baskets/", urlPath), c.returnStub()).Methods("POST").Headers("Accept", "application/json")
	r.HandleFunc(fmt.Sprintf("%v/baskets/{id}/items/", urlPath), c.returnStub()).Methods("POST").Headers("Content-Type", "application/json")
	r.HandleFunc(fmt.Sprintf("%v/baskets/{id}", urlPath), c.returnStub()).Methods("GET").Queries("price", "").Headers("Accept", "application/json")
	r.HandleFunc(fmt.Sprintf("%v/baskets/{id}", urlPath), c.returnStub()).Methods("DELETE")

	return r
}

func (c *CheckoutServerStub) Close() {
	c.httpServer.Close()
}

func (c *CheckoutServerStub) GetUrl() string {
	return c.httpServer.URL
}

func (c *CheckoutServerStub) StubResponse(statusCode int, payload interface{}) {
	c.context.responseStatusCode = statusCode
	c.context.payload = payload
}

func (c *CheckoutServerStub) returnStub() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(c.context.responseStatusCode)

		if c.context.payload != nil {
			w.Header().Set("Content-Type", "application/json")

			jsonEncoded, err := json.Marshal(c.context.payload)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}

			_, err = w.Write(jsonEncoded)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	}
}
