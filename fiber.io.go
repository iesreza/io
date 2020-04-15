package io

import (
	"github.com/gofiber/fiber"
	"net/http"
)

func  Group(prefix string, handlers ...func(*fiber.Ctx)) *fiber.Group {
	if app == nil{
		panic("Access object before call Setup()")
	}
	return app.Group(prefix, handlers ...)
}

func Static(prefix,path string)  {
	statics = append(statics,[2]string{ prefix,path })
}

func  Use(args ...interface{}) *fiber.Fiber {
	if app == nil{
		panic("Access object before call Setup()")
	}
	return app.Use(args ...)
}

func  Connect(path string, handlers... func(*fiber.Ctx)) *fiber.Fiber {
	if app == nil{
		panic("Access object before call Setup()")
	}
	return app.Connect(path, handlers ...)
}

func  Put(path string, handlers ...func(*fiber.Ctx)) *fiber.Fiber {
	if app == nil{
		panic("Access object before call Setup()")
	}
	return app.Put(path, handlers ...)
}

func  Post(path string, handlers ...func(*fiber.Ctx)) *fiber.Fiber {
	if app == nil{
		panic("Access object before call Setup()")
	}
	return app.Post(path, handlers ...)
}

func  Delete(path string, handlers ...func(*fiber.Ctx)) *fiber.Fiber {
	if app == nil{
		panic("Access object before call Setup()")
	}
	return app.Delete(path, handlers ...)
}

func  Head(path string, handlers ...func(*fiber.Ctx)) *fiber.Fiber {
	if app == nil{
		panic("Access object before call Setup()")
	}
	return app.Head(path, handlers ...)
}

func  Patch(path string, handlers ...func(*fiber.Ctx)) *fiber.Fiber {
	if app == nil{
		panic("Access object before call Setup()")
	}
	return app.Patch(path, handlers ...)
}

func  Options(path string, handlers ...func(*fiber.Ctx)) *fiber.Fiber {
	if app == nil{
		panic("Access object before call Setup()")
	}
	return app.Options(path, handlers ...)
}

func  Trace(path string, handlers ...func(*fiber.Ctx)) *fiber.Fiber {
	if app == nil{
		panic("Access object before call Setup()")
	}
	return app.Trace(path, handlers ...)
}

func  Get(path string, handlers ...func(*fiber.Ctx)) *fiber.Fiber {
	if app == nil{
		panic("Access object before call Setup()")
	}
	return app.Get(path, handlers ...)
}

func  All(path string, handlers ...func(*fiber.Ctx)) *fiber.Fiber {
	if app == nil{
		panic("Access object before call Setup()")
	}
	return app.All(path, handlers ...)
}


func  Shutdown() error {
	if app == nil{
		panic("Access object before call Setup()")
	}
	return app.Shutdown()
}

func  Test(request *http.Request, msTimeout ...int) (*http.Response, error) {
	if app == nil{
		panic("Access object before call Setup()")
	}
	return app.Test(request, msTimeout ...)
}


