package main

import (
	"fmt"
	"math"
	"math/rand"
)

//type shapeType int

//shape type
const (
	SIN int = iota
	SQR
	TRI
	SAW
	ISAW
	FLAT
	NOISE
	PULSE

//	CUSTOM
)

//Oscillator elements
const (
	FREQ int = iota
	PHASE
	AMPL
	WIDTH
)

//Oscillator structure
type Oscillator struct {
	shape int
	freq  Signal
	phase Signal
	ampl  Signal
	width Signal // optional, to set the width of the top of the PULSE shape
	name  string
}

//CONSTRUCTORS
func oscillator_full(shape int, freq Signal, phase Signal, ampl Signal, width Signal) Oscillator {
	return Oscillator{shape, freq, phase, ampl, width, ""}
}
func oscillator(shape int, freq Signal, phase Signal, ampl Signal) Oscillator {
	return oscillator_full(shape, freq, phase, ampl, tf(.1))
}
func oscillator_sf(shape int, freq float64) Oscillator {
	return oscillator(shape, tf(freq), tf(.0), tf(1.0))
}
func oscillator_sfpa(shape int, freq float64, phase float64, ampl float64) Oscillator {
	return oscillator(shape, tf(freq), tf(phase), tf(ampl))
}

//CONSTRUCTORS dedicated to shapes
func oscillator_pulse(freq Signal, phase Signal, ampl Signal, width Signal) Oscillator {
	return oscillator_full(PULSE, freq, phase, ampl, width).setname("PULSE")
}
func oscillator_noise(ampl Signal) Oscillator {
	return oscillator(NOISE, tf(1.0), tf(1.0), ampl).setname("NOISE")
}

// See customshape.go for CUSTOM constructors

//SETTERS
func (o Oscillator) setname(s string) Oscillator {
	o.name = s
	return o
}

func (o Oscillator) set(elem int, s Signal) Oscillator {
	switch elem {
	case FREQ:
		o.freq = s
	case PHASE:
		o.phase = s
	case AMPL:
		o.ampl = s
	case WIDTH:
		o.width = s
	default:
		fmt.Printf("WARNING: Oscillator.set:unkown element:%v\n", elem)
	}
	return o
}

//GETERS
func (o Oscillator) getval(t float64) float64 {
	freq := o.freq.getval(t)
	phase := o.phase.getval(t)
	period := 1.0 / freq                 //period
	tmod := math.Mod(t, period)          //adjusted time: O-p
	xmod := math.Mod(tmod*freq+phase, 1) //O-1 - All shapes have a period of 1
	//fmt.Printf("	Oscillator:getval:x:%v,xmod:%v\n", x, xmod)
	y := 0.0
	switch o.shape { //All shapes have a period of 1
	case SIN:
		y = math.Sin(2.0 * math.Pi * xmod)
	case FLAT:
		y = 1.0
	case SQR:
		if xmod < .5 {
			y = 1.0
		} else {
			y = -1.0
		}
	case PULSE:
		if xmod < o.width.getval(t) { //with width=.5 it's just a square
			y = 1.0
		} else {
			y = -1.0
		}
	case SAW: //ramp up
		y = -1.0 + 2*xmod
	case ISAW: //ramp down
		y = 1.0 - 2*xmod
	case TRI: //ramp up and down
		if xmod < .5 {
			y = -1.0 + 4.0*xmod
		} else {
			y = 3.0 - 4.0*xmod
		}
	case NOISE:
		y = 2.0*rand.Float64() - 1.0 //rand is [0.0,1.0)
		/* 	case CUSTOM:
		if o.customshape != nil {
			y = getcustomshapeval(o.customshape, xmod, t)
		} */
	default:
		fmt.Printf("WARNING: Oscillator.getval:unkown shape:%v\n", o.shape)
	}
	return y * o.ampl.getval(t) //can be negative - ex for LFOs
}
