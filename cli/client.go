package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alfcope/checkouttest/api/requests"
	"github.com/alfcope/checkouttest/api/responses"
	"github.com/alfcope/checkouttest/model"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type CheckoutClient struct {
	serverUrl  string
	apiVersion int
	httpClient *http.Client
}

func NewCheckoutClient(serverUrl string, version int) *CheckoutClient {
	return &CheckoutClient{
		serverUrl:  serverUrl,
		apiVersion: version,
		httpClient: &http.Client{
			Timeout: time.Second * 5,
		},
	}
}

func (c *CheckoutClient) AddBasket() (string, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v%d/baskets/", c.serverUrl, c.apiVersion), nil)
	if err != nil {
		return "", fmt.Errorf("there was an error creating http request: %v", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Printf("%s\n", resp.Status)
		return "", fmt.Errorf("%s", resp.Status)
	}

	if resp.Body != nil {
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("error fetching response body: %v", err)
		}

		nb := responses.NewBasketResponse{}
		err = json.Unmarshal(responseBody, &nb)
		if err != nil {
			return "", fmt.Errorf("error fetching response body: %v", err)
		}

		return nb.Id, nil
	}

	return "", errors.New("empty response")
}

func (c *CheckoutClient) AddItem(basketId, productCode string) error {
	if strings.TrimSpace(basketId) == "" || strings.TrimSpace(productCode) == "" {
		return errors.New("invalid request")
	}

	ir := requests.AddItemRequest{Code: model.ProductCode(productCode)}
	jsonRequest, err := json.Marshal(ir)
	if err != nil {
		return fmt.Errorf("there was an error creating http request: %v", err)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v%d/baskets/%s/items/", c.serverUrl, c.apiVersion, strings.TrimSpace(basketId)), bytes.NewBuffer(jsonRequest))
	if err != nil {
		return fmt.Errorf("there was an error creating http request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("%s", resp.Status)
	}

	return nil
}

func (c *CheckoutClient) GetPrice(basketId string) (float64, error) {
	if strings.TrimSpace(basketId) == "" {
		return float64(-1), errors.New("invalid request")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v%d/baskets/%s?price", c.serverUrl, c.apiVersion, strings.TrimSpace(basketId)), nil)
	if err != nil {
		return float64(-1), fmt.Errorf("there was an error creating http request: %v", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return float64(-1), err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return float64(-1), fmt.Errorf("%s", resp.Status)
	}

	if resp.Body != nil {
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return float64(-1), fmt.Errorf("error fetching response body: %v", err)
		}

		pb := responses.PriceBasketResponse{}
		err = json.Unmarshal(responseBody, &pb)
		if err != nil {
			return float64(-1), fmt.Errorf("error fetching response body: %v", err)
		}

		return pb.Total, nil
	}

	return float64(0), errors.New("empty response")
}

func (c *CheckoutClient) DeleteBasket(basketId string) error {
	if strings.TrimSpace(basketId) == "" {
		return errors.New("invalid request")
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v%d/baskets/%s", c.serverUrl, c.apiVersion, strings.TrimSpace(basketId)), nil)
	if err != nil {
		return fmt.Errorf("there was an error creating http request: %v", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s", resp.Status)
	}

	return nil
}
