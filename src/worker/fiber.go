package worker

func (w *Worker) RouteFiber() {
	w.RouteWsGroup()
}

func (w *Worker) RouteWsGroup() {
	w.RouterGroup.Use(w.fiberMiddlewareWorker)
	w.RouterGroup.Get(w.Options.FiberWsPath, w.fiberWorkerWs())
	w.RouterGroup.Post("/upload", w.fiberMediaUpload)

	/*workergroup := w.Fiber.Group("/worker", w.fiberMiddlewareWorker)
	workergroup.Get(w.Options.FiberWsPath, w.fiberWorkerWs())
	workergroup.Post("/upload", w.fiberMediaUpload)*/

}
