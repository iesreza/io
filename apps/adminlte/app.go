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

// Register register the adminlte in io apps
func Register() {
	io.Register(App{})
}

//Path to adminlte app
var Path string

// App adminlte app struct
type App struct{}

var setting = &Settings{
	NavbarColor:     "green",
	NavbarVariation: "dark",
	SidebarColor:    "green",
}
var pages *jet.Set
var elements *jet.Set
var config *io.Configuration

// Register the adminlte
func (App) Register() {
	fmt.Println("AdminLTE Registered")
	Path = gpath.Parent(gpath.WorkingDir()) + "/apps/adminlte"
	pages = io.RegisterView("template", Path+"/pages")
	elements = io.RegisterView("html", Path+"/html")
	config = io.GetConfig()
	settings.Register("AdminLTE Template", setting)

}

// WhenReady called after setup all apps
func (App) WhenReady() {
	pages.AddGlobal("title", config.App.Name)
	pages.AddGlobal("nav", io.AppMenus)
	pages.AddGlobal("settings", setting)
	viewfn.Bind(pages, "thumb")
}

// Router setup routers
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

// Permissions setup permissions of app
func (App) Permissions() []user.Permission { return []user.Permission{} }

// Menus setup menus
func (App) Menus() []menu.Menu {
	return []menu.Menu{}
}
