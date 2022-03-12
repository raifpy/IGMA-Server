package goalchecker

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var OCRMatchRegexNormal = regexp.MustCompile(`(\d+)[.:](\d+)`)

func (g *GoalChecker) TryParseOCR(response string) (normalMatchMinute, normalMatchSecond, extension int, err error) {
	response = strings.TrimSpace(response) // clear
	extension = -1

	if OCRMatchRegexNormal.MatchString(response) {
		rawmatchtime := OCRMatchRegexNormal.FindString(response)
		fmt.Printf("rawmatchtime: %v\n", rawmatchtime)
		if strings.Contains(rawmatchtime, ":") || strings.Contains(rawmatchtime, ".") {
			var splitkey = ":"
			if strings.Contains(rawmatchtime, ".") {
				splitkey = "."
			}
			split := strings.Split(rawmatchtime, splitkey)
			if len(split) == 2 {
				if normalMatchMinute, err = strconv.Atoi(split[0]); err == nil {
					normalMatchSecond, err = strconv.Atoi(split[1])
				}
			}
			//normalMatchTime, err = strconv.Atoi(strings.Split(rawmatchtime, ":")[0])
			return
		}

		log.Println("ELSE :   ", rawmatchtime)

	}

	err = fmt.Errorf("%s and match minute regex not valid", response)
	return

}
