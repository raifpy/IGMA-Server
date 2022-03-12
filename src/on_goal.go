package soccer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path"
	dsystem "soccerapi/src/downloadsystem"
	"soccerapi/src/goalchecker"
	"soccerapi/src/goalchecker/ocr"
	"soccerapi/src/iptv"
	"soccerapi/src/scores"
	"soccerapi/src/worker/types"
	"strings"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

/*
	Eksikler
		Düşük kalite algoritma
		Tekrar eden kod  satırları
		Sabit belleği doldurabilecek io açıkları
		Yetersiz hata ayrıştırma

*/
func (w *Soccer) OnGoal(s scores.GameScore, schan chan scores.GameScore) {
	log.Printf("%s %v %v %s \033[31mGOL!\033[0m\n", s.GameGame.HomeCompetitor.Name, s.GameGame.HomeCompetitor.Score, s.GameGame.AwayCompetitor.Score, s.GameGame.AwayCompetitor.Name)
	var status string = "goal"
	plain := s.ToPlain()

	go func() {
		ss := <-schan

		ssplain := ss.ToPlain()

		if plain.HomeScore > ssplain.HomeScore || plain.AwayScore > ssplain.AwayScore {
			log.Printf("Maç (%s - %s) \033[31mgol iptal\033[0m, ya da algoritma yanlış saydı, her neyse status hot olarak çeviriliyor", plain.HomeName, plain.AwayName)
			status = "hot"

			//! Hatalı ya da eksik çalışıyor
		}

		fmt.Printf("MAÇ %v - %v\n", plain.HomeName, plain.AwayName)

		fmt.Printf("plain.HomeScore: %v\n", plain.HomeScore)
		fmt.Printf("ssplain.HomeScore: %v\n", ssplain.HomeScore)

		fmt.Printf("plain.AwayScore: %v\n", plain.AwayScore)
		fmt.Printf("ssplain.AwayScore: %v\n", ssplain.AwayScore)

		//if maçdbyesetedilecek {

		//}
	}()

	clips, err := w.GetMatchClips(int64(plain.Id))
	if err != nil {
		if err := w.SetMatchClips(MatchClips{ //?? Update status
			MatchID: int64(s.GameGame.ID),
			Clips: []MatchClip{
				{
					Client: ClientMatchClip{
						HomeScore:  plain.HomeScore,
						AwayScore:  plain.AwayScore,
						Update:     "watching",
						UpdateID:   s.LastUpdateID,
						CheckAfter: 10,
						MatchID:    int64(s.GameGame.ID),
						Status:     status,
						Minute:     plain.Minute,
					},
				},
			},
		}); err != nil {
			w.OnGoalError(err, plain, true, "database setmatchclip error")
			return
		}
	}

	clips.Clips = append(clips.Clips, MatchClip{
		Client: ClientMatchClip{
			HomeScore:  plain.HomeScore,
			AwayScore:  plain.AwayScore,
			Update:     "watching",
			UpdateID:   s.LastUpdateID,
			CheckAfter: 10,
			MatchID:    int64(s.GameGame.ID),
			Status:     status,
			Minute:     plain.Minute,
		},
	})

	w.ReplaceMatchClips(clips)

	channel := w.GameScoreWatcher.GameIDs[s.GameGame.ID]
	if channel == "" {
		w.OnGoalError(errors.New("channel not found"), plain, true, "maça uygun kanal bulunamadı")
		return
	}
	fmt.Printf("kanal: %s\n", channel)

	list := w.IpTv.FilterByChannel(channel)
	if len(list) == 0 {
		w.OnGoalError(fmt.Errorf("%s channel not found in %s", channel, w.config.IpTvConfigPath), plain, true, "maça uygun iptv bulunamadı")
		return
	}

	//maçdbyesetedilecek = true

	log.Printf("%d adet iptv uyumlu, izlenecek\n", len(list))
	wait := sync.WaitGroup{}
	var responselist = []IpTvResponse{}
	for _, l := range list {
		wait.Add(1)

		go func(l *iptv.IpTv) {
			defer wait.Done()
			fmt.Printf("l.Name: %v\n", l.Name)

			var r = IpTvResponse{
				WatchTime: time.Now(),
				IpTvName:  l.Name,
			}
			r.MediaPath = path.Join(createTempVideoDir(""), fmt.Sprintf("%d_%d_%s_%d_%s_raw", s.GameGame.ID, s.LastUpdateID, (s.GameGame.HomeCompetitor.Name+"_"+s.GameGame.AwayCompetitor.Name), s.GameGame.GameTimeAndStatusDisplayType, l.Name))
			file, err := os.Create(r.MediaPath)
			if err != nil {
				r.FatalError = err
			} else {
				defer file.Close()
				var counter int
			w:
				r.Size, r.FatalError = l.Watch(io.MultiWriter(file, progressbar.Default(-1, r.MediaPath)), l.Channels[channel], func(e error) {
					log.Println("IPTV Watch OnError: ", e)
					r.ErrorList = append(r.ErrorList, e)
				})
				if r.Size == 0 {
					counter++
					if counter <= 1 {
						goto w
					}
				}
				log.Println("Maç izlendi.")
				r.Saved = r.FatalError == nil
				responselist = append(responselist, r)
			}

		}(l)
	}

	log.Println("maçlar bekleniyor..")
	wait.Wait()
	fmt.Println("STATUS: ", status)

	if len(responselist) == 0 {
		fmt.Println("Eyvah! responselist boş")
		//maçdbyesetedilecek = false
		w.OnGoalError(errors.New("responselist empty"), plain, true, fmt.Sprintf("channel %s eşleşen ipadresleri: %v ancak hiç biri çalışmamış gözüküyor.", channel, list))
		return
	}

	ciddiyeal, err := w.FilterIpTvResponseList(responselist)
	if err != nil {
		w.OnGoalError(err, plain, true, fmt.Sprintf("%d yayın içerisinden düzgün yayın bulunamadı", len(responselist)))
		return
	}
	//! DÜZELECEK!
	/*if len(w.Worker.WorkerMap.Map) == 0 {
		fmt.Println("Eyvah! worker yok ğ")
		//maçdbyesetedilecek = false
		w.OnGoalError(errors.New("workers not exists/enought"), plain, true, "video render etmek için worker gerekli, ve bizde hiç yok")
		return
	}*/

	cr, _ := json.MarshalIndent(ciddiyeal, "", " ")
	fmt.Printf("ciddiyeal:  %v\n", string(cr))

	var ocrengine = "1"
	/*if regexp.MustCompile("[bB][eE][iIİ][nN]").MatchString(channel) { // Kanal bein ise engine 1'i kullanıyorum.
		ocrengine = "1"
	}*/

	var freqstart, freqstop = "0", "150" //120
	var biggestint int

	goalcheck, err := goalchecker.NewGoalChecker(goalchecker.Options{
		OcrConfig: ocr.Config{
			OCREngine: ocrengine,
		},
		RawVideoPath: ciddiyeal.MediaPath, GoalMinute: strings.Replace(plain.Minute, "'", "", -1)})

	if err != nil {
		fmt.Printf("err: %v\n", err)
		w.OnGoalError(err, plain, true, "GoalChecker error")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	ocrsec, err := goalcheck.FindGoalSecondOCR(ctx)
	cancel()
	if err != nil {
		log.Println("goalcheck.FindGoalSecondOCR \033[31mERROR\033[0m: ", err)
		w.OnGoalError(err, plain, false, "clip ocr score match goal clip find err")

	} else {
		if ocrsec < 0 {
			ocrsec = 0
		}
		biggestint = ocrsec
		freqstart = fmt.Sprint(ocrsec)
		freqstop = fmt.Sprint(ocrsec + 60)
	}

	fmt.Printf("biggestint: %v\n", biggestint)
	fmt.Printf("freqstart: %v\n", freqstart)
	fmt.Printf("freqstop: %v\n", freqstop)

	if w.useVoiceFrequencyToo(plain) {
		log.Println("Voice Frequency de dikkate alınacak!")

		vqsec, err := goalcheck.FindGoalMinuteSecondVoiceFrequency(goalchecker.VoiceFrequencyConfig{
			StartSecond: freqstart,
			StopSecond:  freqstop,
		})
		if err != nil {
			log.Println("FindGoalMinuteSecondVoiceFrequency \033[31mERROR\033[0m: ", err)
			w.OnGoalError(err, plain, false, "voice frquency match goal clip find err")
		} else {
			fmt.Printf("vqsec: %v\n", vqsec)
			biggestint += vqsec
		}

	}

	fmt.Printf("biggestint: %v\n", biggestint)

	var kalan int = (biggestint % 60)
	var bolum int = (biggestint / 60)

	ffmpegstart := fmt.Sprintf("00:%s:%s", IntToFfmpegInt(bolum), IntToFfmpegInt(kalan))
	ffmpegstop := fmt.Sprintf("00:%s:%s", IntToFfmpegInt((biggestint+59)/60), IntToFfmpegInt((biggestint+59)%60))

	fmt.Printf("ffmpegstart: %v\n", ffmpegstart)
	fmt.Printf("ffmpegstop: %v\n", ffmpegstop)

	fmt.Printf(" ROW ciddiyeal.MediaPath: %v\n", ciddiyeal.MediaPath)

	clipid := rand.Int63n(999999999998)

	time.AfterFunc(time.Hour*2, func() {
		log.Printf("RAW MEDIA \033[31mREMOVING\033[0m: %v", ciddiyeal.MediaPath)
		os.Remove(ciddiyeal.MediaPath)
	})

	var counter int
bastan: //? goto kullanımı C/C++ geliştiricileri tarafından fazla tercih edilmez
	counter++
	if counter >= 3 {
		w.OnGoalError(errors.New("counter error"), plain, true, "counter 3'ü aştı!")
		return
	}

	worker := w.Worker.WorkerMap.Random()
	if worker == nil {
		log.Println("Worker is nil!")
		w.OnGoalError(errors.New("nil worker"), plain, true, "worker is nil")
		return
	}

	log.Println("Worker: ", worker.Id)
	did := RandStringRunes(30)
	dtoken := RandStringRunes(30)
	w.Dsystem.Set(did, dsystem.DsystemValue{
		Token: dtoken,
		//Ip:     worker.IP //!! Eksik
		Path:   ciddiyeal.MediaPath,
		Delete: time.Hour,
	})

	job := types.Job{
		Status: "newjob",

		Expired: time.Now().Add(time.Minute * 15),
		Exec: &types.Exec{
			Exec: "ffmpeg",
			Args: []string{
				"-y",
				"-i", "$mediain:downloadid=" + did + "&token=" + dtoken,
				//"-filter:v", "crop=in_w:in_h-75:0:out_h",
				"-ss", ffmpegstart, "-to", ffmpegstop,
				fmt.Sprintf("$mediaout:uploadid={jobid}&filename=%s.mp4", "raw_goal_stream"),
			},
			ShareSTD: false,
		},
	}
	/*if err := worker.Conn.WriteJSON(job); err != nil {
		log.Println("worker job kabul etmiyor it: ", err)
		w.OnGoalError(err, plain, false, "worker job write error")
		goto bastan
	}*/
	//client job request
	jobid, err := w.Worker.AddJobDb(context.Background(), worker, job)
	if err != nil {
		log.Println("Err on AddJobDb", err)
		w.OnGoalError(err, plain, false, "worker job write (addjobdb) error")
		goto bastan
	}

	log.Println("worker'a job set edildi!")
	workerchan := make(chan types.WebsocketContact) // client job response
	w.Worker.WorkerMap.SetChan(jobid, workerchan)

	defer func() {
		w.Worker.WorkerMap.DelChan(jobid)
		close(workerchan)
	}()

	var res types.WebsocketContact
yanit:
	select {
	case <-time.NewTicker(time.Minute * 15).C:
		log.Println("worker response timeout")
		goto bastan
	case res = <-workerchan:
		break

	}
	log.Println("\033[32mVERI GELDI\033[0m")
	_jsonobject, _ := json.MarshalIndent(res, "", " ")
	fmt.Println(string(_jsonobject))

	if res.Type != "error" && res.Type != "done" {
		log.Println("Worker status update: ", res.Type)

		go func() { // çorba oldu amk
			response, err := w.GetMatchClips(int64(s.GameGame.ID))
			if err != nil { // yoksa baştan set edelim?
				w.OnGoalError(err, plain, true, "GetMatchClips maçı update edecektim!")
				return
			}
			response.MatchID = int64(plain.Id)

			clips := []MatchClip{}
			for _, v := range response.Clips {
				if v.Client.UpdateID == plain.UpdateID {
					v = MatchClip{
						Client: ClientMatchClip{
							HomeScore:  plain.HomeScore,
							AwayScore:  plain.AwayScore,
							Update:     res.Type,
							MatchID:    int64(plain.Id),
							UpdateID:   plain.UpdateID,
							Status:     status,
							Minute:     plain.Minute,
							ClipID:     clipid,
							CheckAfter: 10,
						},

						WorkerID: worker.Id,
						JobID:    jobid,
					}
				}
				clips = append(clips, v)
			}
			response.Clips = clips
			w.ReplaceMatchClips(response)
		}()

		goto yanit
	}

	if res.Error != nil {
		log.Println("Worker Error: ", res.Error)
		goto bastan //?
	}

	if res.Update != nil {
		switch res.Update.Job.Status {
		case "done":

			t, err := w.Worker.GetJobDb(context.Background(), jobid)
			if err != nil {
				w.OnGoalError(err, plain, true, fmt.Sprintf("GetJobDb: %d", jobid))
				return

			}
			if t.JobResponseStore == nil {
				w.OnGoalError(errors.New("JobResponseStore is nil"), plain, true, "Unexcepted worker response")
				return
			}

			response, err := w.GetMatchClips(int64(s.GameGame.ID))
			if err != nil { // yoksa baştan set edelim?
				w.OnGoalError(err, plain, true, "GetMatchClips maçı update edecektim!")
				return
			}
			response.MatchID = int64(plain.Id)

			clips := []MatchClip{}
			for _, v := range response.Clips {
				if v.Client.UpdateID == plain.UpdateID {
					v = MatchClip{
						Client: ClientMatchClip{
							HomeScore:  plain.HomeScore,
							AwayScore:  plain.AwayScore,
							Update:     "done",
							MatchID:    int64(plain.Id),
							UpdateID:   plain.UpdateID,
							Status:     status,
							Minute:     plain.Minute,
							ClipID:     clipid,
							CheckAfter: 0,
						},
						Path:     t.JobResponseStore.LocalPath,
						GdriveID: t.JobResponseStore.GdriveId,
						WorkerID: worker.Id,
						JobID:    jobid,
					}
				}
				clips = append(clips, v)
			}
			response.Clips = clips
			if err := w.ReplaceMatchClips(response); err != nil {
				w.OnGoalError(err, plain, true, "ReplaceMatchClips")
				return
			}

		case "uploading":
			log.Println("State 'Video upload'")

		}
	}

}
func (s *Soccer) useVoiceFrequencyToo(p scores.PlainGameScore) bool {
	return true
}

type IpTvResponse struct {
	Duration     string
	GoalDuration string

	WatchTime  time.Time
	IpTvName   string
	FatalError error
	ErrorList  []error
	Size       int64

	Saved     bool
	MediaPath string
}
