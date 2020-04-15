package __PACKAGE__
// Package      {{ .name }}
// Generated on {{ .date }}
// Created   by {{ .user }}


import "github.com/iesreza/io"


func Register()  {
	io.Register(App{})
}

type App struct {}

func (App) Register() {
	//implement me
}

func (App) Router() {
	//implement me
}
