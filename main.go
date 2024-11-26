package main

import (
	"fmt"
	"golang.org/x/text/language"
	"kurse/color"
	"kurse/exchangerates"
	"kurse/lang"
	"kurse/portfolio"
	"kurse/yahoo"
	"os"
	"sort"
	"sync"
)

func main() {
	useCache := isUseCache()
	out := NewOut(language.German)

	stocks, syms, secrets, err := portfolio.LoadPortfolio()
	lang.FatalOnError(err)

	results, rates := asyncFetch(secrets, syms, useCache)

	symbols := make([]string, 0, len(stocks))
	for symbol := range stocks {
		symbols = append(symbols, string(symbol))
	}
	sort.Strings(symbols)
	valSum := float64(0)
	buySum := float64(0)
	dividendSum := float64(0)
	dividendSteuerSum := float64(0)
	for _, symbol := range symbols {
		stock := stocks[portfolio.Symbol(symbol)]
		result, ok := results[symbol]
		if ok {
			var (
				orderCount                    float64 = 0
				orderPrice                    float64 = 0
				orderProvision                float64 = 0
				orderFee                      float64 = 0
				orderBuy                      float64 = 0
				dividendAmount                float64 = 0
				dividendQuellensteuer         float64 = 0
				dividendKapitalertragsteuer   float64 = 0
				dividendSolidaritaetszuschlag float64 = 0
				dividendKirchensteuer         float64 = 0
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
				dividendQuellensteuer += dividend.Quellensteuer
				dividendKapitalertragsteuer += dividend.Kapitalertragsteuer
				dividendSolidaritaetszuschlag += dividend.Solidaritaetszuschlag
				dividendKirchensteuer += dividend.Kirchensteuer
			}

			value := orderCount * result.RegularMarketPrice
			eurValue := value * rate
			guvV := eurValue + dividendAmount - orderBuy
			if guvV >= 0 {
				out.Print(color.GreenBackground, color.Black)
			} else {
				out.Print(color.RedBackground, color.Black)
			}
			var name string
			if result.LongName == "" {
				name = result.ShortName
			} else {
				name = fmt.Sprintf("%s (%s)", result.LongName, result.ShortName)
			}
			out.Printf("%s%s\n", name, color.Reset)
			out.Printf("            Wert: %10.2f %s = %10.2f %s x %f\n", value, result.Currency, result.RegularMarketPrice, result.Currency, orderCount)
			if rate != 1.0 {
				out.Printf("               %10.2f EUR = %10.2f EUR x %f\n", eurValue, result.RegularMarketPrice*rate, orderCount)
			}
			out.Printf("            Kauf: %10.2f EUR (%.2fx%.2f=%.2f + %.2f + %.2f)\n", orderBuy, orderCount, orderPrice/orderCount, orderPrice, orderProvision, orderFee)
			guvK := eurValue - orderBuy
			guvKP := (eurValue / orderBuy * 100) - 100
			out.Printf("             GuV: %s %s\n", color.ByAmount(guvK, "%+10.2f EUR"), color.ByAmount(guvKP, "(%+.2f%%)"))
			dividendSteuer := dividendQuellensteuer + dividendKapitalertragsteuer + dividendSolidaritaetszuschlag + dividendKirchensteuer
			out.Printf("       Dividende: %10.2f EUR (Brutto: %10.2f EUR | Steuer: %10.2f EUR)\n", dividendAmount, dividendAmount+dividendSteuer, dividendSteuer)
			guvP := ((eurValue + dividendAmount) / orderBuy * 100) - 100
			out.Printf("  GuV inkl. Div.: %s %s\n", color.ByAmount(guvV, "%+10.2f EUR"), color.ByAmount(guvP, "(%+.2f%%)"))
			valSum += value * rate
			buySum += orderBuy * rate
			dividendSum += dividendAmount
			dividendSteuerSum += dividendSteuer
			out.Println()
		}
	}

	out.Println("Summe:")
	out.Printf("            Wert: %10.2f %s\n", valSum, "EUR")
	out.Printf("            Kauf: %10.2f %s\n", buySum, "EUR")
	out.Printf("             GuV: %s %s\n", color.ByAmount(valSum-buySum, "%+10.2f EUR"), color.ByAmount(valSum/buySum*100-100, "(%+.2f%%)"))
	out.Printf("       Dividende: %10.2[1]f EUR (Brutto: %10.2[2]f EUR | Steuer: %10.2[3]f EUR)\n", dividendSum, dividendSum+dividendSteuerSum, dividendSteuerSum)
	out.Printf("  GuV inkl. Div.: %s %s\n", color.ByAmount(valSum+dividendSum-buySum, "%+10.2f EUR"), color.ByAmount(((valSum+dividendSum)/buySum*100)-100, "(%+.2f%%)"))

}

func asyncFetch(secrets portfolio.Secrets, syms []portfolio.Symbol, cached bool) (yahoo.Results, exchangerates.Rates) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	var results yahoo.Results
	go func(results *yahoo.Results, wg *sync.WaitGroup) {
		*results = yahoo.FetchStocks(syms, secrets, cached)
		wg.Done()
	}(&results, &wg)
	var rates exchangerates.Rates
	go func(rates *exchangerates.Rates, wg *sync.WaitGroup) {
		*rates = exchangerates.FetchExchangeRates(secrets, cached)
		wg.Done()
	}(&rates, &wg)
	wg.Wait()
	return results, rates
}

func isUseCache() bool {
	cache, ok := os.LookupEnv("CACHE")
	if ok && cache == "false" {
		return false
	} else {
		return true
	}
}
