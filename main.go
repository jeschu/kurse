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
	)

	if rq, err = http.NewRequest(http.MethodGet, "https://query1.finance.yahoo.com/v7/finance/quote?symbols=VNA.DE", nil); err != nil {
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

	for _, result := range response.QuoteResponse.Results {
		fmt.Printf("%s\n    %s (%s)\n", result.Symbol, result.LongName, result.ShortName)
		fmt.Printf("    Kurs: %10.2f %s\n", result.RegularMarketPrice, result.Currency)
	}
}
