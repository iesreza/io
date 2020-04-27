package settings

import "github.com/iesreza/io"

type Settings struct {
	io.Model
	Reference string
	Title     string
	Data      string
	Default   string
	Ptr       interface{} `gorm:"-"`
}
