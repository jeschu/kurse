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
		client    = http.Client{Timeout: 5 * time.Second}
		rq        *http.Request
		rs        *http.Response
		stocks    map[string]Stock
		param     string
		depotFile string
		out       = NewOut(language.German)
		rates     Rates
		results   = make(map[string]Result)
	)

	if depotFile, err = findDepot(); err != nil {
		log.Fatal(err)
	}
	log.Printf("loading depot from '%s'\n", depotFile)
	if stocks, param, err = readDepot(depotFile); err != nil {
		log.Fatal(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		if rq, err = http.NewRequest(http.MethodGet, "https://query1.finance.yahoo.com/v7/finance/quote?symbols="+param, nil); err != nil {
			log.Fatal(err)
		}
		if rs, err = client.Do(rq); err != nil {
			log.Fatal(err)
		}
		defer support.Close(rs.Body, "unable to close response body")
		var response = Response{}
		if err = json.NewDecoder(rs.Body).Decode(&response); err != nil {
			log.Fatal(err)
		}

		for _, result := range response.QuoteResponse.Results {
			results[result.Symbol] = result
		}
		wg.Done()
	}()
	go func() {
		if rates, err = fetchExchangeRates(); err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()
	wg.Wait()

	symbols := make([]string, 0, len(stocks))
	for symbol := range stocks {
		symbols = append(symbols, symbol)
	}
	sort.Strings(symbols)
	valSum := float64(0)
	buySum := float64(0)
	for _, symbol := range symbols {
		stock := stocks[symbol]
		result, ok := results[symbol]
		if ok {
			var (
				count     float64 = 0
				price     float64 = 0
				provision float64 = 0
				fee       float64 = 0
				buy       float64 = 0
			)

			var rate = 1.0
			currency, cok := rates.Data[result.Currency]
			if cok {
				rate = 1.0 / currency
			}

			for _, order := range stock.Orders {
				count += order.Count
				price += order.Price
				provision += order.Provision
				fee += order.Fee
				buy += order.Price + order.Provision + order.Fee
			}
			out.Printf("%s (%s)\n", result.LongName, result.ShortName)
			value := count * result.RegularMarketPrice
			out.Printf("    Wert: %10.2f %s = %10.2f %s x %f\n", value, result.Currency, result.RegularMarketPrice, result.Currency, count)
			eurValue := value * rate
			if rate != 1.0 {
				out.Printf("          %10.2f EUR = %10.2f EUR x %f\n", eurValue, result.RegularMarketPrice*rate, count)
			}
			out.Printf("    Kauf: %10.2f EUR (%.2fx%.2f=%.2f + %.2f + %.2f)\n", buy, count, price/count, price, provision, fee)
			guvV := eurValue - buy
			guvP := (eurValue / buy * 100) - 100
			out.Printf("     GuV: %+10.2f EUR (%+.2f%%)\n", guvV, guvP)
			valSum += value * rate
			buySum += buy * rate
		}
	}

	out.Println("Summe:")
	out.Printf("    Wert: %10.2f %s\n", valSum, "EUR")
	out.Printf("    Kauf: %10.2f %s\n", buySum, "EUR")
	out.Printf("     GuV: %+10.2f %s (%+.2f%%)\n", valSum-buySum, "EUR", (valSum/buySum*100)-100)

}
