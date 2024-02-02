package main

import (
	"log"

	"github.com/samjtro/go-tda/data"
	"github.com/gdamore/tcell/v2"
        "github.com/rivo/tview"
)

var (
	app = tview.NewApplication()
	tickerText = tview.NewTextView()
	form = tview.NewForm()
	pages = tview.NewPages()
	text = tview.NewTextView().
		SetTextColor(tcell.ColorGreen).
		SetText("(n) to search for a new ticker\n(q) to quit")
	flex = tview.NewFlex()
	tickers = make(map[int]string)
	tickersList = tview.NewList().ShowSecondaryText(false)
)

func main() {
	//TODO: Something in this block is causing an error
	tickersList.SetSelectedFunc(func(index int, ticker string, secondary string, shortcut rune) {
		pullTickerData((tickers[index]))
	})

	flex.SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(tickersList, 0, 1, true).
			AddItem(tickerText, 0, 4, false), 0, 6, false).
		AddItem(text, 0, 1, false)

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 113 {
			app.Stop()
		} else if event.Rune() == 110 {
			form.Clear(true)
			searchTicker()
			pages.SwitchToPage("Search")
		}
		return event
	})

	pages.AddPage("Menu", flex, true, true)
	pages.AddPage("Search", form, true, false)

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		log.Fatalf(err.Error())
	}
}

func addTickerList() {
	tickersList.Clear()

	for i := 0; i <= len(tickers); i++ {
		tickersList.AddItem(tickers[i], " ", rune(49+i), nil)
	}
}

func searchTicker() *tview.Form {
	var ticker string

        form.AddInputField("Ticker", "", 20, nil, func(t string) {
		ticker = t
	})

	form.AddButton("Save", func() {
		tickers[len(tickers)] = ticker
		addTickerList()
		pages.SwitchToPage("Menu")
	})

	return form
}

//TODO: Potentially broken, need to figure out a way to signal this is triggering
func pullTickerData(ticker string) {
	tickerText.Clear()
	quote, err := data.RealTime(ticker)

	if err != nil {
        	log.Fatalf(err.Error())
	}

	text := "DateTime: " + quote.Datetime + "\nTicker: " + quote.Ticker
	tickerText.SetText(text)
}
