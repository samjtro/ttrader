package main

import (
	"sync"

	"github.com/samjtro/go-tda/data"
)

// Marshal the []FRAME returned by go-tda calls into a DataSlice
func FRAMEToDataSlice(df []data.FRAME) DataSlice {
	d := DataSlice{}

	for _, x := range df {
		d1 := DATA{
			Close:  x.Close,
			Hi:     x.Hi,
			Lo:     x.Lo,
			Volume: x.Volume,
		}

		d = append(d, d1)
	}

	return d
}

// Set all indicators for the given DataSlice
func (d DataSlice) Set() {
	wg := new(sync.WaitGroup)
	wg.Add(7)

	go d.PivotPoints(wg)
	go d.RMA(3, wg)
	go d.EMA(12, wg)
	go d.RSI(wg)
	go d.VWAP(wg)
	go d.MACD(wg)
	go d.BollingerBands(wg)
	// go d.Chaikin(21, df)

	wg.Wait()
}

func InitialAverageGainLoss(data []float64) float64 {
	sum := 0.0

	m3.Lock()
	defer m3.Unlock()

	for _, x := range data {
		sum += x
	}

	return sum
}

func AverageGainLoss(i int, d []DATA, data []float64) float64 {
	var initialAvgGainLoss, avgGainLoss float64

	m4.Lock()
	defer m4.Unlock()

	if len(data) >= (i - 13) {
		initialAvgGainLoss = InitialAverageGainLoss(data[(i - 13):])
		avgGainLoss = ((initialAvgGainLoss * 13) + (d[i].Close - d[i-1].Close)) / 14
	} else {
		initialAvgGainLoss = InitialAverageGainLoss(data)
		avgGainLoss = (initialAvgGainLoss*float64(len(data)-1) + (d[i].Close - d[i-1].Close)) / float64(len(data))
	}

	return avgGainLoss
}
