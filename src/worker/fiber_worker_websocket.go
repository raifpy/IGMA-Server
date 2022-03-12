package worker

import (
	"errors"
	"log"
	"time"

	"soccerapi/src/worker/types"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/websocket/v2"
)

var errDuplicateWsConnection = errors.New("duplicate ws connection")

func (w *Worker) fiberWorkerWs() fiber.Handler {

	return websocket.New(func(c *websocket.Conn) {

		t := c.Locals("worker").(types.DatabaseAuth)
		if client := w.WorkerMap.Get(t.Id); client != nil && client.Conn != nil {
			log.Println("zaten var & bağlı")
			c.SetReadDeadline(time.Now().Add(time.Second * 4))
			c.WriteJSON(types.WebsocketContact{
				Type: "close",
				Error: &types.WebsocketError{
					Error: errDuplicateWsConnection.Error(),
				},
			})
			c.Close()
			return
		}

		wc := &WorkerConn{
			Id:   t.Id,
			Conn: c,
		}
		w.WorkerMap.Set(wc)
		defer func() {
			wc.CloseConnection()
			c.Close()

			w.WorkerMap.Del(&WorkerConn{
				Id: t.Id,
			})
		}()

		log.Printf("\033[32mWorker bağlandı! %d \033[0m\n", wc.Id)

		for {
			var wt types.WebsocketContact
			if err := c.ReadJSON(&wt); err != nil {
				log.Printf("ID: %d close: %v", t.Id, err)
				break
			}
			go w.WorkerStatusUpdate(wt, wc)
		}

	})
}
