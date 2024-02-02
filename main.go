package main

import (
	"log"
	"encoding/json"
	"fmt"

	"github.com/samjtro/go-tda/data"
	"github.com/gdamore/tcell/v2"
        "github.com/rivo/tview"
)

type optionFlags struct {
	ticker string
	periodType string
	period string
	freqType string
	freq string
}

var (
	app = tview.NewApplication()
	tickerText = tview.NewTextView()
	form = tview.NewForm()
	pages = tview.NewPages()
	text = tview.NewTextView().
		SetTextColor(tcell.ColorGreen).
		SetText("(n) to search for a new ticker\n(q) to quit")
	flex = tview.NewFlex()
	tickers = make([]optionFlags, 0)
	tickersList = tview.NewList().ShowSecondaryText(false)
)

func main() {
	tickersList.SetSelectedFunc(func(index int, ticker string, company string, shortcut rune) {
		pullTickerData((&tickers[index]).ticker) //TODO: I think the problem is this func
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

	for index, ticker := range tickers {
		tickersList.AddItem(ticker.ticker, " ", rune(49+index), nil)
	}
}

func searchTicker() *tview.Form {
	list := optionFlags{}

        form.AddInputField("Ticker", "", 20, nil, func(ticker string) {
		list.ticker = ticker
        })

	form.AddButton("Save", func() {
		tickers = append(tickers, list)
		addTickerList()
		pages.SwitchToPage("Menu")
	})

	return form
}

// TODO: This function does not work as intended, potentially it is json.Marshal that is giving me the issue?
func pullTickerData(ticker string) {
	tickerText.Clear()
	quote, err := data.RealTime(ticker)

	if err != nil {
        	log.Fatalf(err.Error())
	}

	text, err := json.Marshal(quote)

	if err != nil {
		log.Fatalf(err.Error())
	}

	tickerText.SetText(string(text))
}
