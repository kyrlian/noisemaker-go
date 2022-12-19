package main

import "fmt"

//SignalSum is a collection of signals
type SignalSum struct {
	elements [](Signal) //slice
	ampl     Signal     //base+lfoa
	name     string
}

//CONSTRUCTORS
func signalSum_n(n string) SignalSum {
	return SignalSum{elements: []Signal{}, ampl: tf(1.0), name: n}
}
func signalSum() SignalSum {
	return signalSum_n("")
}

func (sum SignalSum) set(elem int, s Signal) SignalSum {
	switch elem {
	case AMPL:
		sum.ampl = s
	default:
		fmt.Printf("WARNING: SignalSum.set:unkown element:%v\n", elem)
	}
	return sum
}
func (sum SignalSum) setampl(s Signal) SignalSum {
	return sum.set(AMPL, s)
}
//MODIFIERS
func (sum SignalSum) appendSignal(s Signal) SignalSum {
	sum.elements = append(sum.elements, s)
	return sum
}

//GETTERS
func (sum SignalSum) getval(t float64) float64 {
	var r = .0
	for _, signal := range sum.elements {
		r += signal.getval(t)
	}
	return r * sum.ampl.getval(t)
}
