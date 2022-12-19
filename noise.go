package main

//TODO: proper packages
import (
	"fmt"
)

func main() {
	fmt.Println("MAKE SOME NOISE")

	fmt.Println("Preparing track")
	var signal = example_harmonicsTuning(22.5, 10)

	var nbsecs = 20
	plotsignal(signal, signal.name, nbsecs, 44100/10) //plot track for 30 seconds
	saveToWav(signal, signal.name, nbsecs) //save to wave file
	runSampler(signal, nbsecs)//play on speakers
}
