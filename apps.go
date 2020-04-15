package io

import (
	"github.com/iesreza/io/lib/ref"
	"github.com/iesreza/io/user"
)

type App interface {
	Register()
	Router()
	Permissions() []user.Permission
}

var apps = map[string]interface{}{}

func Register(app App)  {
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

}

func GetRegisteredApps() map[string]interface{}  {
	return apps
}