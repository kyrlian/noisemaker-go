package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

func getMinMax(signal Signal, nbseconds int, samplerate int) (float64, float64){
	ymin:=signal.getval(.0)
	ymax:=ymin
	step := 1 / float64(samplerate)
	xmax := float64(nbseconds)
	for x := .0 ;x<xmax; x+= step{
		y := signal.getval(x)
		ymin=math.Min(ymin,y)
		ymax=math.Max(ymax,y)
	}
	return ymin,ymax
}

func plotsignal(signal Signal, filename string, nbseconds int, samplerate int) {
	fmt.Println("Ploting signal to out/"+filename+".png")
	//nbseconds := 30
	nbsamples := samplerate * nbseconds / 10
	signalfunc := plotter.NewFunction(func(x float64) float64 { return signal.getval(x) })
	signalfunc.Color = color.RGBA{B: 255, A: 255}
	signalfunc.Samples = nbsamples

	p := plot.New()
	p.Add(signalfunc)
	p.Legend.Add("signal", signalfunc)
	p.X.Min = .0
	p.X.Max = float64(nbseconds)
	//p.Y.Min = -2.0
	//p.Y.Max = 2.0
	p.Y.Min,p.Y.Max = getMinMax(signal ,  nbseconds , samplerate)
	err := p.Save(2048, 1024, "out/"+filename+".png")
	if err != nil {
		log.Panic(err)
	}
}
