package test

import (
	"fmt"
	"github.com/CloudyKit/jet"
	"github.com/iesreza/io"
	"github.com/iesreza/io/lib/gpath"
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
	views = io.RegisterView("test", Path+"/views")
}

func (App) Router() {
	io.Get("/admin/list", FilterViewController)
}

func (App) Permissions() []user.Permission {
	return []user.Permission{}
}

func (App) Menus() []menu.Menu {
	return []menu.Menu{}
}
func (App) WhenReady() {

	db.AutoMigrate(MyModel{}, MyGroup{})

	/*	db.Debug().Create(&MyGroup{
			Name:"Group 1",
		})
		db.Debug().Create(&MyGroup{
			Name:"Group 2",
		})
		db.Debug().Create(&MyGroup{
			Name:"Group 3",
		})*/
}
