package main

import (
	"encoding/json"
	"kurse/support"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"

	"golang.org/x/text/language"
)

func main() {
	var (
		err       error
		stocks    map[string]Stock
		param     string
		depotFile string
		out       = NewOut(language.German)
		rates     Rates
		results   = make(map[string]Result)
		key       string
		host      string
	)

	if depotFile, err = findDepot(); err != nil {
		log.Fatal(err)
	}
	log.Printf("loading depot from '%s'\n", depotFile)
	if stocks, param, key, host, err = readDepot(depotFile); err != nil {
		log.Fatal(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(2)

	go fetchStocks(results, &wg, param, key, host) // go fetchStocksOffline(results, &wg, param, key, host)
	go fetchRates(rates, &wg)
	wg.Wait()

	symbols := make([]string, 0, len(stocks))
	for symbol := range stocks {
		symbols = append(symbols, symbol)
	}
	sort.Strings(symbols)
	valSum := float64(0)
	buySum := float64(0)
	dividendSum := float64(0)
	for _, symbol := range symbols {
		stock := stocks[symbol]
		result, ok := results[symbol]
		if ok {
			var (
				orderCount     float64 = 0
				orderPrice     float64 = 0
				orderProvision float64 = 0
				orderFee       float64 = 0
				orderBuy       float64 = 0
				dividendAmount float64 = 0
			)

			var rate = 1.0
			currency, cok := rates.Data[result.Currency]
			if cok {
				rate = 1.0 / currency
			}

			for _, order := range stock.Orders {
				orderCount += order.Count
				orderPrice += order.Price
				orderProvision += order.Provision
				orderFee += order.Fee
				orderBuy += order.Price + order.Provision + order.Fee
			}
			for _, dividend := range stock.Dividends {
				dividendAmount += dividend.Amount
			}
			out.Printf("%s (%s)\n", result.LongName, result.ShortName)
			value := orderCount * result.RegularMarketPrice
			out.Printf("         Wert: %10.2f %s = %10.2f %s x %f\n", value, result.Currency, result.RegularMarketPrice, result.Currency, orderCount)
			eurValue := value * rate
			if rate != 1.0 {
				out.Printf("               %10.2f EUR = %10.2f EUR x %f\n", eurValue, result.RegularMarketPrice*rate, orderCount)
			}
			out.Printf("         Kauf: %10.2f EUR (%.2fx%.2f=%.2f + %.2f + %.2f)\n", orderBuy, orderCount, orderPrice/orderCount, orderPrice, orderProvision, orderFee)
			out.Printf("    Dividende: %10.2f EUR\n", dividendAmount)
			guvV := eurValue + dividendAmount - orderBuy
			guvP := ((eurValue + dividendAmount) / orderBuy * 100) - 100
			out.Printf("          GuV: %+10.2f EUR (%+.2f%%)\n", guvV, guvP)
			valSum += value * rate
			buySum += orderBuy * rate
			dividendSum += dividendAmount
			out.Println()
		}
	}

	out.Println("Summe:")
	out.Printf("            Wert: %10.2f %s\n", valSum, "EUR")
	out.Printf("            Kauf: %10.2f %s\n", buySum, "EUR")
	out.Printf("             GuV: %+10.2f %s (%+.2f%%)\n", valSum-buySum, "EUR", (valSum/buySum*100)-100)
	out.Printf("       Dividende: %10.2f %s\n", dividendSum, "EUR")
	out.Printf("  GuV inkl. Div.: %+10.2f %s (%+.2f%%)\n", valSum+dividendSum-buySum, "EUR", ((valSum+dividendSum)/buySum*100)-100)

}

var fetchStocks = func(results map[string]Result, wg *sync.WaitGroup, param string, key string, host string) {
	var (
		err    error
		client = http.Client{Timeout: 5 * time.Second}
		rq     *http.Request
		rs     *http.Response
	)
	if rq, err = http.NewRequest(http.MethodGet, "https://yahoo-finance15.p.rapidapi.com/api/v1/markets/stock/quotes?ticker="+param, nil); err != nil {
		log.Fatal(err)
	}
	rq.Header.Add("X-RapidAPI-Key", key)
	rq.Header.Add("X-RapidAPI-Host", host)
	if rs, err = client.Do(rq); err != nil {
		log.Fatal(err)
	}
	defer support.Close(rs.Body, "unable to close response body")
	var response = Response{}
	if err = json.NewDecoder(rs.Body).Decode(&response); err != nil {
		log.Fatal(err)
	}

	rateLimitLimit := rs.Header.Get("x-ratelimit-limit")
	rateLimitRemaining := rs.Header.Get("x-ratelimit-remaining")

	log.Printf("fetch stocks => rate remaining/limit: %s/%s\n", rateLimitRemaining, rateLimitLimit)

	for _, result := range response.Results {
		results[result.Symbol] = result
	}
	wg.Done()
}

var fetchRates = func(rates Rates, wg *sync.WaitGroup) {
	var err error
	if rates, err = fetchExchangeRates(); err != nil {
		log.Fatal(err)
	}
	wg.Done()
}
