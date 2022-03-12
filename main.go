package main

import (
	"log"
	"math/rand"
	"os"
	soccer "soccerapi/src"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

var (
	Host      string = "localhost:2095"
	MongoAddr string = "mongodb+srv://mongodbaddress/database"
)

func init() {
	rand.Seed(time.Now().Unix())
	os.Setenv("TZ", "Europe/Istanbul")
}

func Main() {
	_soccer, err := soccer.NewSoccer(soccer.Config{
		FiberHost:    Host,
		MongoAddr:    MongoAddr,
		WorkerWsPath: "/wsworker",
		ErrRaport: func(e error) {
			log.Printf("ErrorReport: %s\n", e)
		},
		SupportedTvChannels: MatchChannelList,
		IpTvConfigPath:      "iptvconf.json",
	})
	if err != nil {
		log.Fatalf("main.go err: %v \n ", err)

	}

	if os.Getenv("GOALTEST") != "" { //? Sistem testi | Kötü kodlama
		log.Println("\033[31mGOAL TEST\033[0m")
		go goaltest(_soccer)
	}

	_soccer.Run()
}

func main() { //? Kötü kodlama
	if len(os.Args) > 1 {
		test()
		os.Exit(2)
	}

	Main()
}
