package main

import (
	"backend-development/pkg/handlers"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func healthcheck(c *fiber.Ctx) error {
	return c.SendString("(Deliverables) Backend Development Assignment GDSC UGM Hacker Role Member")
}

func main() {
	app := fiber.New()

	app.Use("/", func(c *fiber.Ctx) error {
		fmt.Println("Server is listening on port 7070")
		return c.Next()
	})

	app.Get("/", healthcheck)

	app.Post("/v1/posts", handlers.CreatePost)
	app.Get("/v1/posts", handlers.GetAllPosts)
	app.Get("/v1/posts/:id", handlers.GetSinglePost)
	app.Patch("/v1/posts/:id", handlers.UpdatePost)
	app.Delete("/v1/posts/:id", handlers.DeletePost)

	log.Fatal(app.Listen(":7070"))
}
