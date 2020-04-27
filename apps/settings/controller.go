package settings

import (
	"encoding/json"
	"github.com/gofiber/fiber"
	"github.com/iesreza/io"
	"github.com/iesreza/io/errors"
	"github.com/iesreza/io/html"
	"github.com/iesreza/io/lib"
	"github.com/iesreza/io/lib/text"
	"reflect"
	"strings"
)

type Controller struct{}
type Registry struct {
	Title     string
	Form      []html.InputStruct
	Slug      string
	Reference string
}

func (c Controller) set(s string, object interface{}) {
	ref := reflect.ValueOf(object)
	_type := strings.ToLower(ref.Elem().Type().String())

	_default := text.ToJSON(object)
	obj := Settings{
		Reference: _type,
		Title:     s,
		Data:      _default,
		Default:   _default,
		Ptr:       object,
	}
	settings.Set(_type, obj)
	if db.Debug().Where("reference = ?", _type).Take(&obj).RecordNotFound() {
		db.Debug().Create(&obj)
		return
	}
	json.Unmarshal([]byte(obj.Data), object)
}

func (c Controller) view(ctx *fiber.Ctx) {
	r := io.Upgrade(ctx)
	r.Var("heading", "Settings")

	var registries []Registry
	for _, key := range settings.Keys() {
		v := settings.Get(key).(Settings)

		registries = append(registries, Registry{
			v.Title, c.getForm(v.Ptr), text.Slugify(v.Reference), v.Reference,
		})
	}

	r.View(map[string]interface{}{
		"registries": registries,
	}, "settings.settings", "template.default")

	/*	r.View(map[string]interface{}{
		"elements": []html.InputStruct{*s1,*s2,*s3,*s4,*D1,*D5,*D2,*D3,*D4,*elem8,*elem9,*elem10,*elem1,*elem2,*elem3},
	},"registry.settings","template.default")*/
}

func (c Controller) getForm(v interface{}) []html.InputStruct {
	ref := reflect.ValueOf(v).Elem()
	frm := []html.InputStruct{}
	for i := 0; i < ref.NumField(); i++ {
		field := ref.Type().Field(i)
		tags := field.Tag
		if len(tags) == 0 {
			continue
		}
		input := html.Input("text", field.Name, field.Name)

		input.SetValue(ref.Field(i).Interface())

		if typ, ok := tags.Lookup("type"); ok {
			input.Type = typ
		}
		if hint, ok := tags.Lookup("hint"); ok {
			input.Hint = hint
		}
		if label, ok := tags.Lookup("label"); ok {
			input.SetLabel(label)
		}
		if col, ok := tags.Lookup("col"); ok {
			input.SetSize(lib.ParseSafeInt(col))
		}
		if min, ok := tags.Lookup("min"); ok {
			input.Min(min)
		}
		if max, ok := tags.Lookup("max"); ok {
			input.Max(max)
		}
		if col, ok := tags.Lookup("options"); ok {
			chunks := strings.Split(col, ",")
			var options []html.KeyValue
			for _, option := range chunks {
				parts := strings.Split(option, ":")
				if len(parts) == 2 {
					options = append(options, html.KeyValue{parts[0], parts[1]})
				}
			}
			input.SetOptions(options)
		}
		frm = append(frm, *input)
	}
	return frm
}

func (c Controller) save(ctx *fiber.Ctx) {
	r := io.Upgrade(ctx)

	name := r.Params("name")
	if !settings.Has(name) {
		r.WriteResponse(false, e.Context("settings not found"))
		return
	}
	item := settings.Get(name).(Settings)

	err := r.BodyParser(item.Ptr)
	if err != nil {
		r.WriteResponse(false, err)
		return
	}
	b, err := json.Marshal(item.Ptr)
	item.Data = string(b)

	db.Debug().Model(&item).Where("reference = ?", item.Reference).Update("data", item.Data)
	settings.Set(name, item)
	r.WriteResponse(true, item)
}

func (c Controller) reset(ctx *fiber.Ctx) {
	r := io.Upgrade(ctx)

	name := r.Params("name")
	if !settings.Has(name) {
		r.WriteResponse(false, e.Context("settings not found"))
		return
	}
	item := settings.Get(name).(Settings)

	item.Data = item.Default

	err := json.Unmarshal([]byte(item.Default), item.Ptr)
	if err != nil {
		r.WriteResponse(false, err)
		return
	}
	db.Debug().Model(&item).Where("reference = ?", item.Reference).Update("data", item.Data)
	settings.Set(name, item)
	r.WriteResponse(true, item)
}
