package main

type TimedFloat struct {
	base float64
	ampl Signal
}

//CONSTRUCTORS
func timedFloat(base float64, ampl Signal) TimedFloat {
	return TimedFloat{base, ampl}
}

func tf(base float64) TimedFloat { //short alias
	return timedFloat(base, nil)
}

//GETTERS
func (tf TimedFloat) getval(t float64) float64 {
	r := tf.base
	if tf.ampl != nil {
		r *= tf.ampl.getval(t)
	}
	return r
}
