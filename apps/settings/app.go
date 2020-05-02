package settings

import (
	"fmt"
	"github.com/CloudyKit/jet"
	"github.com/iesreza/io"
	"github.com/iesreza/io/lib/concurrent"
	"github.com/iesreza/io/lib/fontawesome"
	"github.com/iesreza/io/lib/gpath"
	"github.com/iesreza/io/menu"
	"github.com/iesreza/io/user"
	"github.com/jinzhu/gorm"
	"reflect"
)

var controller Controller
var settings = concurrent.Map{}
var db *gorm.DB
var config *io.Configuration
var views *jet.Set
var Path string

// App settings app struct
type App struct{}

var initiated = false

// Register register the auth in io apps
func Register(v ...interface{}) {
	if len(v) == 0 {
		io.Register(App{})
		return
	}
	if initiated == false {
		Register()
	}
	var title string
	var object interface{}
	for _, item := range v {
		ref := reflect.ValueOf(item)
		switch ref.Kind() {
		case reflect.String:
			title = item.(string)
			break
		case reflect.Ptr:
			object = item
			break
		}
	}

	if title != "" && object != nil {
		controller.set(title, object)
	}

}

// Register settings app
func (App) Register() {
	fmt.Println("Settings Registered")
	settings.Init()
	db = io.GetDBO()
	config = io.GetConfig()
	if config.Database.Enabled == false {
		panic("Auth App require database to be enabled. solution: enable database at config.yml")
	}
	Path = gpath.Parent(gpath.WorkingDir()) + "/apps/settings"
	fmt.Println(Path + "/views")
	views = io.RegisterView("settings", Path+"/views")
	db.AutoMigrate(&Settings{})
}

// Router setup routers
func (App) Router() {
	controller := Controller{}
	io.Get("admin/settings", controller.view)
	io.Post("admin/settings/:name", controller.save)
	io.Post("admin/settings/reset/:name", controller.reset)
}

// Permissions setup permissions of app
func (App) Permissions() []user.Permission {
	return []user.Permission{
		{Title: "Access Settings", CodeName: "view", Description: "Access list to view list of settings"},
		{Title: "Modify Settings", CodeName: "modify", Description: "Modify Settings"},
	}
}

// Menus setup menus
func (App) Menus() []menu.Menu {
	return []menu.Menu{
		{Title: "Settings", Url: "admin/settings", Permission: "settings.view", Icon: fontawesome.Cog},
	}
}

// WhenReady called after setup all apps
func (App) WhenReady() {}
