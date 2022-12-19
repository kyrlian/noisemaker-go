package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

var valmax = 1.0
var SAMPLERATE = 44100 //Only Capitalized variable are visible outside the package - I used SAMPLERATE for readability
var bufferseconds = 1

//Beep Streamer struct
type TrackStreamer struct {
	trackpointer   Signal
	samplePosition int
	//nbsamples      int //max number of sample based on track's end time and SAMPLERATE
}

func trackStreamer(signal Signal) TrackStreamer { //Infinite streamer
	return TrackStreamer{signal, 0}
}

func timedStreamer(signal Signal, nbseconds int) beep.Streamer { //Finite streamer based on track duration
	trackstr := trackStreamer(signal)
	return beep.Take( SAMPLERATE * nbseconds, &trackstr) 
}

//beep.Streamer interface requires Stream and Err
func (str *TrackStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	var tstart = time.Now()
	for i := range samples {
		time := float64(str.samplePosition) / float64(SAMPLERATE) //seconds
		val := str.trackpointer.getval(time)
		if val > valmax {
			valmax = val //store and display the max value seen
			fmt.Printf("WARNING - sampler:valmax=%v\n", val)
		}
		samples[i][0] = val
		samples[i][1] = val
		str.samplePosition++
	}
	if false { //activate to show computing stats
		var elapsed = time.Since(tstart)
		var microseconds = 1000.0 * 1000.0 * len(samples) / SAMPLERATE
		var duration = time.Duration(microseconds) * time.Microsecond
		fmt.Printf("	sampler.Stream:Computed %v samples for %v in %v\n", len(samples), duration, elapsed)
	}
	return len(samples), true
}

func (str TrackStreamer) Err() error {
	return nil
}

func saveToWav(signal Signal, filename string, nbseconds int) {
	fmt.Printf("Saving to %v\n", filename+".wav")
	file, fileerr := os.CreateTemp("out", filename+"_*.wav")
	if fileerr != nil {
		log.Panic(fileerr)
	}
	//trackstr := trackStreamer(signal)
	//timedStreamer := beep.Take(trackstr.nbsamples, &trackstr)
	timedStreamer := timedStreamer(signal,nbseconds)
	format := beep.Format{SampleRate: beep.SampleRate(SAMPLERATE), NumChannels: 2, Precision: 2}
	writeerr := wav.Encode(file, timedStreamer, format)
	if writeerr != nil {
		log.Panic(writeerr)
	}
}

//Run
func runSampler(signal Signal, nbseconds int) {
	fmt.Printf("Preparing speaker\n")
	var sampleRate = beep.SampleRate(SAMPLERATE)
	var buffSize = bufferseconds * sampleRate.N(time.Second) //buffer for 5 seconds
	fmt.Printf("	sampler.initSpeaker:buffSize=%v\n", buffSize)
	speaker.Init(sampleRate, buffSize)
	fmt.Printf("Preparing streamer\n")
	//var trackstr = trackStreamer(signal)
	var timedStreamer = timedStreamer(signal,nbseconds)
	fmt.Printf("Playing speaker\n")
	speaker.Play(timedStreamer)
	select {}
}
