package main

import (
	"fmt"
	"math"
)

//xyPair is a simple x,y struct, for custom shapes
type xyPair struct {
	x float64
	y float64
}

type TimedPair struct {
	x TimedFloat
	y TimedFloat
}

//Oscillator structure
type CustomShape struct {
	//freq   Signal //unused, we use last x
	phase  Signal
	ampl   Signal
	points [](TimedPair) //x is relative to previous, to enable to change part of the shape only - y is absolute
	name   string
}

//CONSTRUCTORS
func NewTimedPair(x TimedFloat, y TimedFloat) TimedPair {
	return TimedPair{x, y}
}

func tp(x float64, y float64) TimedPair { //alias with conversion float-timedFloat
	return NewTimedPair(tf(x), tf(y))
}

func customShape(phase Signal, ampl Signal, points [](TimedPair)) CustomShape {
	return CustomShape{phase, ampl, points, ""}
}

//CONSTRUCTORS - helpers
func timedPairList_xy(xylist [](xyPair)) []TimedPair { //convert list of xyPair to list of TimedPair
	var res = [](TimedPair){}
	for _, xy := range xylist {
		res = append(res, tp(xy.x, xy.y))
	}
	return res
}

func timedPairList_yd(ylist [](float64), duration float64) [](TimedPair) { //build an array of xy from a list of y and a total duration
	var nbpoints = len(ylist)
	var tfx = tf(duration / (float64(nbpoints - 1)))
	var res = [](TimedPair){}
	res = append(res, NewTimedPair(tf(.0), tf(ylist[0])))
	for i := 1; i < nbpoints; i++ {
		res = append(res, NewTimedPair(tfx, tf(ylist[i]))) //x is relative to previous
	}
	return res
}

func customShape_tp(points [](TimedPair)) CustomShape {
	//lastx := points[len(points)-1].x
	return customShape(tf(.0), tf(1.0), points)
}
func customShape_xy(points [](xyPair)) CustomShape {
	//var lastx = points[len(points)-1].x
	//var freq = 1 / lastx
	//return customShape(tf(freq), tf(.0), tf(1.0), timedPairList_xy(points))
	return customShape_tp(timedPairList_xy(points))
}
func customShape_yd(ylist [](float64), duration float64) CustomShape {
	return customShape_tp(timedPairList_yd(ylist, duration))
}

//SETTERS
func (cs CustomShape) set(elem int, s Signal) CustomShape {
	switch elem {
	case PHASE:
		cs.phase = s
	case AMPL:
		cs.ampl = s
	default:
		fmt.Printf("WARNING: CustomShape.set:unkown element:%v\n", elem)
	}
	return cs
}
func (cs CustomShape) setName(n string) CustomShape {
	cs.name = n
	return cs
}

//GETTERS
func (cs CustomShape) getperiod(t float64) float64 {
	points := cs.points
	period := .0
	for i := 0; i < len(points); i++ {
		period += points[i].x.getval(t)
	}
	return period
}

func (cs CustomShape) getval(t float64) float64 {
	//freq := cs.freq.getval(t)
	phase := cs.phase.getval(t)
	points := cs.points
	period := cs.getperiod(t)
	//fmt.Printf("	getcustomval:period:%v\n", period)
	tmod := math.Mod(t+(phase*period), period) //adjusted time: O-p
	previousx := points[0].x.getval(tmod)
	previousy := points[0].y.getval(tmod)
	y := .0
	for i := 1; i < len(points); i++ {
		point := points[i]
		deltax := point.x.getval(t) //x is relative
		nextx := previousx + deltax
		nexty := point.y.getval(t)
		if previousx <= tmod && tmod < nextx {
			y = previousy + (tmod-previousx)/(deltax)*(nexty-previousy)
			//fmt.Printf("	getcustomval;x;%v;y;%v;\n", x, y)
		}
		previousx = nextx
		previousy = nexty
	}
	return y * cs.ampl.getval(t)
}
