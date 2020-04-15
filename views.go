package io

import (
	"fmt"
	"github.com/CloudyKit/jet"
)

//TODO: Concurrent View Pool
type views map[string]*jet.Set
var viewList = views{}

func RegisterView(prefix,path string) *jet.Set {
	viewList[prefix] = jet.NewHTMLSet(path)
	if config.Server.Debug {
		viewList[prefix].SetDevelopmentMode(true)
	}
	return viewList[prefix]
}

func GetView(prefix,name string) (*jet.Template,error) {
	if t, ok := viewList[prefix]; ok {
		return t.GetTemplate(name)
	}

	return nil,fmt.Errorf("template prefix \"%s\" not found",prefix)

}