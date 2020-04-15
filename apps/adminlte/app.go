package adminlte

import (
	"fmt"
	"github.com/CloudyKit/jet"
	"github.com/iesreza/io"
	"github.com/iesreza/io/lib/gpath"
	"github.com/iesreza/io/user"
)

func Register()  {
	io.Register(App{})
}

var  Path string
type App struct {}
var PageVars  = map[string]string{}
var views       *jet.Set
var config      *io.Configuration
func (App) Register() {
	fmt.Println("AdminLTE Registered")
	Path = gpath.Parent(gpath.WorkingDir())+"/apps/adminlte"
	views = io.RegisterView("template",Path+"/pages")
	config = io.GetConfig()
	views.AddGlobal("title",config.App.Name)
}

func (App) Router() {
	io.Static("/assets",Path+"/assets")
	io.Static("/plugins",Path+"/plugins")
}
func (App) Permissions() []user.Permission {return []user.Permission{}}

