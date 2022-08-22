package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/samjtro/go-tda/data"
)

var (
	tickerFlag     string
	periodTypeFlag string
	periodFlag     string
	freqTypeFlag   string
	freqFlag       string
	m              sync.Mutex
	m1             sync.Mutex
	m2             sync.Mutex
	m3             sync.Mutex
	m4             sync.Mutex
	m5             sync.Mutex
	m6             sync.Mutex
	m7             sync.Mutex
)

type DATA struct {
	Close            float64
	Hi               float64
	Lo               float64
	Volume           float64
	PivotPoint       float64
	ResistancePoints []float64
	SupportPoints    []float64
	SMA              float64
	RMA              float64
	EMA              float64
	RSI              float64
	VWAP             float64
	MACD             float64
	Chaikin          float64
	BollingerBands   []float64
	IMI              float64
	MFI              float64
	PCR              float64
	OI               float64
}

type DataSlice []DATA

func main() {
	start := time.Now()

	flag.StringVar(&tickerFlag, "t", "AAPL", "Ticker of the Stock you want to look up.")
	flag.StringVar(&periodTypeFlag, "pt", "month", "Period Type of the return; e.g. day, month, year.")
	flag.StringVar(&periodFlag, "p", "3", "Number of periodTypes to return in the []FRAME.")
	flag.StringVar(&freqTypeFlag, "ft", "daily", "Frequency Type of the return - Valid fTypes by pType; day: minute / month: daily, weekly / year: daily, weekly, monthly / ytd: daily, weekly.")
	flag.StringVar(&freqFlag, "f", "1", "Frequency of the return in the []FRAME.")
	flag.Parse()

	df, err := data.PriceHistory(tickerFlag, periodTypeFlag, periodFlag, freqTypeFlag, freqFlag)

	if err != nil {
		log.Fatalf(err.Error())
	}

	d := FRAMEToDataSlice(df)
	d.Set()

	fmt.Printf("%+v\n", d)
	fmt.Println(time.Since(start))
}
