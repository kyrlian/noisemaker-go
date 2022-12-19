package main

import (
	"fmt"
	"math"
)

//TrackElement with a start and end
type TrackElement struct {
	signal Signal
	start  float64
	end    float64
	ampl   Signal //TODO maybe remove ampl from here, we can use .amplify on the signal if needed
}

//Track is a collection of track elements
type Track struct {
	elements [](TrackElement) //slice
	start    float64
	end      float64
	ampl     Signal     //base+lfoa
	name     string
}

//CONSTRUCTORS
func NewTrackElement(signal Signal, start float64, end float64, ampl Signal) TrackElement {
	if end < start {
		end = start
	}
	return TrackElement{signal, start, end, ampl}
}
func NewTrackNamed(n string) Track {
	return Track{ []TrackElement{},  .0,  .0, tf(1.0), n}
}
func NewTrack() Track {
	return NewTrackNamed("")
}

//MODIFIERS
func (tr Track) set(elem int, s Signal) Track {
	switch elem {
	case AMPL:
		tr.ampl = s
	default:
		fmt.Printf("WARNING: SignalSum.set:unkown element:%v\n", elem)
	}
	return tr
}
func (tr Track) appendSignal(signal Signal, start float64, end float64, ampl Signal) Track {
	tr.elements = append(tr.elements, NewTrackElement(signal, start, end, ampl))
	tr.start = math.Min(tr.start, start)
	tr.end = math.Max(tr.end, end)
	return tr
}

//GETTERS
func (tr Track) getval(t float64) float64 {
	var r = .0
	//fmt.Printf("	tr.elements.len: %v\n", len(tr.elements))
	for _, elem := range tr.elements {
		if t >= elem.start && t < elem.end { //this signal is active
			var subt = t - elem.start
			v := elem.signal.getval(subt)
			a := elem.ampl.getval(subt)
			r += v * math.Max(0, a)
			//fmt.Printf("	tr.getval.t,r: %v , %v\n", t, r)
		}
	}
	return r * tr.ampl.getval(t)
}
