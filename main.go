package main

import (
	"kurse/exchangerates"
	"kurse/lang"
	"kurse/portfolio"
	"kurse/yahoo"
	"os"
	"sort"
	"sync"
	"time"

	"golang.org/x/text/language"
)

func main() {
	debug := isDebug()
	out := NewOut(language.German)

	stocks, syms, secrets, err := portfolio.LoadPortfolio(debug)
	lang.FatalOnError(err)

	results, rates := asyncFetch(secrets, syms, debug)

	symbols := make([]string, 0, len(stocks))
	for symbol := range stocks {
		symbols = append(symbols, string(symbol))
	}
	sort.Strings(symbols)
	valSum := float64(0)
	buySum := float64(0)
	dividendSum := float64(0)
	for _, symbol := range symbols {
		stock := stocks[portfolio.Symbol(symbol)]
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

func isDebug() bool {
	debug, ok := os.LookupEnv("DEBUG")
	return ok && debug == "true"
}

func asyncFetch(secrets portfolio.Secrets, syms []portfolio.Symbol, offline bool) (yahoo.Results, exchangerates.Rates) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	var results yahoo.Results
	go func(results *yahoo.Results, wg *sync.WaitGroup) {
		client := yahoo.NewClient(secrets.YahooHost, secrets.YahooKey, 10*time.Second, offline)
		rs, e := client.FetchStocks(syms)
		lang.FatalOnError(e)
		*results = rs
		wg.Done()
	}(&results, &wg)
	var rates exchangerates.Rates
	go func(rates *exchangerates.Rates, wg *sync.WaitGroup) {
		client := exchangerates.NewClient(secrets.FreecurrencyApiKey, 10*time.Second, offline)
		rs, e := client.FetchExchangeRates()
		lang.FatalOnError(e)
		*rates = rs
		wg.Done()
	}(&rates, &wg)
	wg.Wait()
	return results, rates
}
