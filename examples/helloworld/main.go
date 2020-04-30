package main

import (
	"github.com/gofiber/fiber"
	"github.com/iesreza/io"
)

func main() {
	io.Setup()
	io.Get("/", func(ctx *fiber.Ctx) {
		ctx.Write("Hello World")
	})
	io.Run()
}
