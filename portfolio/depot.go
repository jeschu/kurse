package portfolio

import (
	"log"
	"os"
	"path"
	"time"

	"gopkg.in/yaml.v2"
)

type Depot struct {
	Stocks  []Stock `yaml:"stocks" json:"stocks"`
	Secrets Secrets `yaml:"secrets" json:"secrets"`
}

type Stock struct {
	Symbol    Symbol     `yaml:"symbol" json:"symbol"`
	Orders    []Order    `yaml:"orders" json:"orders"`
	Dividends []Dividend `yaml:"dividends" json:"dividends"`
}

type Symbol string

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

func LoadPortfolio() (map[Symbol]Stock, []Symbol, Secrets, error) {
	var (
		stocks   map[Symbol]Stock
		symbols  []Symbol
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
