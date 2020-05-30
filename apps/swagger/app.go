package apiman

import (
	"fmt"
	"github.com/iesreza/io"
	"github.com/iesreza/io/menu"
	"github.com/iesreza/io/user"
)

func Register() {
	io.Register(App{})
}

type App struct{}

var config *io.Configuration

// Register the adminlte
func (App) Register() {
	fmt.Println("API Man Registered")
	config = io.GetConfig()
}

// WhenReady called after setup all apps
func (App) WhenReady() {}

// Router setup routers
func (App) Router() {

}

// Permissions setup permissions of app
func (App) Permissions() []user.Permission { return []user.Permission{} }

// Menus setup menus
func (App) Menus() []menu.Menu {
	return []menu.Menu{}
}
