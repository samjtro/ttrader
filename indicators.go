package main

import (
	"fmt"
	"strconv"
	//"math"
	. "github.com/samjtro/go-tda/data"
)

func RMA(n float64,df []float64) float64 {
	var sum float64
	for _,x := range df {
		sum += x
	}
	return sum/n
}

//func EMA() {}
//func RSI() {}
//func VWAP() {}
//func MACD() {}
//func CHAIKIN() {}

func Set(df []FRAME) []float64 {
	var arr []float64
	var rmas []float64
	for i,_ := range df {
		if(i > 3) {
			c1,_ := strconv.ParseFloat(df[i-1].CLOSE,8)
			c2,_ := strconv.ParseFloat(df[i-2].CLOSE,8)
			c3,_ := strconv.ParseFloat(df[i-3].CLOSE,8)
			c4,_ := strconv.ParseFloat(df[i-4].CLOSE,8)
			arr = append(arr,c1,c2,c3,c4)
			rma := RMA(4.0,arr))
		}

		arr = nil
	}

	return rmas
}
