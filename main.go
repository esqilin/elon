package main

import (
	"github.com/esqilin/godsp"
	"github.com/esqilin/godsp/sound"
	"github.com/esqilin/godsp/wave"
	"github.com/esqilin/gojade"

	"fmt"
	"os"
)

const (
	CLIENT_NAME = "Shnookums" // jack client name
)

var (
	soundBank []dsp.Sound
)

func main() {
	soundBank = make([]dsp.Sound, 128)

	c, err := gojade.New(CLIENT_NAME, false)
	defer c.Close()
	exitOnError(err)

	//~ lfoObj, err := wave.NewShapedWave(1, dsp.Triangle)
	//~ waveShape, err := waveshape.NewShapeSwitcher(dsp.Sine, dsp.Triangle, lfoObj)
	//~ pulseShape := waveshape.NewPulse(0.0)
	//~ //moog := waveshape.NewMoog()
	//~ whitenoise := waveshape.NewWhiteNoise()
	//~ waveObj, err := wave.NewShapedWave(440, whitenoise)
	//~ if nil != err {
        //~ fmt.Fprintf(os.Stderr, "error setting up first wave: %s", err)
        //~ os.Exit(1)
	//~ }

	//~ procDat := processData{ waveObj, lfoObj, pulseShape }

    outCh, err := c.AddAudioOut("output", true)
	exitOnError(err)
	//~ exitOnError(c.AddAudioIn("input", false))

	exitOnError(c.ConnectSystemSpeaker("output", 0))
	exitOnError(c.ConnectSystemSpeaker("output", 1))
	//~ exitOnError(c.ConnectSystemAudioSource("input", 0))

    midiCh, err := c.AddMidiIn("midi_in")
	exitOnError(err)

	//~ exitOnError(c.AddMidiOnCallback("midi_in", func(index, velocity byte) {
		//~ s := soundBank[index]
		//~ if nil != s && !s.IsFinished() {
			//~ s.Release()
		//~ }
//~
		//~ freq := dsp.StdTuning[index]
		//~ ws, err := wave.NewShapedWave(freq, dsp.Triangle)
		//~ if nil != err {
			//~ fmt.Fprintf(os.Stderr, err.Error())
			//~ return
		//~ }
		//~ env := dsp.NewEnvelope(0.1, 0.3, 0.8, 0.15)
		//~ s = sound.NewEnvelopeSound(env, sound.NewWaveSound(ws), float64(velocity)/128.0)
		//~ soundBank[index] = s
//~
		//~ c.PlaySound(s)
	//~ }))

	//~ exitOnError(c.AddMidiOffCallback("midi_in", func(index, _ byte) {
		//~ s := soundBank[index]
		//~ if nil == s {
			//~ fmt.Fprintf(
				//~ os.Stderr,
				//~ "released key with no according unreleased sound?\n",
			//~ )
			//~ return
		//~ }
//~
		//~ s.Release()
	//~ }))

	//~ c.PlaySound(wave.NewBiquadWave(
	//~ [...]float64{ 0.0009405039485878813, 0.0018810078971757626, 0.0009405039485878813 },
	//~ [...]float64{ -1.91139741777364, 0.9151594335679912 },
	//~ ))

	//~ select {} // wait

    c.Play(outCh)

    for {
        event := <-midiCh
        if event.IsOn() {
            s := soundBank[event.Index]
            if nil != s && !s.IsFinished() {
                s.Release()
            }

            freq := dsp.StdTuning[event.Index]
            ws, err := wave.NewShapedWave(freq, dsp.Triangle)
            if nil != err {
                fmt.Fprintf(os.Stderr, err.Error())
                continue
            }
            env := dsp.NewEnvelope(0.1, 0.3, 0.8, 0.15)
            s = sound.NewEnvelopeSound(env, sound.NewWaveSound(ws), float64(event.Velocity)/128.0)
            soundBank[event.Index] = s

            c.PlaySound(s)
        } else if event.IsOff() {
            s := soundBank[event.Index]
            if nil == s {
                fmt.Fprintf(
                    os.Stderr,
                    "released key with no according unreleased sound?\n",
                )
            } else {
                s.Release()
            }
        }
    }
}

func exitOnError(err error) {
	if nil == err {
		return
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
