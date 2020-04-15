package main

import (
	"github.com/iesreza/io"
	"github.com/iesreza/io/apps/admin"
	"github.com/iesreza/io/apps/adminlte"
	"github.com/iesreza/io/apps/auth"
	"github.com/iesreza/io/apps/query"
)


func main()  {

	io.Setup()
	adminlte.Register()
	admin.Register()
	auth.Register()
	query.Register()
	io.Run()
}
