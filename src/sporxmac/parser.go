package sporxmac

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/charmap"
)

func (s *SporxMax) parse(reader io.Reader) (m Matchs, err error) {

	//"windows-1254"
	outr := charmap.Windows1254.NewDecoder().Reader(reader)
	/*outr, err := iconv.NewReader(reader, "windows-1254", "utf-8")
	if err != nil {
		return m, err
	}*/

	q, err := goquery.NewDocumentFromReader(outr)
	if err != nil {
		return nil, err
	}

	q.Find("li").Each(func(_ int, s *goquery.Selection) {

		if s.AttrOr("class", "") != "odd" && s.AttrOr("class", "") != "even" {
			return
		}
		//var match Match

		time := s.Find(".ch-time").Text()
		channel := s.Find(".ch-desc .ch-name").Text()

		mactipraw := s.Find(".ch-desc .ch-text").Text()

		//macraw := s.Find(".ch-desc .ch-text .ch-link").Text()
		//fmt.Printf("macraw: %v\n", (macraw))
		/*fmt.Printf("time: %v\n", (time))
		fmt.Printf("channel: %v\n", (channel))
		fmt.Printf("mactipraw: %v\n", (mactipraw))*/

		i := strings.Index(mactipraw, "(")
		if i == -1 {
			return
		}
		macraw := strings.TrimRight(strings.TrimLeft(mactipraw[:i], " "), " ")

		if macraw == "" {
			return
		}

		macsplit := strings.Split(macraw, " - ")
		if len(macsplit) == 1 {
			return
		}

		/*mactipraw =
		strings.TrimFunc(mactipraw[len(macraw):], func(r rune) bool {
			return r == '(' || r == ')'
		})*/
		i2 := strings.Index(mactipraw, ")")
		if i2 == -1 {
			return
		}

		mactipraw = mactipraw[i+1 : i2]
		mactipraw = strings.TrimLeft(mactipraw, " ")
		mactipraw = strings.TrimRight(mactipraw, " ")
		mactipraw = strings.ToLower(mactipraw)

		m = append(m, Match{
			MatchType: mactipraw,
			HomeName:  macsplit[0],
			AwayName:  macsplit[1],
			MatchName: macraw,
			Time:      time, Channels: []string{
				strings.Title(strings.ToLower(channel)),
			},
		})

	})

	return

}

func (s *SporxMax) Test() {
	res, err := s.request()
	if err != nil {
		log.Fatalln(err)
	}

	ü, err := s.parse(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("\n ")

	r, _ := json.MarshalIndent(ü, "", " ")
	fmt.Println(string(r))

}
