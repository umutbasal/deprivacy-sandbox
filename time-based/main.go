package main

import (
	_ "embed"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

//go:embed id.html
var idhtml string

//go:embed id.js
var idjs string

//go:embed id-worklet.js
var idworklet string

type sessionTimeCapture map[string]map[string]time.Time

func main() {

	app := fiber.New()

	sessionTimeCapture := make(sessionTimeCapture)

	app.Use(func(c *fiber.Ctx) error {
		c.Set("Supports-Loading-Mode", "fenced-frame")
		c.Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		c.Set("Pragma", "no-cache")

		return c.Next()
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Deprivacy Sandbox - Time-based")
	})

	app.Get("/id.html", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.Send([]byte(idhtml))
	})

	app.Get("/id.js", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/javascript")
		id, err := uuid.NewRandom()
		id = uuid.Must(id, err)
		idjss := fmt.Sprintf("const uuid = '%s';\n%s", id, idjs)
		return c.Send([]byte(idjss))
	})

	app.Get("/id-worklet.js", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/javascript")
		return c.Send([]byte(idworklet))
	})

	app.Get("/time-capture", func(c *fiber.Ctx) error {
		id := c.Query("id")
		if id == "" {
			return c.SendStatus(400)
		}
		index := c.Query("index")
		if index == "" {
			return c.SendStatus(400)
		}

		if _, ok := sessionTimeCapture[id]; !ok {
			sessionTimeCapture[id] = make(map[string]time.Time)
		}

		sessionTimeCapture[id]["index-"+index] = time.Now()

		return c.SendString("OK")
	})

	app.Get("/time-capture-start", func(c *fiber.Ctx) error {
		id := c.Query("id")
		if id == "" {
			return c.SendStatus(400)
		}

		if _, ok := sessionTimeCapture[id]; !ok {
			sessionTimeCapture[id] = make(map[string]time.Time)
		}

		sessionTimeCapture[id]["start"] = time.Now()

		return c.SendString("OK")
	})

	app.Get("/id", func(c *fiber.Ctx) error {
		id := c.Query("id")
		if id == "" {
			return c.SendStatus(400)
		}

		if _, ok := sessionTimeCapture[id]; !ok {
			return c.SendStatus(404)
		}

		if ok := len(sessionTimeCapture[id]) == 5; !ok {
			return c.SendStatus(404)
		}

		var ds []time.Duration
		ds = append(ds, sessionTimeCapture[id]["index-0"].Sub(sessionTimeCapture[id]["start"]))
		ds = append(ds, sessionTimeCapture[id]["index-1"].Sub(sessionTimeCapture[id]["index-0"]))
		ds = append(ds, sessionTimeCapture[id]["index-2"].Sub(sessionTimeCapture[id]["index-1"]))
		ds = append(ds, sessionTimeCapture[id]["index-3"].Sub(sessionTimeCapture[id]["index-2"]))

		fmt.Println(ds)
		fmt.Println(sessionTimeCapture[id])

		// print unix
		for k, v := range sessionTimeCapture[id] {
			fmt.Println(k, v.Unix())
		}

		ids := durationsToIds(ds)

		return c.JSON(ids)
	})

	app.Listen(":8080")

}

// [
//     "2024-01-26T16:47:20.547433+03:00",
//     "2024-01-26T16:47:20.574793+03:00",
//     "2024-01-26T16:47:20.870537+03:00",
//     "2024-01-26T16:47:21.470618+03:00",
//     "2024-01-26T16:47:22.370607+03:00"
// ]

func durationsToIds(durations []time.Duration) []string {
	ids := make([]string, 0)
	for _, duration := range durations {
		ids = append(ids, fmt.Sprintf("%d", int(duration.Round(time.Second).Seconds())))
	}
	return ids
}
