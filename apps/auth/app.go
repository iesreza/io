package auth

import (
	"fmt"
	"github.com/iesreza/io"
	"github.com/iesreza/io/apps/query"
	"github.com/iesreza/io/menu"
	"github.com/iesreza/io/user"
	"github.com/jinzhu/gorm"
)

func Register() {
	io.Register(App{})
}

var db *gorm.DB
var config *io.Configuration

// App settings app struct
type App struct{}

// Register register the auth in io apps
func (App) Register() {

	fmt.Println("Auth Registered")
	userFilter := query.Filter{
		Object: &user.User{},
		Slug:   "user",
		Allowed: map[string]string{
			"id":         ``,
			"username":   `validate:"format=username"`,
			"name":       `validate:"format=text"`,
			"email":      `validate:"format=text"`,
			"group_id":   `validate:"format=numeric"`,
			"created_at": `validate:"format=date"`,
		},
	}
	query.Register(userFilter)

	db = io.GetDBO()
	config = io.GetConfig()
	if config.Database.Enabled == false {
		panic("Auth App require database to be enabled. solution: enable database at config.yml")
	}
}

// Router setup routers
func (App) Router() {
	controller := Controller{}

	auth := io.Group("/auth")
	//=>
	auth.Post("/user/login", controller.Login)

	auth.Post("/user/create", controller.CreateUser)
	auth.Post("/user/edit/:id", controller.EditUser)
	auth.Get("/user/me", controller.GetMe)
	auth.Get("/user/all/:offset/:limit", controller.GetAllUsers)
	auth.Get("/user/:id", controller.GetUser) //this should be always last router else it will match before others

	auth.Post("/role/create", controller.CreateRole)
	auth.Post("/role/edit/:id", controller.EditRole)
	auth.Get("/role/all", controller.GetRoles)
	auth.Get("/role/:id", controller.GetRole)
	auth.Get("/role/:id/groups", controller.GetRoleGroups)

	auth.Post("/group/create", controller.CreateGroup)
	auth.Post("/group/edit/:id", controller.EditGroup)
	auth.Get("/group/all", controller.GetGroups)
	auth.Get("/group/:id", controller.GetGroup)

	auth.Get("/permission/all", controller.GetAllPermissions)

}

// Permissions setup permissions of app
func (App) Permissions() []user.Permission {
	return []user.Permission{
		{Title: "Access Users", CodeName: "user.view", Description: "Access list to view list of users"},
		{Title: "Create Users", CodeName: "user.create", Description: "Create new user"},
		{Title: "Edit users", CodeName: "user.edit", Description: "Edit user data"},
		{Title: "Remove users", CodeName: "user.remove", Description: "Remove user data"},
		{Title: "Login as user", CodeName: "user.loginas", Description: "Login as another user without credentials"},
		{Title: "Limit to subgroup", CodeName: "user.limited", Description: "User can only access/modify subgroup users"},
		{Title: "Access Groups", CodeName: "group.view", Description: "Access to groups data"},
		{Title: "Create Groups", CodeName: "group.create", Description: "Create new group"},
		{Title: "Edit Groups", CodeName: "group.edit", Description: "Edit groups data"},
		{Title: "Remove Groups", CodeName: "group.remove", Description: "Remove groups"},
		{Title: "Access Roles", CodeName: "role.view", Description: "Access to roles"},
		{Title: "Create Role", CodeName: "role.create", Description: "Create new role"},
		{Title: "Edit Roles", CodeName: "role.edit", Description: "Edit roles data"},
		{Title: "Remove Roles", CodeName: "role.remove", Description: "Remove roles"},
	}
}

// Menus setup menus
func (App) Menus() []menu.Menu {
	return []menu.Menu{}
}

// WhenReady called after setup all apps
func (App) WhenReady() {}
