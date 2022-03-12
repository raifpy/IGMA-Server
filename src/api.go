package soccer

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

var errUnknownToken = errors.New("invalid token")

//var errNeedMatchidAndUpdateid = errors.New("match_id or update_id parameters valid")
var errNeedMatchidAndClipid = errors.New("match_id or clip_id parameters valid")
var errNeedMatchid = errors.New("match_id parameters valid")
var errExpired = errors.New("expired token")

func (s *Soccer) ApiHandle() {
	s.Fiber.Get("/updates", func(c *fiber.Ctx) error {
		// AUTH
		g, err := s.GetMatchsJSONFiber()
		if err != nil {
			return err
		}
		c.Set("Content-Type", "application/json")
		return c.Send(g)
	})

	s.Fiber.Get("/worker/download", s.Dsystem.FiberHandler())

	//s.Worker.RouterGroup.Get("/download", s.Dsystem.FiberHandler())

	paymentgroup := s.Fiber.Group("/service/:token/", func(c *fiber.Ctx) error {

		token := c.Params("token")
		if token == "" {
			log.Println("token empty")
			return errUnknownToken
		}
		user, err := s.GetUserFromToken(c.Context(), token)
		if err != nil {
			log.Printf("token %s kayıtlı değil\n", token)
			return errUnknownToken
		}
		//if time.Until(user.Expired)

		if user.Expired.Before(time.Now()) {
			return errExpired
		}

		c.Locals("user", user)
		return c.Next()
	})
	paymentgroup.Get("/getMe", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Send(c.Locals("user").(UserClient).ToJSON())
	})

	paymentgroup.Get("/getClips", func(c *fiber.Ctx) error {
		matchid := c.Query("match_id")
		//updateid := c.Query("update_id")

		if matchid == "" /* || updateid == "" */ {
			return errNeedMatchid
		}
		mid, err := strconv.Atoi(matchid)
		if err != nil {
			log.Println("matchid int değil anlaşılan")
			return errNeedMatchid
		}
		res, err := s.GetMatchClips(int64(mid))
		if err != nil {
			log.Println("GetMatchClips: ", err)
			return errNeedMatchid
		}
		c.Set("Content-Type", "application/json")
		return c.Send(res.ToClientJSON())

	})

	paymentgroup.Get("/getClip", func(c *fiber.Ctx) error {
		matchid := c.Query("match_id")
		clipid := c.Query("clip_id")

		fmt.Printf("matchid: %v\n", matchid)
		fmt.Printf("clipid: %v\n", clipid)

		if matchid == "" || clipid == "" {
			return errNeedMatchidAndClipid
		}

		mid, err := strconv.Atoi(matchid)
		if err != nil {
			log.Println("matchid int değil anlaşılan")
			return errNeedMatchidAndClipid
		}
		cid, err := strconv.Atoi(clipid)
		if err != nil {
			log.Println("clipid int değil anlaşılan")
			return errNeedMatchidAndClipid
		}

		res, err := s.GetMatchClips(int64(mid))
		if err != nil {
			log.Println("GetMatchClips: ", err)
			return errNeedMatchidAndClipid
		}

		var clip = MatchClip{}
		var ok bool

		for _, v := range res.Clips {
			if v.Client.ClipID == int64(cid) {

				ok = true
				clip = v
				break
			}

		}

		if !ok {
			return errNeedMatchidAndClipid
		}

		fmt.Printf("clip: %v\n", clip)

		if clip.Path != "" {
			fmt.Printf("v.Path: %v\n", clip.Path)
			/*file, err := os.Open(clip.Path) // Burada sadece local'den dosya dönüyor!
			if err != nil {
				log.Println("Eyvah, db'de eşleşti ancak file olarak yok! ", clip.Path)
				return err
			}*/

			c.Set("Content-Type", "video/mp4")

			//c.Set("Content-Length", fmt.Sprint(stat.Size()))
			//return c.SendStream(file, int(stat.Size()))
			return c.SendFile(clip.Path)
		}

		return c.SendString("video erişilebilir değil!")

	})

}
