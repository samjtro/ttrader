package main

import (
	"log"

	"github.com/go-schwab/trader"
)

func main() {
	agent := trader.Initiate()
	data, err := agent.GetPriceHistory("AAPL", "month", "3", "daily", "1", "", "")
	if err != nil {
		log.Fatalf(err.Error())
	}
	d := CandleToDataSlice(data)
	d.Set()
}
