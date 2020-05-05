package test

import (
	"github.com/gofiber/fiber"
	"github.com/iesreza/io"
	"github.com/iesreza/io/html"
)

type Controller struct{}

type ColumnType int

const (
	TEXT   ColumnType = 0
	NUMBER ColumnType = 1
	DATE   ColumnType = 2
	HTML   ColumnType = 3
	RANGE  ColumnType = 4
	SELECT ColumnType = 5
	CUSTOM ColumnType = 6
)

type Column struct {
	Type          ColumnType
	Title         string
	Width         int
	Name          string
	Processor     func(data interface{})
	SearchBuilder func(r *io.Request) string
	Attribs       html.Attributes
}

type FilterView struct {
	Style   string
	Columns []Column
	Model   interface{}
	Join    string
}

func (fv FilterView) Render() string {
	return "abcd"
}

func (col Column) Filter(r *io.Request) string {
	if col.SearchBuilder != nil {
		return col.SearchBuilder(r)
	}
	var el *html.InputStruct
	switch col.Type {
	case NUMBER:
		el = html.Input("number", col.Name, "")
		break
	case DATE:
		el = html.Input("daterange", col.Name, "")
		break
	default:
		el = html.Input("text", col.Name, "")
	}
	if col.Attribs != nil {
		el.Attributes = col.Attribs
	}
	if col.Title != "" {
		el.Placeholder(col.Title)
	}
	if r.Query(col.Name) != "" {
		el.Value = r.Query(col.Name)
	} else {
		el.Value = ""
	}
	return el.Render()
}

func FilterViewController(ctx *fiber.Ctx) {
	r := io.Upgrade(ctx)
	fv := FilterView{
		Columns: []Column{
			{Type: TEXT, Title: "Name", Name: "name"},
			{Type: TEXT, Title: "Username", Name: "username"},
		},
		Model: MyModel{},
	}

	r.View(fv, "test.list")

}
