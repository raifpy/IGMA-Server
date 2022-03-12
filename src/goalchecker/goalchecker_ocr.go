package goalchecker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

//GoalSecond, error
func (g *GoalChecker) FindGoalSecondOCR(ctx context.Context /* trycounter int*/) (int, error) {
	durationstring, err := GetVideoDuration(g.o.RawVideoPath)
	if err != nil {
		return -1, err
	}
	fmt.Printf("durationstring: %v\n", durationstring)
	fduration, err := strconv.ParseFloat(durationstring, 32)
	if err != nil {
		return -1, err
	}
	fmt.Printf("fduration: %v\n", fduration)

	out := "frame_" + RandStringRunes(10) + ".png"
	defer os.Remove(out)

	/*var counterframe int = 0

	lastframe:*/

	lastsec, err := TryGetLastFrameFromVideo(int(fduration), g.o.RawVideoPath, out)
	if err != nil {
		return -1, err
	}
	fmt.Printf("lastsec: %v\n", lastsec)
	file, err := os.Open(out)
	if err != nil {
		return -1, err
	}

	ocrtext, err := g.ocr.Text(ctx, file)
	file.Close()
	if err != nil {
		return -1, err
	}
	fmt.Printf("ocrtext: %v\n", ocrtext)

	//fmt.Printf("g.o.GoalMinute: %v\n", g.o.GoalMinute)
	mm, ms, em, err := g.TryParseOCR(ocrtext)
	if err != nil {
		return -1, err
	}
	fmt.Printf("mm: %v\n", mm)
	fmt.Printf("ms: %v\n", ms)
	fmt.Printf("em: %v\n", em)

	fmt.Printf("g.o.GoalMinute: %v\n", g.o.GoalMinute)

	if g.o.GoalMinute[len(g.o.GoalMinute)-1] == '\'' { //!! Panic
		g.o.GoalMinute = g.o.GoalMinute[0 : len(g.o.GoalMinute)-1]
	}
	fmt.Printf("g.o.GoalMinute: %v\n", g.o.GoalMinute)

	if strings.Contains(g.o.GoalMinute, "+") {
		fmt.Println("Extension'da ciddiye alınmalı")

		//panic("extention henüz codlanmadı")
	}

	apigoalsec, err := strconv.Atoi(g.o.GoalMinute)
	if err != nil {
		log.Println("g.o.GoalMinute int'e evirilemedi")
		return -1, err
	}
	apigoalsec--
	if apigoalsec > mm {
		return -1, errors.New("GoalMinute bigger than match last frame's minute")
	}

	fark := mm - apigoalsec
	fmt.Printf("fark: %v\n", fark)
	/*if fark == 1 {
		fark = 0
	}*/

	fark = fark * 60

	fmt.Printf("fark2: %v\n", fark)

	return lastsec - (fark + ms), nil
}
