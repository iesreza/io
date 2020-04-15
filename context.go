package io

import (
	"bytes"
	"github.com/CloudyKit/jet"
	"github.com/gofiber/fiber"
	"github.com/gofiber/session"
	"github.com/iesreza/io/lib/jwt"
	"github.com/iesreza/io/lib/log"
	"github.com/iesreza/io/user"
	"net/http"
	"strings"
	"time"
)


type Request struct {
	Variables   fiber.Map
	Session     *session.Store
	Context     *fiber.Ctx
	JWT         *jwt.Payload
	User        *user.User
}



func Upgrade(ctx *fiber.Ctx) *Request {
	r := Request{}
	r.Session = Sessions.Start(ctx)
	r.Context = ctx
	if r.Cookies("access_token") != ""{
		token, err := jwt.Verify(r.Cookies("access_token"))
		if err == nil{
			r.JWT = &token
		}else{
			r.Status(http.StatusUnauthorized)
			r.Send("invalid JWT token")
			log.Error(err)
		}
	}else{
		r.JWT = &jwt.Payload{ Empty:true, Data: map[string]interface{}{} }
	}
	r.User = getUser(&r)
	return &r
}

func getUser(request *Request) *user.User {
	// return user using jwt
	log.Warning("implement getUser")
	return &user.User{}
}

func (r *Request)Persist()  {
	r.Session.Save(r.Context,r.Session)
	if !r.JWT.Empty{
		exp := time.Now().Add(config.JWT.Age)
		if d,exist := r.JWT.Get("_extend_duration"); exist{
			duration := d.(time.Duration)
			exp = time.Now().Add(duration)
		}
		token,err := jwt.Generate(r.JWT.Data)
		if err == nil{
			r.Cookie(&fiber.Cookie{
				Name:"access_token",
				Value:token,
				Expires: exp,
			})

		}else{
			log.Error(err)
		}

	}
}

func (r *Request)View(data map[string]interface{},views... string)  {
	buff := r.RenderView(data,views...)
	buff.Bytes()
	r.SendHTML(buff.Bytes())
}

func (r *Request)RenderView(data map[string]interface{},views... string) *bytes.Buffer {
	var buff bytes.Buffer
	if data == nil{
		data = map[string]interface{}{}
	}
	vars := jet.VarMap{}
	vars.Set("base",r.Context.Protocol()+"://"+r.Context.Hostname())
	vars.Set("proto",r.Context.Protocol())
	vars.Set("hostname",r.Context.Hostname())
	for _,view := range views{
		buff = bytes.Buffer{}
		parts := strings.Split(view,".")

		if len(parts) > 1 {
			t,err := GetView(parts[0], strings.Join(parts[1:], "."))
			if err == nil{
				t.Execute(&buff,vars,data)
			}
			vars.Set("body",buff.Bytes())

		}
	}
	return &buff
}