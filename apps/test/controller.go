package test

import (
	"github.com/gofiber/fiber"
	"github.com/iesreza/io"
	"github.com/iesreza/io/html"
	"strings"
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
	Type         ColumnType
	Title        string
	Width        int
	Name         string
	Options      []html.KeyValue
	InputBuilder func(r *io.Request) *html.InputStruct
	Attribs      html.Attributes
	QueryBuilder func(r *io.Request) []string
	SimpleFilter string
}

type FilterView struct {
	Style        string
	Columns      []Column
	Model        interface{}
	Join         string
	Attribs      html.Attributes
	Unscoped     bool
	QueryBuilder func(r *io.Request) []string
	data         *map[string]interface{}
}

func (fv FilterView) Render() string {
	return "abcd"
}

func (fv FilterView) GetData(r *io.Request) {

}

func (fv FilterView) Prepare(r *io.Request) {
	var query []string
	if fv.QueryBuilder != nil {
		query = append(query, fv.QueryBuilder(r)...)
	}
	for _, column := range fv.Columns {
		if column.QueryBuilder != nil {
			query = append(query, column.QueryBuilder(r)...)
		} else {
			v := r.Body(column.Name)
			if v != "" {
				query = append(query, strings.Replace(column.SimpleFilter, "*", v, -1))
			}
		}
	}

	db := io.GetDBO()
	if fv.Unscoped {
		db = db.Unscoped()
	}
	db.Where(strings.Join(query, " AND ")).Find(fv.data)

}

func (col Column) Filter(r *io.Request) *html.InputStruct {
	if col.InputBuilder != nil {
		return col.InputBuilder(r)
	}
	var el *html.InputStruct
	switch col.Type {
	case NUMBER:
		el = html.Input("number", col.Name, "")
		break
	case DATE:
		el = html.Input("daterange", col.Name, "")
		break
	case SELECT:
		el = html.Input("select", col.Name, "").SetOptions(col.Options)
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
	return el
}

func FilterViewController(ctx *fiber.Ctx) {
	r := io.Upgrade(ctx)
	fv := FilterView{
		Columns: []Column{
			{Type: TEXT, Title: "Name", Name: "name", SimpleFilter: "name = '%*%'"},
			{Type: TEXT, Title: "Username", Name: "username", SimpleFilter: "username = '%*%'"},
			{Type: SELECT, Title: "Group", Name: "group", Options: []html.KeyValue{
				{1, "Admin"},
				{2, "Non Admin"},
			}, SimpleFilter: "group = '*'"},
		},
		Model: MyModel{},
	}

	fv.Prepare(r)

	r.View(fv, "test.list")

}
