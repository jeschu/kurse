package exchangerates

import (
	"encoding/json"
	"fmt"
	"kurse/lang"
	"net/http"
	"time"
)

const freeCurrencyApiUrl = "https://api.freecurrencyapi.com/v1/latest?apikey=%s&base_currency=EUR"

type Client struct {
	offline bool
	client  http.Client
	apiKey  string
}

func NewClient(apiKey string, timeout time.Duration, offline bool) *Client {
	return &Client{
		offline: offline,
		client:  http.Client{Timeout: timeout},
		apiKey:  apiKey,
	}
}

type Rates struct {
	Data map[string]float64 `json:"data"`
}

func (client *Client) FetchExchangeRates() (Rates, error) {
	if client.offline {
		return client.fetchExchangeRatesOffline()
	} else {
		return client.fetchExchangeRates()
	}
}

func (client *Client) fetchExchangeRates() (rates Rates, err error) {
	var (
		rq *http.Request
		rs *http.Response
	)
	url := fmt.Sprintf(freeCurrencyApiUrl, client.apiKey)
	if rq, err = http.NewRequest(http.MethodGet, url, nil); err != nil {
		return rates, err
	}
	if rs, err = client.client.Do(rq); err != nil {
		return rates, err
	}
	defer lang.Close(rs.Body, "unable to close response body")
	err = json.NewDecoder(rs.Body).Decode(&rates)
	return rates, err
}

func (client *Client) fetchExchangeRatesOffline() (Rates, error) {
	rates := Rates{}
	err := json.Unmarshal([]byte(offlineResponse), &rates)
	return rates, err
}

const offlineResponse = `{
  "data": {
    "AUD": 1.6488078757,
    "BGN": 1.9491438507,
    "BRL": 5.4918613368,
    "CAD": 1.4739704921,
    "CHF": 0.9784026472,
    "CNY": 7.8399326372,
    "CZK": 25.2903647161,
    "DKK": 7.4601321978,
    "EUR": 1,
    "GBP": 0.858214108,
    "HKD": 8.479604841,
    "HRK": 7.2652037798,
    "HUF": 389.3914865315,
    "IDR": 17186.1667682943,
    "ILS": 4.066590601,
    "INR": 90.2738232893,
    "ISK": 150.3900192176,
    "JPY": 164.3402781195,
    "KRW": 1462.4827271101,
    "MXN": 17.8292249803,
    "MYR": 5.1443104595,
    "NOK": 11.6487965414,
    "NZD": 1.8044646371,
    "PHP": 61.257863968,
    "PLN": 4.2811981823,
    "RON": 4.967197659,
    "RUB": 100.3406809574,
    "SEK": 11.5360316258,
    "SGD": 1.4620827238,
    "THB": 39.7160028943,
    "TRY": 34.5924670397,
    "USD": 1.0836582265,
    "ZAR": 20.271097858
  }
}`
