package main

import (
	//	"fmt"
	"fmt"
	"math"
	"math/rand"
)

func example_simple() Track {
	var o2 = oscillator_sf(PULSE, 440.0)
	o2.setname("pulse 440")
	return NewTrack().appendSignal(o2, .0, 5.0, tf(.8)).amplify(tf(.5))
}

var _ = example_simple() //Avoid not used error

func example_majeurmineur() Track {
	var lfoa = customShape_yd([](float64){.0, .8, .4, .4, .0}, .5).amplify(tf(1 / 1.3))
	var envshape = timedPairList_yd([](float64){.0, .8, .4, .4, .0}, .5)
	var lfoa2 = customShape(tf(.5), tf(1/1.3), envshape)
	var tmajeur = accordMajeur(440.0)
	var tmineur = accordMineur(440.0)
	return NewTrack().
		appendSignal(tmajeur, .0, 5, timedFloat(.8, &lfoa)).
		appendSignal(tmineur, .0, 5.0, timedFloat(.8, &lfoa2))
}

var _ = example_majeurmineur() //Avoid not used error

func example_enveloppe() Track {
	var lfoa = customShape_xy([]xyPair{{.0, .0}, {.1, .8}, {.1, .5}, {.2, .0}, {.1, .0}}) //x is relative
	var o2 = oscillator(SIN, tf(440.0), tf(0), timedFloat(.8, &lfoa))
	return NewTrack().appendSignal(o2, .0, 5.0, tf(1.0))
}

var _ = example_enveloppe() //Avoid not used error

func example_harmonicsTuning(baseFreq float64, nharmonics int) SignalSum {
	var sum = signalSum_n("harmonicsTuning" + fmt.Sprint(nharmonics))
	var outOfTuneMax = .5     //ratio of basefreq
	var tuningDuration = 10.0 //seconds
	var stableDuration = 10.0 //seconds
	for i := 1; i <= nharmonics; i++ {
		var f2pi = math.Pow(2, float64(i))
		var n1Freq = baseFreq * f2pi
		var n2Freq = getSemiToneFreq(n1Freq, 4)
		var n3Freq = getSemiToneFreq(n1Freq, 7)
		//fmt.Printf("	harmonicsTuning:n1Freq:%v	,n2Freq:%v	,n3Freq:%v\n", n1Freq, n2Freq, n3Freq)
		var outOfTuneStart = (2.0*rand.Float64() - 1.0) * outOfTuneMax //+-outOfTuneMax
		var tuningLfof = customShape_xy([]xyPair{{.0, 1.0 + outOfTuneStart}, {tuningDuration, 1.0}, {stableDuration, 1.0}}).setName("tuningLfof")
		plotsignal(tuningLfof, "tuningLfof_"+fmt.Sprint(i), 2*int(tuningDuration+stableDuration), 4410)
		plotsignal(timedFloat(n1Freq, tuningLfof), "timedFloat_n1Freq_"+fmt.Sprint(i), 2*int(tuningDuration+stableDuration), 4410)

		//fmt.Printf("tuningLfof.getperiod(.0):%v\n", tuningLfof.getperiod(.0))
		var accAmp = tf(1.0 / f2pi / 3)
		sum = sum.
			appendSignal(oscillator(SIN, timedFloat(n1Freq, tuningLfof), tf(.0), accAmp)).
			appendSignal(oscillator(SIN, timedFloat(n2Freq, tuningLfof), tf(.0), accAmp)).
			appendSignal(oscillator(SIN, timedFloat(n3Freq, tuningLfof), tf(.0), accAmp))
	}
	return sum.amplify(tf(.8))
}

var _ = example_harmonicsTuning(22.5, 10) //Avoid not used error

func example_drums() Track {
	//var oHighKicks = oscillator_noise(customShape_xy([]xyPair{{.0, .0}, {.1, .8}, {.1, .1}, {.8, .0}}))
	var tOscs = NewTrack().
		//appendSignal(oHighKicks, .0, 5.0, tf(1.0)).
		appendSignal(oscillator(SIN, tf(55.0), tf(0), tf(.7)), .0, 5.0, tf(1.0)).
		appendSignal(oscillator(SIN, tf(110.0), tf(.5), tf(.6)), .0, 5.0, tf(1.0)).
		appendSignal(oscillator_noise(tf(.1)), .0, 5.0, tf(.08))
	var fslope = customShape_yd([]float64{1.0, .1}, 2.0)                                                          //length of the last part of the enveloppe
	var enveloppe = []TimedPair{tp(.0, .0), tp(.2, .8), tp(.1, .6), tp(.1, .1), {timedFloat(.6, fslope), tf(.0)}} //Envelope of the hit
	var ampl = timedFloat(.8, customShape(tf(0), tf(1.0), enveloppe))
	return NewTrack().appendSignal(tOscs, .0, 5.0, ampl)
}

var _ = example_drums() //Avoid not used error

func example_combined1() Track {
	//intro
	var baseFreq = 22.5
	var finalTrack = NewTrack().appendSignal(harmonics(baseFreq, 6), 0, 5, tf(1.0))
	//bip
	var i = 7
	var f2pi = math.Pow(2, float64(i))
	var lfoaBip = oscillator_pulse(tf(4), tf(0), tf(.9), tf(.2))
	var oBip = oscillator(SIN, tf(baseFreq*f2pi), tf(0), tf(1.0/f2pi))
	var oHighKicks = oscillator(NOISE, tf(0), tf(0), timedFloat(.8, oscillator_pulse(tf(8), tf(.1), tf(.2), tf(.1))))
	finalTrack = finalTrack.
		appendSignal(oBip, float64(i)/5, float64(i), timedFloat(.8, lfoaBip)).
		appendSignal(oHighKicks, float64(i)/5, 5.0, tf(1))
	//finalise
	return finalTrack
}

var _ = example_combined1() //Avoid not used error

func example_engine() SignalSum {
	var tduration = 30.0
	var nbh = 3
	var variableSilence = TimedPair{timedFloat(1.0, customShape_yd([](float64){1.0, .01}, tduration)), tf(.0)} //variable factor of the length of the silence
	var otherPistonHit = TimedPair{tf(.2), tf(.0)}

	var piston1 = harmonics(75, nbh)                                                                                                                                                     //.appendSignal(oscillator_noise(tf(.6)))
	var enveloppe1 = []TimedPair{tp(.0, .0), tp(.1, .8), tp(.1, .0), variableSilence, otherPistonHit, variableSilence, otherPistonHit, variableSilence, otherPistonHit, variableSilence} //Envelope of the hit 1
	var cs1 = customShape(tf(0), tf(1.0), enveloppe1)
	var ampl1 = timedFloat(.8, cs1) //global enveloppe uses the hit enveloppe, but with a viariable repetition frequency
	//fmt.Printf("	example_engine:cs1.getperiod(.0):%v\n", cs1.getperiod(.0))

	var piston2 = harmonics(70, nbh)                                                                                                                                                     //.appendSignal(oscillator_noise(tf(.5)))
	var enveloppe2 = []TimedPair{otherPistonHit, variableSilence, tp(.0, .0), tp(.1, .8), tp(.1, .0), variableSilence, otherPistonHit, variableSilence, otherPistonHit, variableSilence} //Envelope of the hit 1
	var ampl2 = timedFloat(.8, customShape(tf(0), tf(1.0), enveloppe2))

	var piston3 = harmonics(65, nbh)                                                                                                                                                     //.appendSignal(oscillator_noise(tf(.4)))
	var enveloppe3 = []TimedPair{otherPistonHit, variableSilence, otherPistonHit, variableSilence, tp(.0, .0), tp(.1, .8), tp(.1, .0), variableSilence, otherPistonHit, variableSilence} //Envelope of the hit 1
	var ampl3 = timedFloat(.8, customShape(tf(0), tf(1.0), enveloppe3))

	var piston4 = harmonics(60, nbh)                                                                                                                                                     //.appendSignal(oscillator_noise(tf(.3)))
	var enveloppe4 = []TimedPair{otherPistonHit, variableSilence, otherPistonHit, variableSilence, otherPistonHit, variableSilence, tp(.0, .0), tp(.1, .8), tp(.1, .0), variableSilence} //Envelope of the hit 1
	var ampl4 = timedFloat(.8, customShape(tf(0), tf(1.0), enveloppe4))

	return signalSum_n("engine").
		appendSignal(piston1.setampl(ampl1)).
		appendSignal(piston2.setampl(ampl2)).
		appendSignal(piston3.setampl(ampl3)).
		appendSignal(piston4.setampl(ampl4))
}

var _ = example_engine() //Avoid not used error
