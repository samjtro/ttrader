package main

import (
	"fmt"
	"strconv"
	. "github.com/samjtro/go-tda/data"
)

func main() {
	fmt.Println(Set(PriceHistory("TSLA","month","3","daily","1")))
}

type RMA struct {
	F		data.FRAME
	RMA		float64
	EMA		float64
	//RSI
	//VWAP
	//MACD
	//CHAIKIN
}

// SMA returns a simple moving average of length n,
// for each FRAME in data
func RMA(n float64,data data.[]FRAME) []FRAME {
	var df []FRAME
	for i,x := range data {
		sum := 0.0
		if i >= int(n) {
			for a:=int(n); a!=0; a-- {
				c,_ := strconv.ParseFloat(data[i-a].CLOSE,8)
				sum += c
			}

			rma := sum/n

			f := FFRAME{
				f:	x,
				rma:	rma,
			}

			df = append(df,f)
		}
	}

	return df
}

// EMA returns the exponential moving average of length n for multiplying factor mult,
// for each FRAME of data
func EMA(n,mult float64,data data.[]FRAME) []FRAME {
	var df []FFRAME
	for i,x := range data {
		sum := 0.0
		if i == int(n) {
			c1,_ := strconv.ParseFloat(data[i-1].CLOSE,8)
			for a:=2; a!=int(n)+1; a++ {
				c,_ := strconv.ParseFloat(data[i-a].CLOSE,8)
				sum += c
			}
			a := c1-sum/n
			b := mult+sum/n

			ema := a * b

			f := FFRAME{
				f:	x,
				rma:	ema,
			}

			df = append(df,f)
		} else if i > int(n) {
			prevEma := df[len(df)-1].rma
			c,_ := strconv.ParseFloat(df[len(df)-1].f.CLOSE,8)
			ema := (c-prevEma)*mult+prevEma

			f := FFRAME{
				f:	x,
				rma:	ema,
			}

			df = append(df,f)
		}
	}

	return df
}

//func RSI() {}
//func VWAP() {}
//func MACD() {}
//func CHAIKIN() {}

func Set(df []FRAME) []FFRAME {
	return EMA(4,.4,df)
}

