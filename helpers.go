package main

import (
	//	"fmt"
	"fmt"
	"math"
)

func root(x float64, n int) float64 {
	var lower = .0
	var upper = x
	var r = .0
	for upper-lower >= 0.000000001 {
		r = (upper + lower) / 2.0
		var temp = math.Pow(r, float64(n))
		if temp > x {
			upper = r
		} else {
			lower = r
		}
	}
	return r
}

var semiToneConst = root(2.0, 12) //1,05946309435929
var toneConst = root(2.0, 7)      //1,104089513673812
var _ = toneConst + semiToneConst //Avoid not used error

func getSemiToneFreq(fstart float64, nsemitones int) float64 {
	var nfreq = fstart * math.Pow(semiToneConst, float64(nsemitones))
	//fmt.Printf("	getFreq(%v,%v)=%v\n", fstart, nsemitones, nfreq)
	return nfreq
}

func accord3(baseFreq float64, gap1 int, gap2 int) SignalSum { //3,4
	return signalSum().
		appendSignal(oscillator_sf(SIN, baseFreq).amplify(tf(.7))).
		appendSignal(oscillator_sf(SIN, getSemiToneFreq(baseFreq, gap1)).amplify(tf(.6))).
		appendSignal(oscillator_sf(SIN, getSemiToneFreq(baseFreq, gap1+gap2)).amplify(tf(.5)))
}

func accordMineur(fstart float64) SignalSum {
	return accord3(fstart, 3, 4)
}
func accordMajeur(fstart float64) SignalSum {
	return accord3(fstart, 4, 3)
}

func harmonics(baseFreq float64, nharmonics int) SignalSum {
	var sum = signalSum_n("harmonics" + fmt.Sprint(nharmonics))
	for i := 1; i <= nharmonics; i++ {
		var f2pi = math.Pow(2, float64(i))
		sum = sum.appendSignal(oscillator_sfpa(SIN, baseFreq*f2pi, .0, 1.0/f2pi))
	}
	return sum
}
