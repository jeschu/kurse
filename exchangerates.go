package main

import (
	"encoding/json"
	"kurse/support"
	"net/http"
)

const freecurrencyapiUrl = "https://api.freecurrencyapi.com/v1/latest?apikey=C84jJiZZ1WbWY2i1wKjZO8cEsqhz2SSN8KI1a5Le&base_currency=EUR"

func fetchExchangeRates() (rates Rates, err error) {
	var (
		rq *http.Request
		rs *http.Response
	)
	if rq, err = http.NewRequest(http.MethodGet, freecurrencyapiUrl, nil); err != nil {
		return
	}
	if rs, err = http.DefaultClient.Do(rq); err != nil {
		return
	}
	defer support.Close(rs.Body, "unable to close response body")
	err = json.NewDecoder(rs.Body).Decode(&rates)
	return

}

type Rates struct {
	Data map[string]float64 `json:"data"`
}
