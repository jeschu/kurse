package main

import (
	"os"
	"path"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type Depot struct {
	Stocks  []Stock `yaml:"stocks" json:"stocks"`
	Secrets Secrets `yaml:"secrets" json:"secrets"`
}

type Stock struct {
	Symbol    string     `yaml:"symbol" json:"symbol"`
	Orders    []Order    `yaml:"orders" json:"orders"`
	Dividends []Dividend `yaml:"dividends" json:"dividends"`
}

type Order struct {
	Date      time.Time `yaml:"date" json:"date"`
	Count     float64   `yaml:"count" json:"count"`
	Price     float64   `yaml:"price" json:"price"`
	Provision float64   `yaml:"provision" json:"provision"`
	Fee       float64   `yaml:"fee" json:"fee"`
}

type Dividend struct {
	Date   time.Time `yaml:"date" json:"date"`
	Count  float64   `yaml:"count" json:"count"`
	Amount float64   `yaml:"amount" json:"amount"`
}

type Secrets struct {
	Key  string `yaml:"key" json:"key"`
	Host string `yaml:"host" json:"host"`
}

func readDepot(filename string) (stocks map[string]Stock, param string, key string, host string, err error) {
	var (
		yml   []byte
		depot = Depot{}
	)
	if yml, err = os.ReadFile(filename); err != nil {
		return
	}
	if err = yaml.Unmarshal(yml, &depot); err != nil {
		return
	}
	stocks = make(map[string]Stock)
	symbols := make([]string, 0, len(depot.Stocks))
	for _, stock := range depot.Stocks {
		stocks[stock.Symbol] = stock
		symbols = append(symbols, stock.Symbol)
	}
	param = strings.Join(symbols, ",")
	key = depot.Secrets.Key
	host = depot.Secrets.Host
	return
}

func findDepot() (filename string, err error) {
	var dir string
	if dir, err = os.UserConfigDir(); err == nil {
		filename = path.Join(dir, "kurse", "depot.yml")
		if _, err = os.Stat(filename); err == nil {
			return
		}
	}
	if dir, err = os.Getwd(); err == nil {
		filename = path.Join(dir, "depot.yml")
		_, err = os.Stat(filename)
	}
	return
}
