package admin

import (
	"fmt"
	"github.com/CloudyKit/jet"
	"github.com/gofiber/fiber"
	"github.com/iesreza/io"
	"github.com/iesreza/io/apps/auth"
	"github.com/iesreza/io/apps/settings"
	"github.com/iesreza/io/lib/fontawesome"
	"github.com/iesreza/io/lib/gpath"
	"github.com/iesreza/io/menu"
	"github.com/iesreza/io/user"
)

func Register() {
	fmt.Println("Dashboard Registered")
	io.Register(App{})
}

var db = io.GetDBO()
var Path string

type App struct{}

var views *jet.Set
var setting Settings

func (App) Register() {
	//Require auth
	setting.SessionAge = fmt.Sprint(io.GetConfig().JWT.Age.Seconds())
	Path = gpath.Parent(gpath.WorkingDir()) + "/apps/admin"
	auth.Register()
	views = io.RegisterView("admin", Path+"/views")
	settings.Register("Admin Panel", &setting)
}

func (App) Router() {
	io.Get("/admin/login", func(ctx *fiber.Ctx) {
		r := io.Upgrade(ctx)
		r.View(nil, "admin.login", "template.login")
	})

	io.Get("/admin/dashboard", func(ctx *fiber.Ctx) {
		r := io.Upgrade(ctx)
		r.View(nil, "template.default")
	})
}

func (App) Permissions() []user.Permission {
	return []user.Permission{
		{Title: "Login to admin", CodeName: "login", Description: "Able login to admin panel"},
		{Title: "Edit own dashboard", CodeName: "dashboard", Description: "Able edit own dashboard otherwise it inherit dashboard assigned by admin"},
	}
}

func (App) Menus() []menu.Menu {
	return []menu.Menu{
		{Title: "Dashboard", Url: "admin/dashboard", Icon: fontawesome.Home},
		{Title: "Parent", Url: "admin/dashboard", Icon: fontawesome.SearchPlus, Children: []menu.Menu{
			{Title: "Child 1", Url: "admin/dashboard", Icon: fontawesome.Image},
			{Title: "Child 2", Url: "admin/dashboard", Icon: fontawesome.Save},
		}},
	}
}
func (App) WhenReady() {}
