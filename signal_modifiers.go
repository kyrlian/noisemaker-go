package main

import "math"

//FILTER - filter(low,high)
type Filter struct {
	in   Signal
	low  float64
	high float64
}

func filter(s Signal, low float64, high float64) Filter {
	return Filter{s, low, high}
}
func (o Oscillator) filter(low float64, high float64) Filter {
	return filter(o, low, high)
}
func (t Track) filter(low float64, high float64) Filter {
	return filter(t, low, high)
}
func (f Filter) getval(t float64) float64 {
	var v = f.in.getval(t)
	if v < f.low {
		v = f.low
	} else if v > f.high {
		v = f.high
	}
	return v
}
func demo_filter() Signal {
	var s_dotfilter = oscillator_sf(SIN, 120).filter(1.0, 2.0)
	var s_funcfilter = filter(oscillator_sf(SIN, 120), 1.0, 2.0)
	return NewTrack().
		appendSignal(s_dotfilter, .0, 1.0, tf(.0)).
		appendSignal(s_funcfilter, .0, 1.0, tf(.0)).
		filter(.0, .9)
}

var _ = demo_filter() //Avoid not used error

//INVERT - invert() - takes an input signal and gives 1/signal(t)
/* type Inverter struct {
	in Signal
}
func (i Inverter) getval(t float64) float64 {
	return 1 / i.in.getval(t)
}
func invert(s Signal) Inverter {
	return Inverter{s}
}
func (tf TimedFloat) invert() Inverter {
	return invert(tf)
} */

//POWER - power(base) - takes an input signal and gives 2^signal(t)
type Power struct {
	in   Signal
	base float64
}

//with base = 2, the output signal of power is a v/oct - raising input by 1 doubles the output, raising 1 octave
//with base = semiToneConst = root(2.0, 12) = 1,05946309435929, the output signal of power is a v/semitone
//with base = semiToneConst = root(2.0, 12) = 1,05946309435929, the output signal of power is a v/semitone

func (p Power) getval(t float64) float64 {
	return math.Pow(p.base, p.in.getval(t))
}
func power(s Signal, base float64) Power {
	return Power{s, base}
}
func (cs CustomShape) power2() Power {
	return power(cs, 2.0)
}
func demo_power() Signal {
	var s_dotfilter = oscillator(SIN, timedFloat(120, customShape_xy([]xyPair{{.0, .0}, {5.0, 1.0}}).power2()), tf(.0), tf(.1))
	var s_funcfilter = oscillator(SIN, timedFloat(120, power(customShape_xy([]xyPair{{.0, .0}, {5.0, 1.0}}), 2.0)), tf(.0), tf(.1))
	return signalSum().
		appendSignal(s_dotfilter).
		appendSignal(s_funcfilter)
}

var _ = demo_power() //Avoid not used error

//AMPLIFY - amplify(factor) - return factor*signal(t) - use native methods when possible to keep the type
type Amplificator struct {
	in     Signal
	factor TimedFloat
}

func (a Amplificator) getval(t float64) float64 {
	return a.in.getval(t) * a.factor.getval(t)
}
func amplify(in Signal, factor TimedFloat) Amplificator {
	return Amplificator{in, factor}
}
func (in Oscillator) amplify(factor TimedFloat) Oscillator {
	return in.set(AMPL, amplify(in.ampl,factor))
}
func (in Track) amplify(factor TimedFloat) Track {
	return in.set(AMPL, amplify(in.ampl,factor))
}
func (in CustomShape) amplify(factor TimedFloat) CustomShape {
	return in.set(AMPL, amplify(in.ampl,factor))
}
func (in SignalSum) amplify(factor TimedFloat) SignalSum {
	return in.set(AMPL, amplify(in.ampl,factor))
}

//TODO
// shift(delay) - return signal(t-delay)
// echo(delay,factor) - return signal(t)+signal.shift(delay).amplify(factor)
