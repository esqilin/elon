package main

import (
	"github.com/esqilin/godsp"
	"github.com/esqilin/gojack"
	//~ "github.com/esqilin/godsp/waveshape"

	"fmt"
	"os"
	"time"
)

const (
	CLIENT_NAME = "Shnookums"
)

var (
	inPort, outPort *gojack.Port
)

func process(in [][]float32, out_ *[][]float32, arg interface{}) error {
	out := *out_
	for i, buf := range in {
		for j, v := range buf {
			out[i][j] = v
		}
	}
	return nil
}

func jackShutdown(arg interface{}) {
	// test! this function is called by another thread, so standard streams do not work
}

func main() {
	client, err := gojack.NewClient(CLIENT_NAME)
	defer client.Close()
	if nil != err {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	client.SetOptionNoStartServer()

	err = client.Open()
	if nil != err {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if client.IsServerStarted() {
		fmt.Println("Jack server started")
	}
	if client.IsNameNotUnique() {
		fmt.Printf("unique name `%s' assigned\n", client.Name())
	}

	dsp.Init(client.SampleRate())

	client.OnProcess(process, nil)
	client.OnShutdown(jackShutdown, nil)

	fmt.Printf("engine sample rate: %d\n", client.SampleRate())

	inPort, err := client.RegisterAudioIn("input", false)
	if nil != err {
		fmt.Fprintln(os.Stderr, "no more JACK input ports available")
		os.Exit(1)
	}
	outPort, err := client.RegisterAudioOut("output", false)
	if nil != err {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = client.Activate()
	if nil != err {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	outPorts, err := client.SystemOutputPorts(true)
	if nil != err {
		fmt.Fprintln(os.Stderr, err)
	} else {
		client.Connect(outPorts[0], inPort)
	}

	inPorts, err := client.SystemInputPorts(true)
	if nil != err {
		fmt.Fprintln(os.Stderr, err)
	} else {
		client.Connect(outPort, inPorts[0])
		client.Connect(outPort, inPorts[1])
	}

	time.Sleep(10 * time.Second)
}
