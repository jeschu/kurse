package portfolio

import (
	"log"
	"os"
	"path"
	"time"

	"gopkg.in/yaml.v2"
)

type Stocks = map[Symbol]Stock

type Depot struct {
	Stocks  []Stock `yaml:"stocks" json:"stocks"`
	Secrets Secrets `yaml:"secrets" json:"secrets"`
}

type Stock struct {
	Symbol    Symbol     `yaml:"symbol" json:"symbol"`
	Name      string     `yaml:"name" json:"name"`
	Wkn       string     `yaml:"wkn" json:"wkn"`
	Isin      string     `yaml:"isin" json:"isin"`
	Orders    []Order    `yaml:"orders" json:"orders"`
	Dividends []Dividend `yaml:"dividends" json:"dividends"`
}

type Symbol string
type Symbols []Symbol

type Order struct {
	Date      time.Time `yaml:"date" json:"date"`
	Count     float64   `yaml:"count" json:"count"`
	Price     float64   `yaml:"price" json:"price"`
	Provision float64   `yaml:"provision" json:"provision"`
	Fee       float64   `yaml:"fee" json:"fee"`
}

type Dividend struct {
	Date                  time.Time `yaml:"date" json:"date"`
	Count                 float64   `yaml:"count" json:"count"`
	Amount                float64   `yaml:"amount" json:"amount"`
	Quellensteuer         float64   `yaml:"quellensteuer" json:"quellensteuer"`
	Kapitalertragsteuer   float64   `yaml:"kapitalertragsteuer" json:"kapitalertragsteuer"`
	Solidaritaetszuschlag float64   `yaml:"solidaritaetszuschlag" json:"solidaritaetszuschlag"`
	Kirchensteuer         float64   `yaml:"kirchensteuer" json:"kirchensteuer"`
}

type Secrets struct {
	YahooKey           string `yaml:"yahooKey" json:"yahooKey"`
	YahooHost          string `yaml:"yahooHost" json:"yahooHost"`
	FreecurrencyApiKey string `yaml:"freecurrencyApiKey" json:"freecurrencyApiKey"`
}

func LoadPortfolio() (Stocks, Symbols, Secrets, error) {
	var (
		stocks   Stocks
		symbols  Symbols
		err      error
		yml      []byte
		filename string
	)
	if filename, err = portfolioConfigurationFile(); err != nil {
		return stocks, symbols, Secrets{}, err
	}
	log.Printf("loading portfolio from '%s'\n", filename)
	if yml, err = os.ReadFile(filename); err != nil {
		return stocks, symbols, Secrets{}, err
	}
	var depot = Depot{}
	if err = yaml.Unmarshal(yml, &depot); err != nil {
		return stocks, symbols, Secrets{}, err
	}
	stocks = make(map[Symbol]Stock)
	symbols = make([]Symbol, 0, len(depot.Stocks))
	for _, stock := range depot.Stocks {
		stocks[stock.Symbol] = stock
		symbols = append(symbols, stock.Symbol)
	}

	return stocks, symbols, depot.Secrets, err
}

func portfolioConfigurationFile() (filename string, err error) {
	var dir string
	if dir, err = os.UserConfigDir(); err == nil {
		filename = path.Join(dir, "kurse", "portfolio.yml")
		if _, err = os.Stat(filename); err == nil {
			return
		}
	}
	if dir, err = os.Getwd(); err == nil {
		filename = path.Join(dir, "portfolio.yml")
		_, err = os.Stat(filename)
	}
	return
}

func (s *Stock) Stueck() float64 {
	count := float64(0)
	for _, order := range s.Orders {
		count += order.Count
	}
	return count
}

func (s *Stock) Kaufkosten() (prices float64, fees float64, provisions float64, sum float64) {
	prices = 0
	fees = 0
	provisions = 0
	for _, order := range s.Orders {
		prices += order.Price
		fees += order.Fee
		provisions += order.Provision
	}
	sum = prices + fees + provisions
	return
}

func (s *Stock) Dividenden() (amount float64, quellensteuer float64, kapitalertragsteuer float64, solidaritaetszuschlag float64, kirchensteuer float64) {
	amount = 0
	quellensteuer = 0
	kapitalertragsteuer = 0
	solidaritaetszuschlag = 0
	kirchensteuer = 0
	for _, dividend := range s.Dividends {
		amount += dividend.Amount
		quellensteuer += dividend.Quellensteuer
		kapitalertragsteuer += dividend.Kapitalertragsteuer
		solidaritaetszuschlag += dividend.Solidaritaetszuschlag
		kirchensteuer += dividend.Kirchensteuer
	}
	return
}
