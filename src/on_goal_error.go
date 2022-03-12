package soccer

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"soccerapi/src/scores"
	"strings"
	"time"

	"github.com/mitchellh/colorstring"
)

func (s *Soccer) OnGoalErrorFn(err error, score scores.PlainGameScore, fatal bool, ek ...string) {
	log.Println(colorstring.Color(fmt.Sprintf("Maç [blue]%s [red]GoalError[white]: [red]%v[white] [red]FATAL[white]: [green]%v %v[white]", score.HomeName+" - "+score.AwayName, err, fatal, ek)))

	go func() { //? Kişiselleştirilebilir
		if res, _ := http.Get("https://api.telegram.org/bot<BOTTOKEN>/sendMessage?chat_id=<TELEGRAM_USER_ID>&text=" + url.QueryEscape(fmt.Sprintf("Fatal %v\n\nError: %v\n\nMaç: %s - %s\nDakika: %s\n\nEk: %s", fatal, err, score.HomeName, score.AwayName, score.Minute, strings.Join(ek, " ")))); res != nil {
			res.Body.Close()
		}

	}()

	if fatal {

		if r, err := os.ReadDir(createTempVideoDir("")); err == nil {
			for _, v := range r {
				if v.IsDir() {
					continue
				}

				if strings.HasPrefix(v.Name(), fmt.Sprintf("%d_%d", score.Id, score.UpdateID)) {
					log.Println("ölü dosya bulundu; siliniyor!") // ? ?
					os.Remove(path.Join(createTempVideoDir(""), v.Name()))
				}
			}
		}

		var bidaha = false
	bida: //

		response, err := s.GetMatchClips(int64(score.Id))
		if err != nil {
			fmt.Printf("err: %v\n", err)
			if !bidaha {
				bidaha = true
				time.Sleep(time.Second * 4)
				goto bida
			}
			return
		}
		var newclips = []MatchClip{}
		for _, v := range response.Clips {
			if v.Client.UpdateID == score.UpdateID {
				v.Client.Update = "error"
			}

			if v.Client.UpdateID == 0 {
				continue
			}

			newclips = append(newclips, v)
		}

		response.Clips = newclips
		if err := s.ReplaceMatchClips(response); err != nil {
			log.Println("OnGoalError ReplaceMatchClips: ", err)
		}

	}

}
