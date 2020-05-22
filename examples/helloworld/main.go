package main

import (
	"github.com/gofiber/fiber"
	"github.com/iesreza/io"
)

type MyModel struct {
	io.Model
	Name     string
	Username string
}

func main() {
	io.Setup()

	db := io.GetDBO()
	db.AutoMigrate(MyModel{})

	io.Get("test", func(ctx *fiber.Ctx) {
		request := io.Upgrade(ctx)

		obj := MyModel{
			Name: "Allan", Username: "allan@ies",
		}

		db.Create(&obj)

		request.WriteResponse(obj)
	})

	io.Get("/", func(ctx *fiber.Ctx) {
		ctx.Write("Hello World")
	})
	io.Run()
}
