package main

import (
	_ "embed"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	cmap "github.com/orcaman/concurrent-map/v2"
)

//go:embed id.html
var idhtml string

//go:embed id.js
var idjs string

//go:embed id-worklet.js
var idworklet string

type IdentityExtractionURLS []IdentifierCheckURLS

type IdentifierCheckURLS []IdentifierCheckURL

type IdentifierCheckURL struct {
	URL string `json:"url"`
}

type UserStore map[string]string

func main() {

	userStore := make(UserStore)
	identifierTriggered := cmap.New[string]()

	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Set("Supports-Loading-Mode", "fenced-frame")
		c.Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		c.Set("Pragma", "no-cache")

		return c.Next()
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Deprivacy Sandbox - Entropy based")
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

	app.Get("/identifier-check-urls", func(c *fiber.Ctx) error {
		id := c.Query("uuid")
		if id == "" {
			return c.SendStatus(400)
		}

		return c.JSON(IdentifierCheckURLS{
			IdentifierCheckURL{
				URL: "http://localhost:8080/unkown?uuid=" + id,
			},
			IdentifierCheckURL{
				URL: "http://localhost:8080/known?uuid=" + id,
			},
		})
	})

	app.Get("/identity-extraction-urls", func(c *fiber.Ctx) error {
		id := c.Query("uuid")
		if id == "" {
			return c.SendStatus(400)
		}

		digits := 4
		maxNum := 4

		idExtractUris := make(IdentityExtractionURLS, 0)

		for i := 0; i < digits; i++ {
			idChckerUris := make(IdentifierCheckURLS, 0)
			for j := 0; j < maxNum; j++ {
				idChckerUris = append(idChckerUris, IdentifierCheckURL{
					URL: fmt.Sprintf("http://localhost:8080/identity-extraction?uuid=%s&digit=%d&num=%d", id, i, j),
				})
			}
			idExtractUris = append(idExtractUris, idChckerUris)
		}

		return c.JSON(idExtractUris)
	})

	app.Get("/identity-extraction", func(c *fiber.Ctx) error {
		id := c.Query("uuid")
		if id == "" {
			return c.SendStatus(400)
		}

		digit := c.Query("digit")
		if digit == "" {
			return c.SendStatus(400)
		}

		num := c.Query("num")
		if num == "" {
			return c.SendStatus(400)
		}

		fmt.Println("Identity extraction triggered:" + id + " digit:" + digit + " num:" + num)

		userStore[id] = userStore[id] + "digit:" + digit + " num:" + num + ";"

		return c.SendString("digit:" + digit + " num:" + num)
	})

	app.Get("/identity-extraction-result", func(c *fiber.Ctx) error {
		id := c.Query("uuid")
		if id == "" {
			return c.SendStatus(400)
		}

		if _, ok := userStore[id]; !ok {
			return c.SendStatus(404)

		}

		return c.JSON(userStore[id])
	})

	app.Get("/unkown", func(c *fiber.Ctx) error {
		id := c.Query("uuid")
		if id == "" {
			return c.SendStatus(400)
		}

		if _, ok := identifierTriggered.Get(id); !ok {
			identifierTriggered.Set(id, "unknown")
			fmt.Println("Identifier triggered unknown:" + id)
		}

		c.Set("Content-Type", "text/html")
		return c.SendString(`<html><body><script>window.sharedStorage.set('has-identifier', 1);window.sharedStorage.set('identifier-1', 1);window.sharedStorage.set('identifier-2', 3);window.sharedStorage.set('identifier-3', 2);window.sharedStorage.set('identifier-4', 2)</script></body></html>`)
	})

	app.Get("/known", func(c *fiber.Ctx) error {
		id := c.Query("uuid")
		if id == "" {
			return c.SendStatus(400)
		}

		if _, ok := identifierTriggered.Get(id); !ok {
			identifierTriggered.Set(id, "known")
			fmt.Println("Identifier triggered known:" + id)
		}

		return c.SendString("known")
	})

	app.Get("/is-identified", func(c *fiber.Ctx) error {
		id := c.Query("uuid")
		if id == "" {
			return c.SendStatus(400)
		}

		if v, ok := identifierTriggered.Get(id); ok {
			fmt.Println("Identifier available:" + id)
			if v == "unknown" {
				return c.SendStatus(200)
			}
		}

		return c.SendStatus(404)
	})

	app.Listen(":8080")

}
