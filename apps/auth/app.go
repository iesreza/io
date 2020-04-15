package auth

import (
	"fmt"
	"github.com/iesreza/io"
	"github.com/iesreza/io/apps/query"
	"github.com/iesreza/io/user"
	"github.com/jinzhu/gorm"
)

func Register()  {
	io.Register(App{})
}

var db *gorm.DB
var config *io.Configuration
type App struct {}
func (App) Register() {

	fmt.Println("Auth Registered")
	userFilter := query.Filter{
		Object:&user.User{},
		Slug: "user",
		Allowed: map[string]string{
			"id":``,
			"username":`validate:"format=username"`,
			"name":`validate:"format=text"`,
			"email":`validate:"format=text"`,
			"group_id":`validate:"format=numeric"`,
			"created_at":`validate:"format=date"`,
		},
	}
	userFilter.SetFilter("id != ?",4)
	query.Register(userFilter)

	db     = io.GetDBO()
	config = io.GetConfig()
	if config.Database.Enabled == false{
		panic("Auth App require database to be enabled. solution: enable database at config.yml")
	}
}

func (App) Router() {
	controller := Controller{}

	auth := io.Group("/auth")
	//=>
		auth.Post("/user/login",controller.Login)

		auth.Post("/user/create",controller.CreateUser)
	    auth.Post("/user/edit/:id",controller.EditUser)
		auth.Get("/user/me",controller.GetMe)
		auth.Get("/user/all/:offset/:limit",controller.GetAllUsers)
		auth.Get("/user/:id",controller.GetUser) //this should be always last router else it will match before others

		auth.Post("/role/create",controller.CreateRole)
	    auth.Post("/role/edit/:id",controller.EditRole)
		auth.Get("/role/all",controller.GetRoles)
		auth.Get("/role/:id",controller.GetRole)
	    auth.Get("/role/:id/groups",controller.GetRoleGroups)

		auth.Post("/group/create",controller.CreateGroup)
		auth.Post("/group/edit/:id",controller.EditGroup)
		auth.Get("/group/all",controller.GetGroups)
		auth.Get("/group/:id",controller.GetGroup)


		auth.Get("/permission/all",controller.GetAllPermissions)

}


func (App) Permissions() []user.Permission{
	return []user.Permission{
		{Title:"Access Users",CodeName:"access",Description:"Access list to view list of users"},
		{Title:"Create Users",CodeName:"create",Description:"Create new user"},
		{Title:"Edit users",CodeName:"edit",Description:"Edit user data"},
		{Title:"Remove users",CodeName:"remove",Description:"Remove user data"},
		{Title:"Login as user",CodeName:"loginas",Description:"Login as another user without credentials"},
		{Title:"Limit to subgroup",CodeName:"limited",Description:"User can only access/modify subgroup users"},
	}
}