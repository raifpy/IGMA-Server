package dsystem

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func (d *Dsystem) FiberHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Query("id")
		if id == "" {
			log.Println("dsystem: id boş")
			return fiber.ErrUnauthorized
		}
		v, ok := d.Get(id)
		if !ok {
			log.Println("dsystem: id eşleşmedi")
			return fiber.ErrUnauthorized
		}
		fmt.Printf("v: %v\n", v)
		if (v.Ip != "" && v.Ip != c.IP()) || v.Token != c.Query("token") {
			log.Println("dsystem: tokenler eş değil")
			return fiber.ErrUnauthorized
		}

		/*file, err := os.Open(v.Path)
		if err != nil {
			return err
		}
		time.AfterFunc(time.Minute*6, func() {
			file.Close() //!!! Leak?
		})
		stat, err := file.Stat()
		if err != nil {
			return err
		}

		return c.SendStream(file, int(stat.Size()))*/
		return c.SendFile(v.Path)

	}
}
