package soccer

import (
	"crypto/tls"
	"log"
	"net/http"
	dsystem "soccerapi/src/downloadsystem"
	"soccerapi/src/iptv"
	"soccerapi/src/scores"
	"soccerapi/src/worker"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

type Config struct {
	FiberHost      string
	MongoAddr      string
	WorkerWsPath   string
	IpTvConfigPath string

	SupportedTvChannels []string
	ErrRaport           func(error)
}
type Soccer struct {
	config Config
	Fiber  *fiber.App
	IpTv   *iptv.IpTvStand

	Worker           *worker.Worker
	Mongo            *mongo.Client
	Dsystem          *dsystem.Dsystem
	GameScoreWatcher *GameScoreWatcher //TODO ekle
	MatchApi         *MatchApi
	OnGoalError      func(err error, score scores.PlainGameScore, fatal bool, ek ...string)
	ErrRaport        func(error)
}

func (s *Soccer) Run() error {
	return s.Fiber.Listen(s.config.FiberHost)
}

func NewSoccer(c Config) (sc *Soccer, err error) {
	sc = &Soccer{
		config:    c,
		ErrRaport: c.ErrRaport,

		//Worker:    &Worker{},
	}
	sc.Fiber = fiber.New(fiber.Config{
		BodyLimit:               1204 * 1024 * 100,
		TrustedProxies:          []string{"127.0.0.1"}, //!! REVERSE PROXY için
		ProxyHeader:             "X-Forwarded-For",     //!! REVERSE PROXY için
		EnableTrustedProxyCheck: true,                  //!! REVERSE PROXY için

		ErrorHandler: func(c *fiber.Ctx, e error) error {
			log.Printf("Error: %v - path: %s", e, c.Request().URI().String())
			if err, ok := c.Locals("custom_error").(error); ok {
				e = err
			}
			if c.Context().Response.StatusCode() == 200 {
				c.Status(400)
			}
			return c.SendString("error: " + e.Error())
		},
	})

	sc.Fiber.Use(logger.New())

	if sc.Mongo, err = mongo.NewClient(options.Client().ApplyURI(sc.config.MongoAddr)); err != nil {
		return
	}
	timeoutcontext, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err = sc.Mongo.Connect(timeoutcontext); err != nil {
		return
	}

	sc.Worker, err = worker.NewWorker(worker.Options{
		Fiber:       sc.Fiber,
		Mongo:       sc.Mongo,
		FiberWsPath: c.WorkerWsPath,
		RouterGroup: sc.Fiber.Group("/worker"),
	})
	if err != nil {
		return nil, err
	}
	sc.Dsystem = dsystem.NewDSystem()
	sc.IpTv, err = iptv.NewIpTvStand(iptv.Options{
		IpTvConfigPath: c.IpTvConfigPath,
		OnUpdate: func() {
			log.Printf("%s güncellendi\n", sc.config.IpTvConfigPath)
		},
		OnError: func(e error) { log.Printf("\033[31m%s update error: %v\033[0m", sc.config.IpTvConfigPath, e) },
		HttpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	sc.GameScoreWatcher, err = NewGameScoreWatcher(GameScoreWatcherOptions{
		IpTvChannelList: c.SupportedTvChannels,
		OnGoal:          sc.OnGoal,
		OnRequestError:  sc.ErrRaport,
		OnCanceledGoal:  sc.OnCanceledGoal,
	})
	if err != nil {
		return nil, err
	}

	sc.SetMatchApi()

	//log.Printf("GameScoreWatcher up\n")

	//sc.Worker.init()

	if sc.ErrRaport == nil {
		sc.ErrRaport = func(e error) {
			log.Printf("\033[31m[error] %v\033[0m", e)
		}
	}

	sc.OnGoalError = sc.OnGoalErrorFn
	sc.GameScoreWatcher.Scores.OnRequestError = func(e error) {
		log.Println("Soccer Request errror: ", err)
	}

	sc.ApiHandle()

	return
}
