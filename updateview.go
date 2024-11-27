package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var out = message.NewPrinter(language.German)

type View struct {
	view                  *tview.Grid
	name                  *tview.TextView
	symbol                *tview.TextView
	wkn                   *tview.TextView
	isin                  *tview.TextView
	kurs                  *tview.TextView
	stueck                *tview.TextView
	wert                  *tview.TextView
	kauf                  *tview.TextView
	fees                  *tview.TextView
	provisions            *tview.TextView
	guvKurs               *tview.TextView
	amount                *tview.TextView
	quellensteuer         *tview.TextView
	kapitalertragsteuer   *tview.TextView
	solidaritaetszuschlag *tview.TextView
	kirchensteuer         *tview.TextView
	guv                   *tview.TextView
	guvK                  *tview.TextView
}

func (v *View) addTextView(label string, row, column int) *tview.TextView {
	tv := tview.NewTextView().SetLabel(label + ": ")
	v.view.AddItem(tv, row, column, 1, 1, 0, 0, false)
	return tv
}

func NewView() *View {
	v := &View{}
	v.view = tview.NewGrid().
		SetRows(1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1).
		SetColumns(-1, -1)
	v.view.SetBorder(true)
	v.name = v.addTextView("                Name", 0, 0)
	v.symbol = v.addTextView("              Symbol", 1, 0)
	v.wkn = v.addTextView("                 WKN", 2, 0)
	v.isin = v.addTextView("                ISIN", 3, 0)
	v.kurs = v.addTextView("                Kurs", 5, 0)
	v.stueck = v.addTextView("               Stück", 6, 0)
	v.wert = v.addTextView("                Wert", 7, 0)
	v.kauf = v.addTextView("                Kauf", 8, 0)
	v.guvKurs = v.addTextView("          GuV (Kurs)", 9, 0)
	v.fees = v.addTextView("                Fees", 10, 0)
	v.provisions = v.addTextView("         Provisionen", 11, 0)
	v.amount = v.addTextView("          Dividenden", 12, 0)
	v.quellensteuer = v.addTextView("       Quellensteuer", 13, 0)
	v.kapitalertragsteuer = v.addTextView(" Kapitalertragsteuer", 14, 0)
	v.solidaritaetszuschlag = v.addTextView("Solidaritätszuschlag", 15, 0)
	v.kirchensteuer = v.addTextView("       Kirchensteuer", 16, 0)
	v.guv = v.addTextView("                 GuV", 17, 0)
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
	kauf, fees, provisions, kosten := stock.Kaufkosten()
	v.kauf.SetText(out.Sprintf("%10.2f %s", -kauf, cur))
	guvKurs := wert - kauf
	guvKursP := (wert / kosten * 100) - 100
	v.guvKurs.SetText(out.Sprintf("%10.2f %s (%+.2f%%)", guvKurs, cur, guvKursP)).SetTextColor(colorForAmount(guvKurs))
	v.fees.SetText(out.Sprintf("%10.2f %s", -fees, cur))
	v.provisions.SetText(out.Sprintf("%10.2f %s", -provisions, cur))
	dividenden, quellensteuer, kapitalertragsteuer, solidaritaetszuschlag, kirchensteuer := stock.Dividenden()
	v.amount.SetText(out.Sprintf("%10.2f %s", dividenden, cur))
	v.quellensteuer.SetText(out.Sprintf("%10.2f %s", -quellensteuer, cur))
	v.kapitalertragsteuer.SetText(out.Sprintf("%10.2f %s", -kapitalertragsteuer, cur))
	v.solidaritaetszuschlag.SetText(out.Sprintf("%10.2f %s", -solidaritaetszuschlag, cur))
	v.kirchensteuer.SetText(out.Sprintf("%10.2f %s", -kirchensteuer, cur))
	guv := wert - kosten + dividenden - quellensteuer - kapitalertragsteuer - solidaritaetszuschlag - kirchensteuer
	v.guv.SetText(out.Sprintf("%10.2f %s", guv, cur)).SetTextColor(colorForAmount(guv))
}

func colorForAmount(amount float64) tcell.Color {
	if amount < 0 {
		return tcell.ColorRed
	} else {
		return tcell.ColorGreen
	}
}
