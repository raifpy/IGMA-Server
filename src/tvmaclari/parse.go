package tvmaclari

import (
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (t *TvMaclari) parse(r io.Reader) (m Matchs, err error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return m, err
	}
	var s *goquery.Selection
	doc.Find("ul").EachWithBreak(func(i int, s2 *goquery.Selection) bool {
		if i == 2 {
			s = s2
			return false
		}
		return true
	})
	if s == nil {
		return m, t.Error
	}

	s.Find("li").Each(func(_ int, s *goquery.Selection) {

		match := Match{
			Time: s.Find("div.schedule-time").Text(),
		}

		/*r, _ := s.Html()
		fmt.Println("R:")
		fmt.Println(string(r))
		fmt.Scanln()*/

		rawnames := s.Find("div.schedule-detail > div.match-info").Text()
		if rawnames == "" {
			return
		}
		match.MatchName = rawnames
		rawnamesplit := strings.Split(rawnames, " - ")
		if len(rawnamesplit) > 1 {
			match.HomeName = rawnamesplit[0]
			match.AwayName = rawnamesplit[1]
		}

		s.Find("span").Each(func(i int, s *goquery.Selection) {
			switch i {
			case 0:
				{
					match.MatchType = s.Text()
					return
				}
			case 1:
				{

					sp := strings.Split(s.Text(), ",")
					for _, a := range sp {
						if a == "" {
							continue
						}
						match.Channels = append(match.Channels, strings.TrimSpace(a))
					}
				}
			}
		})
		if len(match.Channels) != 0 {
			m = append(m, match)
		}
	})

	return m, err
}
