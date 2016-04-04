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
	CLIENT_NAME = "Elon" // jack client name
)

var (
	soundBank []dsp.Sound
)

func main() {
	soundBank = make([]dsp.Sound, 128)

	c, err := jade.New(CLIENT_NAME, false)
	exitOnError(err)
	defer c.Close()

	outName := fmt.Sprintf("%s:output", c.Jack.GetName())

	outCh := c.AddAudioOut("output", 1.0)
	midiCh := c.AddMidiIn("midi_in")

	exitOnError(c.Connect(outName, "system:playback_1"))
	exitOnError(c.Connect(outName, "system:playback_2"))

	c.Play(outCh)

	for {
		midiData := <-midiCh
		status := midiData.Buffer[0]
		index := midiData.Buffer[1]
		velocity := midiData.Buffer[2]

		if 0x90 == status { // NOTE ON
			s := soundBank[index]
			if nil != s && !s.IsFinished() {
				s.Release()
			}

			freq := dsp.StdTuning[index]
			ws, err := wave.NewShapedWave(freq, dsp.Triangle)
			if nil != err {
				fmt.Fprintf(os.Stderr, err.Error())
				continue
			}
			env := dsp.NewEnvelope(0.025, 0.025, 0.8, 0.05)
			s = sound.NewEnvelopeSound(env, sound.NewWaveSound(ws), 0.2*float64(velocity)/128.0)
			soundBank[index] = s

			c.PlaySound(s)
		} else if 0x80 == status { // NOTE OFF
			s := soundBank[index]
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
