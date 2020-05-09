package main

///
import (
	"fmt"

	"github.com/gofiber/fiber"
	"github.com/iesreza/io"
)

///
func main() {
	io.Setup()
	io.Get("/", func(ctx *fiber.Ctx) {

		fmt.Println("DB Examples Registered")

		r := io.Upgrade(ctx)
		//b = io.GetDBO()
		/*type User struct {
			username  	string
			firstname  	string
		}
		user_response = User{
			"allan_nava"
			"Allan"
		}*/
		//r.WriteResponse(user_response)*/
		r.WriteResponse("Trying to add db")
	})
	io.Run()
}

///
