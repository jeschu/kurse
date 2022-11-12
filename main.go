package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	var (
		err    error
		client = http.Client{Timeout: 5 * time.Second}
		rq     *http.Request
		rs     *http.Response
		stocks map[string]Stock
		param  string
	)

	if stocks, param, err = readDepot("depot.yml"); err != nil {
		log.Fatal(err)
	}
	if rq, err = http.NewRequest(http.MethodGet, "https://query1.finance.yahoo.com/v7/finance/quote?symbols="+param, nil); err != nil {
		log.Fatal(err)
	}
	if rs, err = client.Do(rq); err != nil {
		log.Fatal(err)
	}
	defer rs.Body.Close()
	var response = Response{}
	if err = json.NewDecoder(rs.Body).Decode(&response); err != nil {
		log.Fatal(err)
	}
	var results = make(map[string]Result)
	for _, result := range response.QuoteResponse.Results {
		results[result.Symbol] = result
	}

	valSum := float64(0)
	buySum := float64(0)
	for symbol, stock := range stocks {
		result, ok := results[symbol]
		fmt.Printf("%s\n", symbol)
		if ok {
			fmt.Printf("    %s (%s)\n", result.LongName, result.ShortName)
			fmt.Printf("    Kurs: %10.2f %s x %f\n", result.RegularMarketPrice, result.Currency, stock.Count)
			value := stock.Count * result.RegularMarketPrice
			fmt.Printf("    Wert: %10.2f %s\n", value, result.Currency)
			buy := stock.Price + stock.Provision + stock.Fee
			fmt.Printf("    Kauf: %10.2f %s (%.2fx%.2f=%.2f + %.2f + %.2f)\n", buy, result.Currency,
				stock.Count, stock.Price/stock.Count, stock.Price, stock.Provision, stock.Fee)
			guvV := value - buy
			guvP := (value / buy * 100) - 100
			fmt.Printf("          %+10.2f %s (%+.2f%%)\n", guvV, result.Currency, guvP)

			valSum += value
			buySum += buy
		}
	}
	fmt.Println("Summe:")
	fmt.Printf("    Wert: %10.2f %s\n", valSum, "EUR")
	fmt.Printf("    Kauf: %10.2f %s\n", buySum, "EUR")
	fmt.Printf("          %+10.2f %s (%+.2f%%)\n", valSum-buySum, "EUR", (valSum/buySum*100)-100)

}
