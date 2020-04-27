package query

import (
	"fmt"
	"github.com/iesreza/io"
	"github.com/iesreza/io/menu"
	"github.com/iesreza/io/user"
	"github.com/jinzhu/gorm"
)

var c Controller

func Register(v ...Filter) {
	if len(v) == 0 {
		io.Register(App{})
		return
	}
	if objects.data == nil {
		objects.Init()
	}
	for _, item := range v {
		c.Register(item)
	}

}

func (App) WhenReady() {}

var db *gorm.DB

type App struct{}

func (App) Register() {
	fmt.Println("Filter Registered")
	db = io.GetDBO()
}

func (App) Router() {

}

func (App) Permissions() []user.Permission {
	return []user.Permission{}
}

func (App) Menus() []menu.Menu {
	return []menu.Menu{}
}
