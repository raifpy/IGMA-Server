package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	soccer "soccerapi/src"
	"soccerapi/src/goalchecker"
	"soccerapi/src/goalchecker/ocr"
	"soccerapi/src/iptv"
	"soccerapi/src/scores"
	"soccerapi/src/sporxmac"
	"soccerapi/src/tvmaclari"
	"strconv"
	"strings"
	"time"

	progressbar "github.com/schollz/progressbar/v3"
)

func test() {
	switch os.Args[1] {
	case "iptv":
		iptvtest()

	case "score":
		scoretest()
	case "tvmaclari":
		tvmaclaritest()
	case "tvmaclariall":
		tvmaclariall()

	case "sporxmacall":
		sporxmacalltest()

	case "sporxmac":
		sporxmactest()

	case "customscore":
		customscore()
	case "gamescore":
		gamescoretest()
	case "videogoal":
		videogoaltest()
	case "mongouseradd":
		mongouseradd()

	case "ocrspace":
		ocrspacetest()

	case "goalchecker":
		goalcheckertest()

	case "videoduration":
		videodurationtest()

	case "lastframe":
		lastframetest()

	default:
		log.Fatalf("\033[31m%s\033[0m test eşleşmedi.\n", os.Args[1])
	}

	os.Exit(0)
}

func tvmaclaritest() {

	res, err := tvmaclari.Get()
	if err != nil {
		log.Fatalln(err)
	}
	a := res.FilterWithSportName("futbol").FilterWithChannelsName(MatchChannelList)
	ires, err := json.MarshalIndent(a, "", " ")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(ires))
}

func tvmaclariall() {
	res, err := tvmaclari.Get()
	if err != nil {
		log.Fatalln(err)
	}
	ires := res.FilterWithSportName("futbol")
	b, _ := json.MarshalIndent(ires, "", " ")
	fmt.Println(string(b))
}

func iptvtest() {

	ip, err := iptv.NewIpTvStand(iptv.Options{
		IpTvConfigPath: iptvconfigpath,
		OnUpdate: func() {
			log.Println("config değişti")
		},
		OnError: func(e error) { fmt.Printf("e: %v\n", e) },
		HttpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	})

	if err != nil {
		log.Fatalln(err)
	}
	for index, l := range ip.List {
		log.Printf("\033[32m%d\033[0m \033[34mname:\033[0m %s\n", index, l.Name)
	}
	var index int
	var indexs string
	fmt.Scanln(&indexs)

	index, err = strconv.Atoi(indexs)
	if err != nil {
		log.Fatalln(err)
	}

	for channel := range ip.List[index].Channels {
		fmt.Printf("\033[35mKanal:\033[0m %s\n", channel)
	}

	indexs, err = bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Fatalln(err)
	}
	indexs = strings.TrimRight(indexs, "\n")

	c := ip.List[index].Channels[indexs]
	if c == "" {
		log.Fatalf("\033[37m|%s|\033[0m için path bulunamadı\n", indexs)
	}

	file, err := os.Create("video.mp4")
	if err != nil {
		log.Fatalln(err)
	}

	size, err := ip.List[index].Watch(io.MultiWriter(file, progressbar.Default(-1, fmt.Sprintf(""))), c, func(e error) {
		fmt.Println("Dostum watch error: ", err)
	})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("size: %v\n", size)

}

func scoretest() {

	s := scores.NewScores()
	s.OnWhenGoal = func(gs scores.GameScore, c chan scores.GameScore) {
		fmt.Printf("\033[31mgol olmuş dayı\033[0m%s %v - %v %s\n", gs.GameGame.HomeCompetitor.Name, gs.GameGame.HomeCompetitor.Score, gs.GameGame.AwayCompetitor.Score, gs.GameGame.AwayCompetitor.Name)
		<-c
	}

	res, err := s.RequestAllScores(true)
	if err != nil {
		log.Fatalln(err)
	}
	//fmt.Printf("res.ToPrettyJson(): %v\n", string(res.ToPrettyJson()))

	for _, a := range res.Games {
		s.GameIds = append(s.GameIds, strconv.Itoa(a.ID))
	}

	s.Loop()
	<-make(chan bool)

}

func customscore() {
	id := os.Args[2] //3477004
	gs, err := scores.NewScores().RequestGameScore(id)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Print(string(gs.ToPrettyJson()))

}

func gamescoretest() {
	now := time.Now()
	s, err := soccer.NewGameScoreWatcher(soccer.GameScoreWatcherOptions{
		IpTvChannelList: MatchChannelList,
		OnGoal: func(gs scores.GameScore, c chan scores.GameScore) {
			//var field string = "gol"
			fmt.Printf("Lig: %v\n", gs.GameGame.CompetitionDisplayName)
			fmt.Printf("Skor: %v - %v", gs.GameGame.HomeCompetitor.Score, gs.GameGame.AwayCompetitor.Score)
			fmt.Println()
			dogrulama := <-c
			fmt.Println("Doğrulama: ", dogrulama.GameGame.HomeCompetitor.Score, dogrulama.GameGame.AwayCompetitor.Score)
		},
		OnRequestError: func(err error) {
			log.Println("Scores request error: ", err)
		},
	})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("\033[32m:", math.Abs(time.Until(now).Seconds()), "\033[0m saniyede sonuca ulaşıldı")

	for key, value := range s.GameIDs {
		fmt.Println(key, value)
		r, _ := s.Scores.RequestGameScore(fmt.Sprint(key))
		fmt.Println(r.GameGame.HomeCompetitor.Name, " : ", r.GameGame.AwayCompetitor.Name, "\n ")
	}

	s.Scores.Loop()

	<-make(chan bool)

}

func videogoaltest() {
	videopath := os.Args[2]
	check, err := goalchecker.NewGoalChecker(goalchecker.Options{
		RawVideoPath: videopath,
	})
	if err != nil {
		log.Fatalln(err)
	}
	sec, err := check.FindGoalMinuteSecondVoiceFrequency()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("En yüksek saniye: ", sec)
}

func mongouseradd() {
	s, err := soccer.NewSoccer(soccer.Config{
		FiberHost:           Host,
		MongoAddr:           MongoAddr,
		WorkerWsPath:        "/wsworker",
		SupportedTvChannels: MatchChannelList,
		IpTvConfigPath:      "iptvconf.json",
	})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Print("Name: ")
	var name string
	fmt.Scanln(&name)

	token := soccer.RandStringRunes(40)

	err = s.SetUser(context.Background(), soccer.UserClient{
		Name:         name,
		Expired:      time.Now().Add(time.Hour * 24 * 30 /* * 12*/), // 1 ay
		RegisterTime: time.Now(),
		Token:        token,
	})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("token: %v\n", token)
}

func goaltest(s *soccer.Soccer) {
	reader := bufio.NewReader(os.Stdin)
	var err error
	for index, value := range s.GameScoreWatcher.ChannelList {
		fmt.Println(index, value)
	}

	fmt.Print("\nKanal seçin: ")
	kanal, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalln(err)
	}
	kanal = strings.TrimRight(kanal, "\n")
	kanali, err := strconv.Atoi(kanal)
	if err != nil {
		log.Fatalln(err)
	}
	kanal = s.GameScoreWatcher.ChannelList[kanali]

	matchid := rand.Intn(9999)
	updateid := int64(rand.Intn(999999))

	s.GameScoreWatcher.GameIDs[matchid] = kanal

	veri := scores.GameScore{
		LastUpdateID: updateid,

		//RequestedUpdateID: 1111,
		GameGame: scores.GameGame{

			ID:              matchid,
			GameTimeDisplay: fmt.Sprintf("%d'", rand.Intn(90)),
			StatusText:      "live",
			JustEnded:       false,
			HomeCompetitor: scores.HomeCompetitor{
				Score: 3,
				Name:  "Takım1",
			},
			AwayCompetitor: scores.AwayCompetitor{
				Score: 0,
				Name:  "Takım2",
			},
		},
	}

	schan := make(chan scores.GameScore)
	go func() {
		sleep := 2
		log.Printf("scan'a veri gönderilecek. %d saniye bekleniyor!\n", sleep)
		time.Sleep(time.Second * time.Duration(sleep))
		schan <- veri
	}()

	s.OnGoal(veri, schan)

}

func sporxmacalltest() {
	res, err := sporxmac.Get()
	if err != nil {
		log.Fatalln(err)
	}

	r, _ := json.MarshalIndent(res.FilterWithSportName("futbol"), "", " ")
	fmt.Println(string(r))
}

func sporxmactest() {
	res, err := sporxmac.Get()
	if err != nil {
		log.Fatalln(err)
	}

	r, _ := json.MarshalIndent(res.FilterWithSportName("futbol").FilterWithChannelsName(MatchChannelList), "", " ")
	fmt.Println(string(r))
}

func ocrspacetest() {
	path := os.Args[2]
	//fmt.Printf("path: %v\n", path)

	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)

	}

	c := ocr.NewOcrSpaceOcr()
	response, err := c.OCR(context.Background(), file, c.Body, "image.png")
	if err != nil {
		log.Fatalln(err)
	}

	p, _ := json.MarshalIndent(response, "", " ")
	fmt.Println(string(p))
}

func goalcheckertest() {

	path := os.Args[2]

	if len(os.Args) == 4 && os.Args[3] == "voice" {
		log.Println("\033[31mUsing voice!\033[0m")
		check, err := goalchecker.NewGoalChecker(goalchecker.Options{RawVideoPath: path})
		if err != nil {
			log.Fatalln(err)
		}
		sec, err := check.FindGoalMinuteSecondVoiceFrequency()
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("sec: %v\n", sec)
		return
	}

	var min string
	fmt.Print("Gol dakikası: ")
	fmt.Scanln(&min)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	check, err := goalchecker.NewGoalChecker(goalchecker.Options{RawVideoPath: path, GoalMinute: min})
	if err != nil {
		log.Fatalln(err)
	}
	sec, err := check.FindGoalSecondOCR(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("sec: %v\n", sec)
	fmt.Printf("goalchecker.IntToFfmpegInt(sec): %v\n", goalchecker.IntToFfmpegInt(sec))
}

func videodurationtest() {
	path := os.Args[2]
	out, err := goalchecker.GetVideoDuration(path)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("out: %v\n", out)

}

func lastframetest() {
	path := os.Args[2]
	outpath := "frame.png"

	rawduration, err := goalchecker.GetVideoDuration(path)
	if err != nil {
		log.Fatalln(err)
	}

	fduration, err := strconv.ParseFloat(rawduration, 32)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("fduration: %v\n", fduration)

	var counter = 0
birkezdaha:

	lastframesec, err := goalchecker.TryGetLastFrameFromVideo(int(fduration), path, outpath)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("lastframesec: %v\n", lastframesec)

	file, err := os.Open(outpath)
	if err != nil {
		log.Fatalln(err)

	}
	defer file.Close()

	ocrveri, err := ocr.NewOcrSpaceOcr().Text(context.Background(), file)
	if err != nil {
		log.Fatalln(err)
	}
	if ocrveri == "" && counter == 0 {
		log.Println("Goto with fduration - 10")
		counter++
		fduration = fduration - 10
		goto birkezdaha
	}
	fmt.Printf("ocrveri: %v\n", ocrveri)
}
