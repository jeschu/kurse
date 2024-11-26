package yahoo

import (
	"encoding/json"
	"kurse/cached"
	"kurse/lang"
	"kurse/portfolio"
	"log"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	client http.Client
	host   string
	key    string
}

func FetchStocks(symbols []portfolio.Symbol, secrets portfolio.Secrets, useCache bool) Results {
	var (
		results Results
		err     error
	)
	if useCache {
		r, ok := cached.Load("kurse", "yahoo", 24*time.Hour, func(data []byte) *Results {
			r := &Results{}
			e := json.Unmarshal(data, r)
			lang.FatalOnError(e)
			return r
		})
		if ok {
			return *r
		}
	}
	client := NewClient(secrets.YahooHost, secrets.YahooKey, 10*time.Second)
	results, err = client.FetchStocks(symbols)
	lang.FatalOnError(err)
	cached.Save("kurse", "yahoo", &results, func(results *Results) (data []byte) {
		data, err = json.MarshalIndent(results, "", "  ")
		lang.FatalOnError(err)
		return
	})
	return results
}

func NewClient(host string, key string, timeout time.Duration) *Client {
	return &Client{
		client: http.Client{Timeout: timeout},
		host:   host,
		key:    key,
	}
}

type Results map[string]Result

type Result struct {
	Language                          string  `json:"language"`
	Region                            string  `json:"region"`
	QuoteType                         string  `json:"quoteType"`
	TypeDisp                          string  `json:"typeDisp"`
	QuoteSourceName                   string  `json:"quoteSourceName"`
	Triggerable                       bool    `json:"triggerable"`
	CustomPriceAlertConfidence        string  `json:"customPriceAlertConfidence"`
	Market                            string  `json:"market"`
	EsgPopulated                      bool    `json:"esgPopulated"`
	ShortName                         string  `json:"shortName"`
	LongName                          string  `json:"longName"`
	MessageBoardId                    string  `json:"messageBoardId"`
	ExchangeTimezoneName              string  `json:"exchangeTimezoneName"`
	ExchangeTimezoneShortName         string  `json:"exchangeTimezoneShortName"`
	GmtOffSetMilliseconds             int     `json:"gmtOffSetMilliseconds"`
	Currency                          string  `json:"currency"`
	MarketState                       string  `json:"marketState"`
	RegularMarketChangePercent        float64 `json:"regularMarketChangePercent"`
	RegularMarketPrice                float64 `json:"regularMarketPrice"`
	Exchange                          string  `json:"exchange"`
	FirstTradeDateMilliseconds        int64   `json:"firstTradeDateMilliseconds"`
	PriceHint                         int     `json:"priceHint"`
	RegularMarketChange               float64 `json:"regularMarketChange"`
	RegularMarketTime                 int     `json:"regularMarketTime"`
	RegularMarketDayHigh              float64 `json:"regularMarketDayHigh"`
	RegularMarketDayRange             string  `json:"regularMarketDayRange"`
	RegularMarketDayLow               float64 `json:"regularMarketDayLow"`
	RegularMarketVolume               int     `json:"regularMarketVolume"`
	RegularMarketPreviousClose        float64 `json:"regularMarketPreviousClose"`
	Bid                               float64 `json:"bid"`
	Ask                               float64 `json:"ask"`
	BidSize                           int     `json:"bidSize"`
	AskSize                           int     `json:"askSize"`
	FullExchangeName                  string  `json:"fullExchangeName"`
	FinancialCurrency                 string  `json:"financialCurrency"`
	RegularMarketOpen                 float64 `json:"regularMarketOpen"`
	AverageDailyVolume3Month          int     `json:"averageDailyVolume3Month"`
	AverageDailyVolume10Day           int     `json:"averageDailyVolume10Day"`
	FiftyTwoWeekLowChange             float64 `json:"fiftyTwoWeekLowChange"`
	FiftyTwoWeekLowChangePercent      float64 `json:"fiftyTwoWeekLowChangePercent"`
	FiftyTwoWeekRange                 string  `json:"fiftyTwoWeekRange"`
	FiftyTwoWeekHighChange            float64 `json:"fiftyTwoWeekHighChange"`
	FiftyTwoWeekHighChangePercent     float64 `json:"fiftyTwoWeekHighChangePercent"`
	FiftyTwoWeekLow                   float64 `json:"fiftyTwoWeekLow"`
	FiftyTwoWeekHigh                  float64 `json:"fiftyTwoWeekHigh"`
	EarningsTimestamp                 int     `json:"earningsTimestamp"`
	EarningsTimestampStart            int     `json:"earningsTimestampStart"`
	EarningsTimestampEnd              int     `json:"earningsTimestampEnd"`
	TrailingAnnualDividendRate        float64 `json:"trailingAnnualDividendRate"`
	TrailingPE                        float64 `json:"trailingPE"`
	TrailingAnnualDividendYield       float64 `json:"trailingAnnualDividendYield"`
	EpsTrailingTwelveMonths           float64 `json:"epsTrailingTwelveMonths"`
	EpsForward                        float64 `json:"epsForward"`
	EpsCurrentYear                    float64 `json:"epsCurrentYear"`
	PriceEpsCurrentYear               float64 `json:"priceEpsCurrentYear"`
	SharesOutstanding                 int     `json:"sharesOutstanding"`
	BookValue                         float64 `json:"bookValue"`
	FiftyDayAverage                   float64 `json:"fiftyDayAverage"`
	FiftyDayAverageChange             float64 `json:"fiftyDayAverageChange"`
	FiftyDayAverageChangePercent      float64 `json:"fiftyDayAverageChangePercent"`
	TwoHundredDayAverage              float64 `json:"twoHundredDayAverage"`
	TwoHundredDayAverageChange        float64 `json:"twoHundredDayAverageChange"`
	TwoHundredDayAverageChangePercent float64 `json:"twoHundredDayAverageChangePercent"`
	MarketCap                         int64   `json:"marketCap"`
	ForwardPE                         float64 `json:"forwardPE"`
	PriceToBook                       float64 `json:"priceToBook"`
	SourceInterval                    int     `json:"sourceInterval"`
	ExchangeDataDelayedBy             int     `json:"exchangeDataDelayedBy"`
	AverageAnalystRating              string  `json:"averageAnalystRating"`
	Tradeable                         bool    `json:"tradeable"`
	CryptoTradeable                   bool    `json:"cryptoTradeable"`
	Symbol                            string  `json:"symbol"`
}

type response struct {
	Results []Result `json:"body"`
	Meta    meta     `json:"meta"`
}

type meta struct {
	Copywrite     string    `json:"copywrite"`
	ProcessedTime time.Time `json:"processedTime"`
	Status        int       `json:"status"`
	Symbol        string    `json:"symbol"`
	Version       string    `json:"version"`
}

func (client *Client) FetchStocks(symbols []portfolio.Symbol) (Results, error) {
	return client.fetchStocks(symbols)
}

func (client *Client) fetchStocks(symbols []portfolio.Symbol) (Results, error) {
	var (
		rq  *http.Request
		rs  *http.Response
		err error
	)
	ticker := sliceOfSymbolsToQueryParam(symbols)
	rq, err = http.NewRequest(http.MethodGet, "https://yahoo-finance15.p.rapidapi.com/api/v1/markets/stock/quotes?ticker="+ticker, nil)
	if err != nil {
		return nil, err
	}
	rq.Header.Add("X-RapidAPI-Key", client.key)
	rq.Header.Add("X-RapidAPI-Host", client.host)

	rs, err = client.client.Do(rq)
	if err != nil {
		return nil, err
	}
	defer lang.Close(rs.Body, "unable to close response body")
	var resp = response{}
	err = json.NewDecoder(rs.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}

	rateLimitLimit := rs.Header.Get("x-ratelimit-limit")
	rateLimitRemaining := rs.Header.Get("x-ratelimit-remaining")
	log.Printf("fetch stocks => rate remaining/limit: %s/%s\n", rateLimitRemaining, rateLimitLimit)

	results := make(Results, len(symbols))
	for _, result := range resp.Results {
		results[result.Symbol] = result
	}

	return results, nil
}

func sliceOfSymbolsToQueryParam(symbols []portfolio.Symbol) string {
	params := make([]string, len(symbols))
	for idx, symbol := range symbols {
		params[idx] = string(symbol)
	}
	return strings.Join(params, ",")
}
