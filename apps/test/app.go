package test

import (
	"fmt"
	"github.com/CloudyKit/jet"
	"github.com/gofiber/fiber"
	"github.com/iesreza/io"
	"github.com/iesreza/io/lib"
	"github.com/iesreza/io/lib/gpath"
	"github.com/iesreza/io/lib/text"
	"github.com/iesreza/io/menu"
	"github.com/iesreza/io/user"
	"github.com/jinzhu/gorm"
)

func Register() {
	fmt.Println("Test Registered")
	io.Register(App{})
}

var db *gorm.DB
var Path string

type App struct{}

var views *jet.Set

func (App) Register() {
	//Require auth
	db = io.GetDBO()
	Path = gpath.Parent(gpath.WorkingDir()) + "/apps/test"

}

func (App) Router() {
	io.Get("/admin/list", func(ctx *fiber.Ctx) {
		r := io.Upgrade(ctx)
		r.View(nil, "test.list", "template.login")
	})

}

func (App) Permissions() []user.Permission {
	return []user.Permission{}
}

func (App) Menus() []menu.Menu {
	return []menu.Menu{}
}
func (App) WhenReady() {

	db.AutoMigrate(MyModel{})
	for i := 0; i < 20; i++ {
		item := MyModel{}
		item.Username = text.Random(6)
		item.Type = lib.RandomBetween(1, 3)
		item.Name = text.Random(8)
		item.Group = lib.RandomBetween(1, 4)
		item.Alias = text.Random(4)
		db.Create(&item)
	}
}