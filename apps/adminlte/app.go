package adminlte

import (
	"fmt"
	"github.com/CloudyKit/jet"
	"github.com/gofiber/fiber"
	"github.com/iesreza/io"
	"github.com/iesreza/io/apps/settings"
	"github.com/iesreza/io/lib/gpath"
	"github.com/iesreza/io/menu"
	"github.com/iesreza/io/user"
	"github.com/iesreza/io/viewfn"
)

func Register() {
	io.Register(App{})
}

var Path string

type App struct{}

var setting = &Settings{
	NavbarColor:     "green",
	NavbarVariation: "dark",
	SidebarColor:    "green",
}
var pages *jet.Set
var elements *jet.Set
var config *io.Configuration

func (App) Register() {
	fmt.Println("AdminLTE Registered")
	Path = gpath.Parent(gpath.WorkingDir()) + "/apps/adminlte"
	pages = io.RegisterView("template", Path+"/pages")
	elements = io.RegisterView("html", Path+"/html")
	config = io.GetConfig()
	settings.Register("AdminLTE Template", setting)

}

func (App) WhenReady() {
	pages.AddGlobal("title", config.App.Name)
	pages.AddGlobal("nav", io.AppMenus)
	pages.AddGlobal("settings", setting)
	viewfn.Bind(pages, "thumb")
}

func (App) Router() {
	io.Static("/assets", Path+"/assets")
	io.Static("/plugins", Path+"/plugins")

	io.Get("", func(ctx *fiber.Ctx) {
		r := io.Upgrade(ctx)
		r.Var("heading", "Test")
		r.View(nil, "template.default")
	})

	io.Get("/test", func(ctx *fiber.Ctx) {
		r := io.Upgrade(ctx)
		r.Var("heading", "Test1")
		r.View(nil, "template.default")
	})

	io.Get("/test2", func(ctx *fiber.Ctx) {
		r := io.Upgrade(ctx)
		r.Var("heading", "Test2")
		r.View(nil, "template.default")
	})
}
func (App) Permissions() []user.Permission { return []user.Permission{} }

func (App) Menus() []menu.Menu {
	return []menu.Menu{}
}
