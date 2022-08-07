package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/samjtro/go-tda/data"
)

var (
	m  sync.Mutex
	m1 sync.Mutex
	m2 sync.Mutex
	m3 sync.Mutex
	m4 sync.Mutex
	m5 sync.Mutex
	m6 sync.Mutex
)

type DATA struct {
	Price          float64
	RMA            float64
	EMA            float64
	RSI            float64
	VWAP           float64
	MACD           float64
	Chaikin        float64
	BollingerUpper float64
	BollingerLower float64
	IMI            float64
	MFI            float64
	PCR            float64
	OI             float64
}

type DataSlice []DATA

func main() {
	start := time.Now()
	df, err := data.PriceHistory("TSLA", "month", "3", "daily", "1")

	if err != nil {
		log.Fatalf(err.Error())
	}

	d := FRAMEToDataSlice(df)
	d.Set(df)
	fmt.Println(d)
	fmt.Println(time.Since(start))
}

func FRAMEToDataSlice(df []data.FRAME) DataSlice {
	d := DataSlice{}

	for _, x := range df {
		d1 := DATA{
			Price: x.Close,
		}

		d = append(d, d1)
	}

	return d
}

func (d DataSlice) Set(df []data.FRAME) {
	wg := new(sync.WaitGroup)
	wg.Add(5)
	go d.RMA(3, df, wg)
	go d.EMA(10, wg)
	go d.RSI(wg)
	go d.VWAP(df, wg)
	go d.MACD(wg)
	wg.Wait()
	// go d.Chaikin(21, df)
}

// Calculates RMA, creates []DATA structure for use in the rest of the App, returns it
func (d DataSlice) RMA(n float64, df []data.FRAME, wg *sync.WaitGroup) {
	defer wg.Done()
	m.Lock()

	for i, _ := range df {
		for j, _ := range d {
			sum := 0.0

			if i >= int(n) {
				for a := int(n); a != 0; a-- {
					sum += df[i-a].Close
				}

				d[j].RMA = sum / n
			}
		}
	}

	m.Unlock()
}

// Calculates EMA, adds to []DATA from RMA, return the []DATA
func (d DataSlice) EMA(n float64, wg *sync.WaitGroup) {
	defer wg.Done()
	m1.Lock()

	mult := 2 / (n + 1)

	for i, _ := range d {
		sum := 0.0

		if i == int(n) {
			for a := 2; a != int(n)+1; a++ {
				sum += d[i-a].Price
			}

			a := d[i-1].Price - sum/n
			b := mult + sum/n
			ema := a * b

			d[i].EMA = ema
		} else if i > int(n) {
			prevEma := d[i-2].EMA
			ema := (d[len(d)-1].Price-prevEma)*mult + prevEma

			d[i].EMA = ema
		}
	}

	m1.Unlock()
}

// Calculates EMA, adds to []DATA from RMA, return the []DATA
func Ema(n float64, d DataSlice) DataSlice {
	m1.Lock()

	mult := 2 / (n + 1)

	for i, _ := range d {
		sum := 0.0

		if i == int(n) {
			for a := 2; a != int(n)+1; a++ {
				sum += d[i-a].Price
			}

			a := d[i-1].Price - sum/n
			b := mult + sum/n
			ema := a * b

			d[i].EMA = ema
		} else if i > int(n) {
			prevEma := d[i-2].EMA
			ema := (d[len(d)-1].Price-prevEma)*mult + prevEma

			d[i].EMA = ema
		}
	}

	m1.Unlock()

	return d
}

func (d DataSlice) RSI(wg *sync.WaitGroup) {
	defer wg.Done()
	m2.Lock()

	gain := []float64{}
	loss := []float64{}
	var avgGain, avgLoss float64

	for i, _ := range d {
		if i > 0 {
			diff := d[i].Price - d[i-1].Price

			if diff < 0 {
				loss = append(loss, diff)
			} else {
				gain = append(gain, diff)
			}
		}

		if i > 14 {
			avgGain = AverageGainLoss(i, d, gain)
			avgLoss = AverageGainLoss(i, d, loss)
			rs := avgGain / avgLoss
			d[i].RSI = 100 - (100 / (1 + rs))
		}
	}

	m2.Unlock()
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
		avgGainLoss = ((initialAvgGainLoss * 13) + (d[i].Price - d[i-1].Price)) / 14
	} else {
		initialAvgGainLoss = InitialAverageGainLoss(data)
		avgGainLoss = (initialAvgGainLoss*float64(len(data)-1) + (d[i].Price - d[i-1].Price)) / float64(len(data))
	}

	return avgGainLoss
}

func (d DataSlice) VWAP(df []data.FRAME, wg *sync.WaitGroup) {
	defer wg.Done()
	m5.Lock()

	var averagePrice float64

	for i, x := range d {
		for _, y := range df {
			averagePrice += (x.Price + y.Lo + y.Hi)
			vwap := averagePrice / y.Volume

			d[i].VWAP = vwap
		}
	}

	m5.Unlock()
}

func (d DataSlice) MACD(wg *sync.WaitGroup) {
	defer wg.Done()
	var macd float64

	twentySixDayEMA := Ema(26, d)
	twelveDayEMA := Ema(12, d)

	m6.Lock()

	for _, x := range twentySixDayEMA {
		for _, y := range twelveDayEMA {
			for i, _ := range d {
				macd = y.EMA - x.EMA
				d[i].MACD = macd
			}
		}
	}

	m6.Unlock()
}

func (d DataSlice) Chaikin(p int, df []data.FRAME, wg *sync.WaitGroup) {
	defer wg.Done()
	var Helper DataSlice

	for _, x := range d {
		for _, y := range df {
			n := ((x.Price - y.Lo) - (y.Hi - x.Price)) / (y.Hi - y.Lo)
			m := n * (y.Volume * float64(p))
			adl := (m * float64(p-1)) + (m * float64(p))
			d1 := DATA{Price: adl}

			Helper = append(Helper, d1)
		}
	}

	threeDayEMA := Ema(3, Helper)
	tenDayEMA := Ema(10, Helper)

	for _, x := range threeDayEMA {
		for _, y := range tenDayEMA {
			for i, _ := range d {
				d[i].Chaikin = x.EMA - y.EMA
			}
		}
	}
}

func BollingerUpper() {

}

func BollingerLower() {

}

func IMI() {

}

func MFI() {

}

func PCR() {

}

func OI() {

}
