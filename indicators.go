package main

import (
	"log"
	"sync"

	"github.com/montanaflynn/stats"
)

/*

RealTime

*/

/*

PriceHistory

*/

// Set SMA values for the given DataSlice
func (d DataSlice) SMA(n int, wg *sync.WaitGroup) {

}

func Sma(n int, d DataSlice) DataSlice {
	for i := range d {
		sum := 0.0
		if i >= n {
			for a := 0; a < n; a++ {
				sum += (d[i-a].Close + d[i-a].Hi + d[i-a].Lo) / 3
			}

			d[i].SMA = sum / float64(n)
		}
	}

	return d
}

// Set RMA values for the given DataSlice
func (d DataSlice) RMA(n float64, wg *sync.WaitGroup) {
	defer wg.Done()
	m.Lock()

	for i := range d {
		sum := 0.0

		if i >= int(n) {
			for a := int(n); a != 0; a-- {
				sum += d[i-a].Close
			}

			d[i].RMA = sum / n
		}
	}

	m.Unlock()
}

// Set EMA values for the given DataSlice
func (d DataSlice) EMA(n float64, wg *sync.WaitGroup) {
	defer wg.Done()
	m1.Lock()

	mult := 2 / (n + 1)

	for i := range d {
		sum := 0.0

		if i == int(n) {
			for a := 2; a != int(n)+1; a++ {
				sum += d[i-a].Close
			}

			a := d[i-1].Close - sum/n
			b := mult + sum/n
			ema := a * b

			d[i].EMA = ema
		} else if i > int(n) {
			prevEma := d[i-2].EMA
			ema := (d[len(d)-1].Close-prevEma)*mult + prevEma

			d[i].EMA = ema
		}
	}

	m1.Unlock()
}

func Ema(n float64, d DataSlice) DataSlice {
	m1.Lock()

	mult := 2 / (n + 1)

	for i := range d {
		sum := 0.0

		if i == int(n) {
			for a := 2; a != int(n)+1; a++ {
				sum += d[i-a].Close
			}

			a := d[i-1].Close - sum/n
			b := mult + sum/n
			ema := a * b

			d[i].EMA = ema
		} else if i > int(n) {
			prevEma := d[i-2].EMA
			ema := (d[len(d)-1].Close-prevEma)*mult + prevEma

			d[i].EMA = ema
		}
	}

	m1.Unlock()

	return d
}

// Calculate Pivot, Resistance & Support Points
func (d DataSlice) PivotPoints(wg *sync.WaitGroup) {
	defer wg.Done()
	m7.Lock()

	for i, x := range d {
		pivotPoint := (x.Hi + x.Lo + x.Close) / 3
		firstResistance := (2 * pivotPoint) - x.Lo
		firstSupport := (2 * pivotPoint) - x.Hi
		secondResistance := pivotPoint + (x.Hi - x.Lo)
		secondSupport := pivotPoint - (x.Hi - x.Lo)
		thirdResistance := x.Hi + (2 * (pivotPoint - x.Lo))
		thirdSupport := x.Lo - (2 * (x.Hi - pivotPoint))

		d[i].PivotPoint = pivotPoint
		d[i].ResistancePoints = []float64{firstResistance, secondResistance, thirdResistance}
		d[i].SupportPoints = []float64{firstSupport, secondSupport, thirdSupport}
	}
}

// Set RSI values for the given DataSlice
func (d DataSlice) RSI(wg *sync.WaitGroup) {
	defer wg.Done()
	m2.Lock()

	var gain, loss []float64
	var avgGain, avgLoss float64

	for i := range d {
		if i > 0 {
			diff := d[i].Close - d[i-1].Close

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

// Set VWAP values for the given DataSlice
func (d DataSlice) VWAP(wg *sync.WaitGroup) {
	defer wg.Done()
	m5.Lock()

	for i, x := range d {
		typicalPrice := (x.Close + x.Lo + x.Hi) / 3
		vwap := (typicalPrice * x.Volume) / x.Volume

		d[i].VWAP = vwap
	}

	m5.Unlock()
}

// Set MACD values for the given DataSlice
func (d DataSlice) MACD(wg *sync.WaitGroup) {
	wg1 := new(sync.WaitGroup)
	defer wg.Done()

	m6.Lock()
	d2 := make(DataSlice, len(d))
	copy(d2, d)
	wg1.Add(1)

	go d2.EMA(26, wg1)

	wg1.Wait()

	for i, x := range d {
		for _, y := range d2 {
			macd := x.EMA - y.EMA
			d[i].MACD = macd
		}
	}

	m6.Unlock()
}

// Set Chaikin values for the given DataSlice
func (d DataSlice) Chaikin(p int, wg *sync.WaitGroup) {
	defer wg.Done()
	var Helper DataSlice

	for _, x := range d {
		n := ((x.Close - x.Lo) - (x.Hi - x.Close)) / (x.Hi - x.Lo)
		m := n * (x.Volume * float64(p))
		adl := (m * float64(p-1)) + (m * float64(p))
		d1 := DATA{Close: adl}

		Helper = append(Helper, d1)
	}

	threeDayEMA := Ema(3, Helper)
	tenDayEMA := Ema(10, Helper)

	for _, x := range threeDayEMA {
		for _, y := range tenDayEMA {
			for i := range d {
				d[i].Chaikin = x.EMA - y.EMA
			}
		}
	}
}

// Set Bollinger Bands for the given DataSlice
func (d DataSlice) BollingerBands(wg *sync.WaitGroup) {
	defer wg.Done()
	var stdDevHelper []float64
	twentyDaySMA := Sma(20, d)

	for i, x := range d {
		for _, y := range twentyDaySMA {
			stdDevHelper = append(stdDevHelper, (x.Close+x.Hi+x.Lo)/3)

			if len(stdDevHelper) == 21 {
				stdDevHelper = stdDevHelper[1:]
			}

			if i >= 20 {
				stdDev, err := stats.StandardDeviation(stdDevHelper)

				if err != nil {
					log.Fatalf(err.Error())
				}

				upperBand := y.SMA + (2 * stdDev)
				lowerBand := y.SMA - (2 * stdDev)
				d[i].BollingerBands = []float64{upperBand, x.SMA, lowerBand}
			} else {
				d[i].BollingerBands = []float64{0, x.SMA, 0}
			}

		}
	}
}

// Set IMI values for the given DataSlice
func (d DataSlice) IMI(wg *sync.WaitGroup) {
	defer wg.Done()
	m2.Lock()

	var gain, loss []float64
	var sumOfGain, sumOfLoss float64

	for i := range d {
		if i > 0 {
			diff := d[i].Close - d[i-1].Close

			if diff < 0 {
				loss = append(loss, diff)
			} else {
				gain = append(gain, diff)
			}
		}

		if i > 14 {
			for i := range gain {
				sumOfGain += gain[i]
			}

			for i := range loss {
				sumOfLoss += loss[i]
			}

			im := sumOfGain / (sumOfGain + sumOfLoss)
			imi := im * 100
			d[i].IMI = imi
		}
	}

	m2.Unlock()
}

// Set MFI values for the given DataSlice
func MFI() {

}

// Set PCR values for the given DataSlice
func PCR() {

}

// Set OI values for the given DataSlice
func OI() {

}
