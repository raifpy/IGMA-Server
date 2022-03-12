package goalchecker

import (
	"log"
	"os"
	"soccerapi/src/goalchecker/waveparser"
	"strconv"
)

type VoiceFrequencyConfig struct {
	StartSecond string
	StopSecond  string
}

func (g *GoalChecker) FindGoalMinuteSecondVoiceFrequency(vc ...VoiceFrequencyConfig) (int, error) {
	var start = "0"
	var stop = "120"

	if len(vc) != 0 {
		if vc[0].StartSecond != "" {
			start = vc[0].StartSecond
		}

		if vc[0].StopSecond != "" {
			stop = vc[0].StopSecond
		}
	}

	out := "/tmp/video2_" + RandStringRunes(10) + "_c"
	path, err := waveparser.VideoToWavFFMPEG(g.o.RawVideoPath, out, start, stop)
	defer os.Remove(path)
	log.Println("\033[31m" + out + "\033[0m")
	if err != nil {

		return -1, err
	}
	out2 := out + ".dat"
	w := waveparser.WaveParser{
		Input:     path,
		Frequency: "1",
		Output:    out2,
		Stdout:    nil,
	}
	defer os.Remove(out2)
	parserR, err := w.Parse()
	if err != nil {
		return -1, err
	}

	biggest := parserR.GetBiggestUpFrequency()
	biggestint, err := strconv.Atoi(biggest.Second)
	if err != nil {
		log.Println("Cloudn't convert string into int:", biggest.Second)
		//w.OnGoalError(err, plain, false, "biggest (frequency) duration convert int")
		biggestint = 0
	}

	biggestint -= 20 // gol'den 20 saniye öncesini alalım
	if biggestint < 0 {
		biggestint = 0 // - olmamalı!
	}

	return biggestint, nil
}
