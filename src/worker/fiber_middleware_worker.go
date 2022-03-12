package worker

import (
	"errors"
	"log"

	"soccerapi/src/worker/types"

	"strconv"

	"github.com/gofiber/fiber/v2"
)

var errRequestIdNotValid = errors.New("id not valid")

//var errRequestTokenExpried = errors.New("expired")

func (w *Worker) fiberMiddlewareWorker(c *fiber.Ctx) error {
	log.Println("Ws GET: ", c.IP())
	token := c.Get("token")
	id, err := strconv.Atoi(c.Get("id"))

	if token == "" || len(token) < 20 || err != nil {
		log.Println("token ya da id elendi")
		c.Status(fiber.StatusUnauthorized)
		return fiber.ErrUnauthorized
	}

	res := w.MongoDB.Database.Database("worker").Collection("workers").FindOne(c.Context(), map[string]interface{}{
		"token": token,
	})

	var t types.DatabaseAuth
	if err := res.Decode(&t); err != nil {
		c.Locals("custom_error", fiber.ErrUnauthorized)
		c.Status(fiber.StatusUnauthorized)
		return err
	}

	if id != int(t.Id) {
		c.Status(fiber.StatusUnauthorized)
		return errRequestIdNotValid
	}

	/*m := time.Until(t.Until).Minutes()
	if m < 0 {
		c.Status(fiber.StatusUnauthorized)
		return errRequestTokenExpried
	}*/ // IP kontrol

	c.Locals("worker", t)

	log.Println("c.Next() auth")

	return c.Next()

}
