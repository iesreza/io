package io

import "github.com/iesreza/io/ui"

type Template interface{
	Name()       func()string
	Version()    func()string
	Pages()      func()map[string]ui.Page
	Elements()   func()map[string]ui.Element
}


type TemplateEngine struct {
	Path    string
}

func (t *TemplateEngine)Setup() {

}