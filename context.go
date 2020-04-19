package io

import (
	"bytes"
	"github.com/CloudyKit/jet"
	"github.com/gofiber/fiber"
	"github.com/gofiber/session"
	"github.com/iesreza/io/errors"
	"github.com/iesreza/io/lib/jwt"
	"github.com/iesreza/io/lib/log"
	"github.com/iesreza/io/user"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type Request struct {
	Variables fiber.Map
	Session   *session.Store
	Context   *fiber.Ctx
	JWT       *jwt.Payload
	User      *user.User
	Response  Response
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Error   e.Errors    `json:"errors"`
	Data    interface{} `json:"data"`
}

func (response Response) HasError() bool {
	return response.Error.Exist()
}

func Upgrade(ctx *fiber.Ctx) *Request {
	r := Request{}
	r.Variables = fiber.Map{}
	r.Session = Sessions.Start(ctx)
	r.Context = ctx
	r.Response = Response{}
	r.Response.Error = e.Errors{}
	if r.Cookies("access_token") != "" {
		token, err := jwt.Verify(r.Cookies("access_token"))
		if err == nil {
			r.JWT = &token
		} else {
			r.Status(http.StatusUnauthorized)
			r.Send("invalid JWT token")
			log.Error(err)
		}
	} else {
		r.JWT = &jwt.Payload{Empty: true, Data: map[string]interface{}{}}
	}
	r.User = getUser(&r)
	return &r
}

func getUser(request *Request) *user.User {
	// return user using jwt
	log.Warning("implement getUser")
	return &user.User{}
}

func (r *Request) Persist() {
	r.Session.Save(r.Context, r.Session)
	if !r.JWT.Empty {
		exp := time.Now().Add(config.JWT.Age)
		if d, exist := r.JWT.Get("_extend_duration"); exist {
			duration := d.(time.Duration)
			exp = time.Now().Add(duration)
		}
		token, err := jwt.Generate(r.JWT.Data)
		if err == nil {
			r.Cookie(&fiber.Cookie{
				Name:    "access_token",
				Value:   token,
				Expires: exp,
			})

		} else {
			log.Error(err)
		}

	}
}

func (r *Request) View(data map[string]interface{}, views ...string) {
	buff := r.RenderView(data, views...)
	buff.Bytes()
	r.SendHTML(buff.Bytes())
}

func (r *Request) RenderView(data map[string]interface{}, views ...string) *bytes.Buffer {
	var buff bytes.Buffer
	if data == nil {
		data = map[string]interface{}{}
	}
	vars := jet.VarMap{}
	vars.Set("base", r.Context.Protocol()+"://"+r.Context.Hostname())
	vars.Set("proto", r.Context.Protocol())
	vars.Set("hostname", r.Context.Hostname())
	for _, view := range views {
		buff = bytes.Buffer{}
		parts := strings.Split(view, ".")

		if len(parts) > 1 {
			t, err := GetView(parts[0], strings.Join(parts[1:], "."))
			if err == nil {
				t.Execute(&buff, vars, data)
			}
			vars.Set("body", buff.Bytes())

		}
	}
	return &buff
}

func (r *Request) WriteResponse(resp ...interface{}) {
	if len(resp) == 0 {
		r._writeResponse(r.Response)
		return
	}

	var message = false
	for _, item := range resp {
		ref := reflect.ValueOf(item)

		switch ref.Kind() {
		case reflect.Struct:
			if v, ok := item.(Response); ok {
				r._writeResponse(v)
				return
			} else if v, ok := item.(e.Error); ok {
				r.Response.Error.Push(&v)
				r._writeResponse(r.Response)
			} else {
				r.Response.Data = item
			}

		case reflect.Ptr:
			obj := ref.Elem().Interface()
			if v, ok := obj.(Response); ok {
				r._writeResponse(v)
				return
			} else if v, ok := obj.(e.Error); ok {
				r.Response.Error.Push(&v)
			} else if v, ok := item.(error); ok {

				r.Response.Error.Push(e.Context(v.Error()))
			} else {
				r.Response.Data = obj
			}

			break
		case reflect.Bool:
			r.Response.Success = item.(bool)
			break
		case reflect.String:
			if !message {
				r.Response.Message = item.(string)
				message = true
			} else {
				r.Response.Data = item
			}
			break
		default:
			r.Response.Data = item
		}

	}
	r._writeResponse(r.Response)

}

func (r *Request) _writeResponse(resp Response) {
	if resp.HasError() {
		r.Response.Success = false
	} else {
		r.Response.Success = true
	}
	r.JSON(resp)
}

func (r *Request) SetError(err interface{}) {
	if v, ok := err.(error); ok {
		r.Response.Error.Push(e.Context(v.Error()))
		return
	}
	if v, ok := err.(e.Error); ok {
		r.Response.Error.Push(&v)
		return
	}
	if v, ok := err.(*e.Error); ok {
		r.Response.Error.Push(v)
		return
	}
	log.Error("invalid error provided %+v", err)

}

func (r *Request) Throw(e *e.Error) {
	r.Response.Error.Push(e)
	r.WriteResponse()
}

func (r *Request) HasError() bool {
	return r.Response.Error.Exist()
}
