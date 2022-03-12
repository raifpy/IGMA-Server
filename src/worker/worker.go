package worker

import (
	"soccerapi/src/worker/database"
	"soccerapi/src/worker/logger"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type Options struct {
	FiberWsPath string
	LogDirPath  string

	Mongo       *mongo.Client
	Fiber       *fiber.App
	RouterGroup fiber.Router
}

type Worker struct {
	Options     Options
	Fiber       *fiber.App
	RouterGroup fiber.Router
	MongoDB     *database.MongoDB
	Log         *logger.Log
	WorkerMap   *WorkerMap
}

func NewWorker(o Options) (w *Worker, err error) {
	w = &Worker{
		Options:     o,
		WorkerMap:   NewWorkerMap(),
		RouterGroup: o.RouterGroup,
	}

	if w.Log, err = logger.NewLog(o.LogDirPath); err != nil {
		return
	}
	w.Fiber = o.Fiber
	/*if w.MongoDB, err = database.NewMongoDB(o.MongoAddr); err != nil {
		return
	}*/
	w.MongoDB = &database.MongoDB{
		Database: o.Mongo,
	}

	w.RouteFiber()
	return

}
