package main

import (
	"os"
	"path"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type Depot struct {
	Stocks []Stock `yaml:"stocks" json:"stocks"`
}

type Stock struct {
	Symbol string  `yaml:"symbol" json:"symbol"`
	Orders []Order `yaml:"orders" json:"orders"`
}

type Order struct {
	Date      time.Time `yaml:"date" json:"date"`
	Count     float64   `yaml:"count" json:"count"`
	Price     float64   `yaml:"price" json:"price"`
	Provision float64   `yaml:"provision" json:"provision"`
	Fee       float64   `yaml:"fee" json:"fee"`
}

func readDepot(filename string) (stocks map[string]Stock, param string, err error) {
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
