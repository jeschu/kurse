package exchangerates

import (
	"encoding/json"
	"fmt"
	"kurse/cached"
	"kurse/lang"
	"kurse/portfolio"
	"net/http"
	"time"
)

const freeCurrencyApiUrl = "https://api.freecurrencyapi.com/v1/latest?apikey=%s&base_currency=EUR"

type Client struct {
	client http.Client
	apiKey string
}

func NewClient(apiKey string, timeout time.Duration) *Client {
	return &Client{
		client: http.Client{Timeout: timeout},
		apiKey: apiKey,
	}
}

type Rates struct {
	Data map[string]float64 `json:"data"`
}

func FetchExchangeRates(secrets portfolio.Secrets, useCache bool) Rates {
	if useCache {
		r, ok := cached.Load("kurse", "exchangerates", 24*time.Hour, func(data []byte) *Rates {
			r := &Rates{}
			e := json.Unmarshal(data, r)
			lang.FatalOnError(e)
			return r
		})
		if ok {
			return *r
		}
	}
	client := NewClient(secrets.FreecurrencyApiKey, 10*time.Second)
	rates, err := client.FetchExchangeRates()
	lang.FatalOnError(err)
	cached.Save("kurse", "exchangerates", &rates, func(rates *Rates) (data []byte) {
		data, err = json.MarshalIndent(rates, "", "  ")
		lang.FatalOnError(err)
		return
	})
	return rates
}

func (client *Client) FetchExchangeRates() (Rates, error) { return client.fetchExchangeRates() }

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
