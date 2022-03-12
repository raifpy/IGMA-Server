package goalchecker

import (
	"soccerapi/src/goalchecker/ocr"
)

type Options struct {
	RawVideoPath string
	GoalMinute   string
	OcrConfig    ocr.Config
}

type GoalChecker struct {
	o   Options
	ocr ocr.OCR
}

func NewGoalChecker(o Options) (g *GoalChecker, err error) {
	g = &GoalChecker{
		o: o,
		//ocr: ocr.NewOcrSpaceOcr(o.OcrConfig),
	}

	g.ocr, err = ocr.NewEasyOcr(ocr.Options{
		Path: "easyocrbin", // /usr/bin/easyocrbin
	})
	return
}
