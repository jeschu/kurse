package yahoo

import (
	"encoding/json"
	"kurse/lang"
	"kurse/portfolio"
	"log"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	offline bool
	client  http.Client
	host    string
	key     string
}

func NewClient(host string, key string, timeout time.Duration, offline bool) *Client {
	return &Client{
		offline: offline,
		client:  http.Client{Timeout: timeout},
		host:    host,
		key:     key,
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
	if client.offline {
		return client.fetchStocksOffline()
	} else {
		return client.fetchStocks(symbols)
	}
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

func (client *Client) fetchStocksOffline() (Results, error) {
	const responseBody = `
{
	"body": [
		{
			"ask": 114.38,
			"askSize": 0,
			"averageDailyVolume10Day": 12435,
			"averageDailyVolume3Month": 27660,
			"bid": 105.2,
			"bidSize": 0,
			"cryptoTradeable": false,
			"currency": "USD",
			"customPriceAlertConfidence": "LOW",
			"dividendYield": 1.9,
			"epsTrailingTwelveMonths": 6.599189,
			"esgPopulated": false,
			"exchange": "LSE",
			"exchangeDataDelayedBy": 15,
			"exchangeTimezoneName": "Europe/London",
			"exchangeTimezoneShortName": "GMT",
			"fiftyDayAverage": 108.91272,
			"fiftyDayAverageChange": 5.5122833,
			"fiftyDayAverageChangePercent": 0.05061193,
			"fiftyTwoWeekChangePercent": 10.716009,
			"fiftyTwoWeekHigh": 116.25,
			"fiftyTwoWeekHighChange": -1.824997,
			"fiftyTwoWeekHighChangePercent": -0.015698899,
			"fiftyTwoWeekLow": 98.07,
			"fiftyTwoWeekLowChange": 16.355003,
			"fiftyTwoWeekLowChangePercent": 0.16676867,
			"fiftyTwoWeekRange": "98.07 - 116.25",
			"firstTradeDateMilliseconds": 1337670000000,
			"fullExchangeName": "LSE",
			"gmtOffSetMilliseconds": 0,
			"language": "en-US",
			"longName": "Vanguard FTSE All-World UCITS ETF",
			"market": "gb_market",
			"marketState": "CLOSED",
			"messageBoardId": "finmb_215299893",
			"netAssets": 18377644000,
			"postMarketChange": null,
			"postMarketChangePercent": null,
			"postMarketPrice": null,
			"postMarketTime": null,
			"preMarketChange": null,
			"preMarketChangePercent": null,
			"preMarketPrice": null,
			"preMarketTime": null,
			"priceHint": 2,
			"quoteSourceName": "Delayed Quote",
			"quoteType": "ETF",
			"region": "US",
			"regularMarketChange": 0.7250061,
			"regularMarketChangePercent": 0.6376483,
			"regularMarketDayHigh": 114.37,
			"regularMarketDayLow": 113.72,
			"regularMarketDayRange": "113.72 - 114.37",
			"regularMarketOpen": 114,
			"regularMarketPreviousClose": 113.7,
			"regularMarketPrice": 114.425,
			"regularMarketTime": 1701447794,
			"regularMarketVolume": 9499,
			"shortName": "VANGUARD FUNDS PLC VANGUARD FTS",
			"sourceInterval": 15,
			"symbol": "VWRD.L",
			"tradeable": false,
			"trailingPE": 0.17339253,
			"trailingThreeMonthNavReturns": 0.99036,
			"trailingThreeMonthReturns": 0.99036,
			"triggerable": false,
			"twoHundredDayAverage": 108.97443,
			"twoHundredDayAverageChange": 5.450577,
			"twoHundredDayAverageChangePercent": 0.050017025,
			"typeDisp": "ETF",
			"ytdReturn": 16.15106
		},
		{
			"ask": 7.531,
			"askSize": 30008,
			"averageDailyVolume10Day": 536014,
			"averageDailyVolume3Month": 507524,
			"bid": 7.528,
			"bidSize": 35454,
			"cryptoTradeable": false,
			"currency": "EUR",
			"customPriceAlertConfidence": "LOW",
			"esgPopulated": false,
			"exchange": "GER",
			"exchangeDataDelayedBy": 15,
			"exchangeTimezoneName": "Europe/Berlin",
			"exchangeTimezoneShortName": "CET",
			"fiftyDayAverage": 7.4121,
			"fiftyDayAverageChange": 0.120900154,
			"fiftyDayAverageChangePercent": 0.016311187,
			"fiftyTwoWeekChangePercent": -34.54119,
			"fiftyTwoWeekHigh": 11.55,
			"fiftyTwoWeekHighChange": -4.017,
			"fiftyTwoWeekHighChangePercent": -0.3477922,
			"fiftyTwoWeekLow": 6.88,
			"fiftyTwoWeekLowChange": 0.6529999,
			"fiftyTwoWeekLowChangePercent": 0.094912775,
			"fiftyTwoWeekRange": "6.88 - 11.55",
			"firstTradeDateMilliseconds": 1199174400000,
			"fullExchangeName": "XETRA",
			"gmtOffSetMilliseconds": 3600000,
			"language": "en-US",
			"longName": "iShares Global Clean Energy UCITS ETF",
			"market": "de_market",
			"marketState": "CLOSED",
			"messageBoardId": "finmb_37568857",
			"netAssets": 3483972100,
			"postMarketChange": null,
			"postMarketChangePercent": null,
			"postMarketPrice": null,
			"postMarketTime": null,
			"preMarketChange": null,
			"preMarketChangePercent": null,
			"preMarketPrice": null,
			"preMarketTime": null,
			"priceHint": 2,
			"quoteSourceName": "Delayed Quote",
			"quoteType": "ETF",
			"region": "US",
			"regularMarketChange": 0.08799982,
			"regularMarketChangePercent": 1.181999,
			"regularMarketDayHigh": 7.534,
			"regularMarketDayLow": 7.379,
			"regularMarketDayRange": "7.379 - 7.534",
			"regularMarketOpen": 7.427,
			"regularMarketPreviousClose": 7.445,
			"regularMarketPrice": 7.533,
			"regularMarketTime": 1701448571,
			"regularMarketVolume": 853540,
			"shortName": "ISHSII-GL.CL.ENERGY DLDIS",
			"sourceInterval": 15,
			"symbol": "IQQH.DE",
			"tradeable": false,
			"trailingAnnualDividendRate": 0,
			"trailingAnnualDividendYield": 0,
			"triggerable": false,
			"twoHundredDayAverage": 9.036425,
			"twoHundredDayAverageChange": -1.5034246,
			"twoHundredDayAverageChangePercent": -0.16637383,
			"typeDisp": "ETF"
		},
		{
			"ask": 155.04,
			"askSize": 1148,
			"averageDailyVolume10Day": 44784,
			"averageDailyVolume3Month": 47580,
			"bid": 155.02,
			"bidSize": 1148,
			"cryptoTradeable": false,
			"currency": "EUR",
			"customPriceAlertConfidence": "LOW",
			"dividendYield": 0,
			"epsTrailingTwelveMonths": 10.288846,
			"esgPopulated": false,
			"exchange": "GER",
			"exchangeDataDelayedBy": 15,
			"exchangeTimezoneName": "Europe/Berlin",
			"exchangeTimezoneShortName": "CET",
			"fiftyDayAverage": 145.3572,
			"fiftyDayAverageChange": 9.662811,
			"fiftyDayAverageChangePercent": 0.06647632,
			"fiftyTwoWeekChangePercent": 12.709034,
			"fiftyTwoWeekHigh": 156.4,
			"fiftyTwoWeekHighChange": -1.3799896,
			"fiftyTwoWeekHighChangePercent": -0.008823464,
			"fiftyTwoWeekLow": 131.24,
			"fiftyTwoWeekLowChange": 23.779999,
			"fiftyTwoWeekLowChangePercent": 0.18119474,
			"fiftyTwoWeekRange": "131.24 - 156.4",
			"firstTradeDateMilliseconds": 1199260800000,
			"fullExchangeName": "XETRA",
			"gmtOffSetMilliseconds": 3600000,
			"language": "en-US",
			"longName": "Xtrackers DAX UCITS ETF",
			"market": "de_market",
			"marketState": "CLOSED",
			"messageBoardId": "finmb_35505419",
			"netAssets": 3785940740,
			"netExpenseRatio": 0.09,
			"postMarketChange": null,
			"postMarketChangePercent": null,
			"postMarketPrice": null,
			"postMarketTime": null,
			"preMarketChange": null,
			"preMarketChangePercent": null,
			"preMarketPrice": null,
			"preMarketTime": null,
			"priceHint": 2,
			"quoteSourceName": "Delayed Quote",
			"quoteType": "ETF",
			"region": "US",
			"regularMarketChange": 1.6200104,
			"regularMarketChangePercent": 1.0560694,
			"regularMarketDayHigh": 155.04,
			"regularMarketDayLow": 154.02,
			"regularMarketDayRange": "154.02 - 155.04",
			"regularMarketOpen": 154.3,
			"regularMarketPreviousClose": 153.4,
			"regularMarketPrice": 155.02,
			"regularMarketTime": 1701448583,
			"regularMarketVolume": 44397,
			"shortName": "XTR.DAX 1C",
			"sourceInterval": 15,
			"symbol": "DBXD.DE",
			"tradeable": false,
			"trailingAnnualDividendRate": 0,
			"trailingAnnualDividendYield": 0,
			"trailingPE": 15.066802,
			"trailingThreeMonthNavReturns": 1.61632,
			"trailingThreeMonthReturns": 1.61632,
			"triggerable": false,
			"twoHundredDayAverage": 148.6137,
			"twoHundredDayAverageChange": 6.406311,
			"twoHundredDayAverageChangePercent": 0.043107137,
			"typeDisp": "ETF",
			"ytdReturn": 15.75611
		},
		{
			"ask": 26.47,
			"askSize": 1083,
			"averageAnalystRating": "2.4 - Buy",
			"averageDailyVolume10Day": 2721716,
			"averageDailyVolume3Month": 3251945,
			"bid": 26.46,
			"bidSize": 1333,
			"bookValue": 34.192,
			"cryptoTradeable": false,
			"currency": "EUR",
			"customPriceAlertConfidence": "LOW",
			"dividendRate": 0.85,
			"dividendYield": 3.2,
			"earningsTimestamp": 1699016400,
			"earningsTimestampEnd": 1699016400,
			"earningsTimestampStart": 1699016400,
			"epsCurrentYear": 2.2,
			"epsForward": 2.05,
			"epsTrailingTwelveMonths": -8.03,
			"esgPopulated": false,
			"exchange": "GER",
			"exchangeDataDelayedBy": 15,
			"exchangeTimezoneName": "Europe/Berlin",
			"exchangeTimezoneShortName": "CET",
			"fiftyDayAverage": 23.0714,
			"fiftyDayAverageChange": 3.4785995,
			"fiftyDayAverageChangePercent": 0.1507754,
			"fiftyTwoWeekChangePercent": 12.02532,
			"fiftyTwoWeekHigh": 28.72,
			"fiftyTwoWeekHighChange": -2.17,
			"fiftyTwoWeekHighChangePercent": -0.075557105,
			"fiftyTwoWeekLow": 15.27,
			"fiftyTwoWeekLowChange": 11.279999,
			"fiftyTwoWeekLowChangePercent": 0.73870325,
			"fiftyTwoWeekRange": "15.27 - 28.72",
			"financialCurrency": "EUR",
			"firstTradeDateMilliseconds": 1373526000000,
			"forwardPE": 12.95122,
			"fullExchangeName": "XETRA",
			"gmtOffSetMilliseconds": 3600000,
			"language": "en-US",
			"longName": "Vonovia SE",
			"market": "de_market",
			"marketCap": 21628823552,
			"marketState": "CLOSED",
			"messageBoardId": "finmb_5537158",
			"postMarketChange": null,
			"postMarketChangePercent": null,
			"postMarketPrice": null,
			"postMarketTime": null,
			"preMarketChange": null,
			"preMarketChangePercent": null,
			"preMarketPrice": null,
			"preMarketTime": null,
			"priceEpsCurrentYear": 12.068181,
			"priceHint": 2,
			"priceToBook": 0.77649736,
			"quoteSourceName": "Delayed Quote",
			"quoteType": "EQUITY",
			"region": "US",
			"regularMarketChange": 1.0299988,
			"regularMarketChangePercent": 4.0360456,
			"regularMarketDayHigh": 26.58,
			"regularMarketDayLow": 25.47,
			"regularMarketDayRange": "25.47 - 26.58",
			"regularMarketOpen": 25.64,
			"regularMarketPreviousClose": 25.52,
			"regularMarketPrice": 26.55,
			"regularMarketTime": 1701448718,
			"regularMarketVolume": 3115648,
			"sharesOutstanding": 814644992,
			"shortName": "VONOVIA SE NA O.N.",
			"sourceInterval": 15,
			"symbol": "VNA.DE",
			"tradeable": false,
			"trailingAnnualDividendRate": 0.85,
			"trailingAnnualDividendYield": 0.03330721,
			"triggerable": false,
			"twoHundredDayAverage": 20.652676,
			"twoHundredDayAverageChange": 5.8973236,
			"twoHundredDayAverageChangePercent": 0.28554767,
			"typeDisp": "Equity"
		}
	],
	"meta": {
		"copywrite": "https://devAPI.ai",
		"processedTime": "2023-12-02T21:10:08.006637Z",
		"status": 200,
		"symbol": "Quotes Data",
		"version": "v1.0"
	}
}
	`
	reader := strings.NewReader(responseBody)
	var err error
	var resp = response{}
	err = json.NewDecoder(reader).Decode(&resp)
	if err != nil {
		return nil, err
	}
	var results = make(Results, len(resp.Results))
	for _, result := range resp.Results {
		results[result.Symbol] = result
	}
	return results, nil
}
