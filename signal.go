package main

//Signal interface - anything that gives a val for a time - both oscillator and track are signals
type Signal interface {
	getval(t float64) float64
	//amplify(factor TimedFloat) Signal 
}

// Audio signal is expected -1 to 1
// Control signal is V/Oct like (-1 = lower freq by 1 octave)
// Trigger, Gate or Clock are pulses - unused yet
