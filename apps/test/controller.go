package test

import (
	"fmt"
	"github.com/gofiber/fiber"
	"github.com/iesreza/io"
	"github.com/iesreza/io/html"
	"github.com/iesreza/io/lib/fontawesome"
	"reflect"
	"strings"
)

type Controller struct{}

type ColumnType int
type ActionType string

const (
	RESET    ActionType = "reset"
	SEARCH   ActionType = "search"
	PAGESIZE ActionType = "pagesize"
	ORDER    ActionType = "order"
)

const (
	TEXT    ColumnType = 0
	NUMBER  ColumnType = 1
	DATE    ColumnType = 2
	HTML    ColumnType = 3
	RANGE   ColumnType = 4
	SELECT  ColumnType = 5
	CUSTOM  ColumnType = 6
	ACTIONS ColumnType = 7
)

type Join struct {
	Model  interface{}
	MainFK string
	DestFK string
}

type Column struct {
	Type         ColumnType
	Title        string
	Width        int
	Resize       bool
	Order        bool
	Name         string
	Alias        string
	Options      []html.KeyValue
	InputBuilder func(r *io.Request) *html.InputStruct
	Attribs      html.Attributes
	QueryBuilder func(r *io.Request) []string
	SimpleFilter string
	Processor    func(column Column, data map[string]interface{}, r *io.Request) string
	Model        interface{}
}

type FilterView struct {
	Style        string
	Columns      []Column
	Model        interface{}
	Join         []Join
	Attribs      html.Attributes
	Unscoped     bool
	QueryBuilder func(r *io.Request) []string
	data         []map[string]interface{}
}

func (fv FilterView) GetData() []map[string]interface{} {
	return fv.data
}

func getName(t reflect.Type) string {
	parts := strings.Split(t.Name(), ".")
	return parts[len(parts)-1]
}

func quote(s string) string {
	return "\"" + s + "\""
}

func actionProcessor(column Column, data map[string]interface{}, r *io.Request) {

}

func defaultProcessor(column Column, data map[string]interface{}, r *io.Request) string {
	if column.Alias == "" {
		if v, ok := data[column.Name]; ok {
			return fmt.Sprint(v)
		}
	} else {
		if v, ok := data[column.Alias]; ok {
			return fmt.Sprint(v)
		}
	}
	return ""
}

func (fv *FilterView) Prepare(r *io.Request) {
	var query = []string{"true"}
	var _select []string
	var _join string
	var models = map[string]string{}
	var tables []string

	tables = append(tables, db.NewScope(fv.Model).TableName())
	for _, join := range fv.Join {
		t := db.NewScope(join.Model).TableName()
		models[getName(reflect.TypeOf(join.Model))] = t
		tables = append(tables, t)
		_join += " INNER JOIN " + quote(t) + " ON " + quote(tables[0]) + "." + quote(join.MainFK) + " = " + quote(t) + "." + quote(join.DestFK)
	}

	if fv.QueryBuilder != nil {
		query = append(query, fv.QueryBuilder(r)...)
	}
	for k, column := range fv.Columns {
		if column.Model == nil {
			column.Model = quote(tables[0])
		} else {
			if _, ok := column.Model.(string); !ok {
				column.Model = quote(db.NewScope(column.Model).TableName())
			}
		}
		if column.Alias == "" {
			column.Alias = column.Name
		}
		if column.Processor == nil {
			fv.Columns[k].Processor = defaultProcessor
		}

		if column.Name != "" {
			fmt.Println(column.Model.(string))
			_select = append(_select, column.Model.(string)+"."+quote(column.Name)+" AS "+quote(column.Alias))
		}

		if column.QueryBuilder != nil {
			query = append(query, column.QueryBuilder(r)...)
		} else {
			v := r.FormValue(column.Name)
			if v != "" {
				query = append(query, strings.Replace(column.SimpleFilter, "*", v, -1))
			}
		}
	}

	db := io.GetDBO()
	if fv.Unscoped {
		db = db.Unscoped()
	}

	q := fmt.Sprintf("SELECT %s FROM %s %s WHERE %s",
		strings.Join(_select, ","),
		quote(tables[0]), //main table
		_join,
		strings.Join(query, " AND "),
	)
	fmt.Println(q)
	rows, err := db.Raw(q).Rows()
	if err != nil {
		return
	}
	columns, err := rows.Columns()
	if err != nil {
		return
	}
	length := len(columns)
	fv.data = make([]map[string]interface{}, 0)
	for rows.Next() {

		current := makeResultReceiver(length)
		if err := rows.Scan(current...); err != nil {
			panic(err)
		}
		value := make(map[string]interface{})
		for i := 0; i < length; i++ {
			k := columns[i]
			v := reflect.ValueOf(current[i]).Elem().Interface()
			value[k] = v
		}
		fv.data = append(fv.data, value)
	}

	fmt.Println(fv.data)
}

func makeResultReceiver(length int) []interface{} {
	result := make([]interface{}, 0, length)
	for i := 0; i < length; i++ {
		var current interface{}
		current = struct{}{}
		result = append(result, &current)
	}
	return result
}

func (col Column) Filter(r *io.Request) html.Renderable {
	if col.InputBuilder != nil {
		return col.InputBuilder(r)
	}

	var el *html.InputStruct
	switch col.Type {
	case NUMBER:
		el = html.Input("number", col.Name, "")
		el.SetAttr("onpressenter", "fv.filter(this)")
		break
	case DATE:
		el = html.Input("daterange", col.Name, "")
		el.SetAttr("onpressenter", "fv.filter(this)")
		break
	case SELECT:
		el = html.Input("select", col.Name, "").SetOptions(col.Options)
		el.SetAttr("onchange", "fv.filter(this)")
	case ACTIONS:
		var actions []html.Element
		for _, item := range col.Options {
			action := item.Key.(ActionType)
			switch action {
			case SEARCH:
				actions = append(actions, *html.Tag("button", item.Value).Set("class", "btn fv-action-btn").Set("onclick", "fv.filter(this)"))
				break
			case RESET:
				actions = append(actions, *html.Tag("button", item.Value).Set("class", "btn fv-action-btn").Set("onclick", "fv.reset(this)"))
				break
			case PAGESIZE:
				actions = append(actions, *html.Tag("div", html.Input("select", "size", fmt.Sprint(item.Value)).SetOptions([]html.KeyValue{
					{10, "10"},
					{25, "25"},
					{50, "50"},
					{100, "100"},
				})).Set("class", "btn fv-action-pagesize").Set("onclick", "fv.setSize(this)"))
				break
			}

		}
		return html.Tag("div", actions).Set("class", "fv-actions")
	default:
		el = html.Input("text", col.Name, "")
		el.SetAttr("onpressenter", "fv.filter(this)")
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
		Join: []Join{
			{MyGroup{}, "group", "id"},
		},
		Columns: []Column{
			{Type: TEXT, Title: "Name", Name: "name", SimpleFilter: "name = '%*%'"},
			{Type: TEXT, Title: "Username", Name: "username", SimpleFilter: "username = '%*%'"},
			{Type: SELECT, Title: "Group", Model: MyGroup{}, Name: "name", Alias: "group", Options: []html.KeyValue{
				{1, "Admin"},
				{2, "Non Admin"},
			}, SimpleFilter: "group = '*'",
				Processor: func(column Column, data map[string]interface{}, r *io.Request) string {

					return html.Tag("a", data["group"]).Set("href", "#").Render()
				},
			},

			{Type: ACTIONS, Title: "Actions", Options: []html.KeyValue{
				{SEARCH, html.Icon(fontawesome.Search)},
				{RESET, html.Icon(fontawesome.Undo)},
				{PAGESIZE, "Size:"},
			},
			},
		},
		Model: MyModel{},
	}

	fv.Prepare(r)
	r.View(fv, "test.list", "template.default")

}
