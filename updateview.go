package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var out = message.NewPrinter(language.German)

type View struct {
	view       *tview.Grid
	name       *tview.TextView
	symbol     *tview.TextView
	wkn        *tview.TextView
	isin       *tview.TextView
	kurs       *tview.TextView
	stueck     *tview.TextView
	wert       *tview.TextView
	prices     *tview.TextView
	fees       *tview.TextView
	provisions *tview.TextView
	kosten     *tview.TextView
	guvKurs    *tview.TextView
}

func (v *View) addTextView(label string, row, column int) *tview.TextView {
	tv := tview.NewTextView().SetLabel(label + ": ")
	v.view.AddItem(tv, row, column, 1, 1, 0, 0, false)
	return tv
}

func NewView() *View {
	v := &View{}
	v.view = tview.NewGrid().
		SetRows(1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1).
		SetColumns(-1, -1)
	v.view.SetBorder(true)
	v.name = v.addTextView("  Name", 0, 0)
	v.symbol = v.addTextView("Symbol", 1, 0)
	v.wkn = v.addTextView("   WKN", 2, 0)
	v.isin = v.addTextView("  ISIN", 3, 0)
	v.kurs = v.addTextView("      Kurs", 5, 0)
	v.stueck = v.addTextView("     Stück", 6, 0)
	v.wert = v.addTextView("      Wert", 7, 0)
	v.prices = v.addTextView("      Kauf", 8, 0)
	v.guvKurs = v.addTextView("       GuV", 9, 0)
	v.fees = v.addTextView("        Fee", 6, 1)
	v.provisions = v.addTextView("  Provision", 7, 1)
	v.kosten = v.addTextView("     Kosten", 8, 1)
	return v
}

func updateView(v *View, entry Entry) {
	stock := entry.Stock
	result := entry.Result
	cur := result.Currency
	kurs := result.Bid
	kursType := "Bid"
	if kurs == 0 {
		kurs = result.RegularMarketPrice
		kursType = "Reg"
	}
	v.name.SetText(stock.Name)
	v.symbol.SetText(result.Symbol)
	v.isin.SetText(stock.Isin)
	v.wkn.SetText(stock.Wkn)
	v.kurs.SetText(out.Sprintf("%10.2f %s (%s)", kurs, cur, kursType))
	stueck := stock.Stueck()
	wert := stueck * kurs
	v.stueck.SetText(out.Sprintf("%10f", stueck))
	v.wert.SetText(out.Sprintf("%10.2f %s", wert, cur))
	prices, fees, provisions, kosten := stock.Kaufkosten()
	v.prices.SetText(out.Sprintf("%10.2f %s", -prices, cur))
	v.fees.SetText(out.Sprintf("%10.2f %s", fees, cur))
	v.provisions.SetText(out.Sprintf("%10.2f %s", provisions, cur))
	v.kosten.SetText(out.Sprintf("%10.2f %s", kosten, cur))
	guvKurs := wert - prices
	guvP := (wert / kosten * 100) - 100
	v.guvKurs.SetText(out.Sprintf("%10.2f %s (%+.2f%%)", guvKurs, cur, guvP)).SetTextColor(colorForAmount(guvKurs))
}

func colorForAmount(amount float64) tcell.Color {
	if amount < 0 {
		return tcell.ColorRed
	} else {
		return tcell.ColorGreen
	}
}
