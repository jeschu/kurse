package main

import (
	"github.com/rivo/tview"
	"kurse/portfolio"
	"kurse/yahoo"
	"sort"
)

func maxLen(i int, s string) int {
	l := len(s)
	if i > l {
		return i
	}
	return l
}

type Entries = map[string]Entry

type Entry struct {
	Stock  portfolio.Stock
	Result yahoo.Result
}

func prepare(stocks portfolio.Stocks, results yahoo.Results) (symList *tview.List, entries Entries, nameMaxLen int) {
	symList = tview.NewList().ShowSecondaryText(false).SetHighlightFullLine(true)
	symList.SetTitle(" Portfolio ").SetBorder(true)

	entries = make(Entries)
	names := make([]string, 0, len(stocks))
	nameMaxLen = 0
	for _, stock := range stocks {
		name := stock.Name
		entries[name] = Entry{Stock: stock, Result: results[string(stock.Symbol)]}
		nameMaxLen = maxLen(nameMaxLen, name)
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		entry := entries[name]
		symList.AddItem(name, entry.Result.Symbol, 0, nil)
	}
	return
}
