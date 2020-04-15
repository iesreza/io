package admin

import (
	"fmt"
	"github.com/CloudyKit/jet"
	"github.com/gofiber/fiber"
	"github.com/iesreza/io"
	"github.com/iesreza/io/apps/auth"
	"github.com/iesreza/io/lib/gpath"
	"github.com/iesreza/io/user"
)

func Register()  {
	fmt.Println("Dashboard Registered")
	io.Register(App{})
}

var  db   = io.GetDBO()
var  Path string
type App struct {}
var views *jet.Set
func (App) Register() {
	//Require auth
	Path = gpath.Parent(gpath.WorkingDir())+"/apps/admin"
	auth.Register()
	views = io.RegisterView("admin",Path+"/views")



}

func (App) Router() {
	io.Get("/admin/login", func(ctx *fiber.Ctx) {
		r := io.Upgrade(ctx)
		r.View(nil,"admin.login","template.login")
	})
}



func (App) Permissions() []user.Permission {
	return []user.Permission{
		{Title:"Login to admin",CodeName:"login",Description:"Able login to admin panel"},
		{Title:"Edit own dashboard",CodeName:"dashboard",Description:"Able edit own dashboard otherwise it inherit dashboard assigned by admin"},

	}
}

