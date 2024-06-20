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
	Date          time.Time `yaml:"date" json:"date"`
	Count         float64   `yaml:"count" json:"count"`
	Amount        float64   `yaml:"amount" json:"amount"`
	Quellensteuer float64   `yaml:"quellensteuer" json:"quellensteuer"`
}

type Secrets struct {
	YahooKey           string `yaml:"yahooKey" json:"yahooKey"`
	YahooHost          string `yaml:"yahooHost" json:"yahooHost"`
	FreecurrencyApiKey string `yaml:"freecurrencyApiKey" json:"freecurrencyApiKey"`
}

func LoadPortfolio(debug bool) (map[Symbol]Stock, []Symbol, Secrets, error) {
	var (
		stocks  map[Symbol]Stock
		symbols []Symbol
		err     error
		yml     []byte
	)
	if debug {
		yml = []byte(debugDepot)
	} else {
		var (
			filename string
			err      error
		)
		if filename, err = portfolioConfigurationFile(); err != nil {
			return stocks, symbols, Secrets{}, err
		}
		log.Printf("loading portfolio from '%s'\n", filename)
		if yml, err = os.ReadFile(filename); err != nil {
			return stocks, symbols, Secrets{}, err
		}
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

const debugDepot = `
stocks:
  - symbol: "VWRD.L" # Vanguard FTSE All-World UCITS ETF
    orders:
      - date: 2023-01-03
        count: 0.265111
        price: 25.00
      - date: 2023-02-02
        count: 0.252118
        price: 25.00
      - date: 2023-03-01
        count: 0.257083
        price: 25.00
      - date: 2023-04-03
        count: 0.255678
        price: 25.00
      - date: 2023-05-03
        count: 0.256932
        price: 25.00
      - date: 2023-06-01
        count: 0.250495
        price: 25.00
      - date: 2023-07-03
        count: 0.243024
        price: 25.00
      - date: 2023-08-01
        count: 0.237752
        price: 25.00
      - date: 2023-09-02
        count: 0.23842
        price: 25.00
      - date: 2023-10-02
        count: 0.245285
        price: 25.00
    dividends:
      - date: 2023-03-29
        count: 0.77
        amount: 0.30
      - date: 2023-06-28
        count: 1.54
        amount: 1.02
      - date: 2023-09-27
        count: 2.26
        amount: 1.02
      - date: 2023-12-27
        count: 2.50
        amount: 0.92

  - symbol: "IQQH.DE" # iShares Global Clean Energy UCITS ETF USD (Dist) (ISHSII-GL.CL.ENERGY DLDIS)
    orders:
      - date: 2022-11-02
        count: 2.267615
        price: 25.00
      - date: 2022-12-02
        count: 2.162499
        price: 25.00
      - date: 2023-01-02
        count: 2.329916
        price: 25.00
      - date: 2023-02-02
        count: 2.275396
        price: 25.00
      - date: 2023-03-01
        count: 2.415319
        price: 25.00
      - date: 2023-04-03
        count: 2.424689
        price: 25.00
      - date: 2023-05-02
        count: 2.620435
        price: 25.00
      - date: 2023-06-01
        count: 2.557754
        price: 25.00
      - date: 2023-07-03
        count: 2.561161
        price: 25.00
      - date: 2023-08-02
        count: 2.626768
        price: 25.00
      - date: 2023-09-01
        count: 2.914364
        price: 25.00
      - date: 2023-10-02
        count: 3.195256
        price: 25.00
    dividends:
      - date: 2022-11-30
        count: 2.27
        amount: 0.06
      - date: 2023-05-30
        count: 16.50
        amount: 0.40
      - date: 2023-11-29
        count: 30.35
        amount: 1.24

  - symbol: "DBXD.DE" # Xtrackers DAX UCITS ETF
    orders:
      - date: 2023-03-01
        count: 0.171553
        price: 25.00
      - date: 2023-04-04
        count: 0.168183
        price: 25.00
      - date: 2023-05-02
        count: 0.167238
        price: 25.00
      - date: 2023-06-01
        count: 0.166725
        price: 25.00
      - date: 2023-07-03
        count: 0.164211
        price: 25.00
      - date: 2023-08-01
        count: 0.162056
        price: 25.00
      - date: 2023-09-01
        count: 0.166525
        price: 25.00
      - date: 2023-10-02
        count: 0.173242
        price: 25.00
    dividends:
      - date: 2024-02-07
        count: 1.339733
        amount: 1.21

  - symbol: "VNA.DE" # Vonovia
    orders:
      - date: 2022-12-18
        count: 4
        price: 87.60
        fee: 2.00
      - date: 2022-11-05
        count: 4
        price: 89.80
        provision: 5.90
        fee: 2.84
    dividends:
      - date: 2023-06-14
        count: 8.00
        amount: 6.80

secrets:
  yahooKey: "yahooKey"
  yahooHost: "yahoo-finance15.p.rapidapi.com"
  freecurrencyApiKey: "freecurrencyApiKey"
`
