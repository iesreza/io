package io

import (
	"github.com/iesreza/io/lib/ref"
	"github.com/iesreza/io/menu"
	"github.com/iesreza/io/user"
)

type App interface {
	Register()
	Router()
	WhenReady()
	Permissions() []user.Permission
	Menus() []menu.Menu
}

var onReady = []func(){}
var apps = map[string]interface{}{}
var AppMenus = []menu.Menu{}

func Register(app App) {
	name := ref.Parse(app).Package

	//app already exist
	if _, ok := apps[name]; ok {
		return
	}

	apps[name] = app
	app.Register()
	app.Router()
	permissions := user.Permissions(app.Permissions())
	permissions.Sync(name)
	n := app.Menus()
	AppMenus = append(AppMenus, n...)

	onReady = append(onReady, app.WhenReady)
}

func GetRegisteredApps() map[string]interface{} {
	return apps
}
