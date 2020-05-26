package test

import (
	"github.com/iesreza/io"
)

type MyModel struct {
	io.Model
	Name     string
	Username string
	Group    int
	Type     int
	Alias    string
}

type MyGroup struct {
	io.Model
	Name string
}
