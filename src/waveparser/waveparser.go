package waveparser

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type WaveParser struct {
	Input     string // in.wav
	Frequency string // 1
	Output    string // out.dat?

	Stdout io.Writer
}

func (w WaveParser) exec() error {
	if w.Stdout == nil {
		w.Stdout = os.Stdout
	}
	cmd := exec.Command("sox", w.Input, "-r", w.Frequency, w.Output)
	cmd.Stdout = w.Stdout
	cmd.Stderr = w.Stdout

	return cmd.Run()

}

func (w WaveParser) parse() (list FrameInfos, err error) {
	out, err := os.ReadFile(w.Output)
	if err != nil {
		return nil, err
	}
	out = bytes.ReplaceAll(out, []byte("\r"), []byte(""))

	for _, value := range bytes.Split(out, []byte("\n")) {

		if len(value) == 0 || value[0] == ';' {
			continue
		}

		s := strings.Split(string(value), " ")

		var ns = []string{}

		for _, a := range s {
			if a != "" {
				ns = append(ns, a)
			}
		}
		if len(ns) == 3 {
			var i FrameInfo
			i.Second = ns[0]
			i.UpWave, _ = strconv.ParseFloat(ns[1], 64)
			i.DownWave, _ = strconv.ParseFloat(ns[2], 64)
			list = append(list, i)
		}

	}

	return

}

type FrameInfo struct {
	Second   string
	UpWave   float64
	DownWave float64

	Index int
}

func (w WaveParser) Parse() (FrameInfos, error) {
	if err := w.exec(); err != nil {
		return nil, err
	}
	return w.parse()
}

func (w WaveParser) Remove() {
	os.Remove(w.Input)
	os.Remove(w.Output)
}
