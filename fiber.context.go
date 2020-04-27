package io

import (
	"fmt"
	"github.com/gofiber/fiber"
	"github.com/iesreza/io/lib/log"
	"github.com/iesreza/io/lib/text"
	"mime/multipart"
	"reflect"
	"time"
)

func (r *Request) Accepts(offers ...string) (offer string) {
	return r.Context.Accepts(offers...)
}

func (r *Request) AcceptsCharsets(offers ...string) (offer string) {
	return r.Context.AcceptsCharsets(offers...)
}

func (r *Request) AcceptsEncodings(offers ...string) (offer string) {
	return r.Context.AcceptsEncodings(offers...)
}

func (r *Request) AcceptsLanguages(offers ...string) (offer string) {
	return r.Context.AcceptsLanguages(offers...)
}

func (r *Request) Append(field string, values ...string) {
	r.Context.Append(field, values...)
}

func (r *Request) Attachment(name ...string) {
	r.Context.Attachment(name...)
}

func (r *Request) BaseURL() string {
	return r.Context.BaseURL()
}

func (r *Request) Body(key ...string) string {
	return r.Context.Body(key...)
}

func (r *Request) BodyParser(out interface{}) error {
	return r.Context.BodyParser(out)
}

func (r *Request) ClearCookie(key ...string) {
	r.Context.ClearCookie(key...)
}

func (r *Request) Cookie(cookie *fiber.Cookie) {
	r.Context.Cookie(cookie)
}

func (r *Request) Cookies(key ...string) (value string) {
	return r.Context.Cookies(key...)
}

func (r *Request) Download(file string, name ...string) {
	r.Context.Download(file, name...)
}

func (r *Request) Error() error {
	return r.Context.Error()
}

//@Override format
func (r *Request) Format(body interface{}) {
	var b string
	accept := r.Context.Accepts("html", "json")

	switch val := body.(type) {
	case string:
		b = val
	case []byte:
		b = string(val)
	default:
		b = fmt.Sprintf("%+v", val)
	}
	switch accept {
	case "html":
		r.Context.SendString(b)
	case "json":
		if err := r.Context.JSON(body); err != nil {
			log.Error("Format: error serializing json ", err)
		}
	default:
		r.Context.SendString(b)
	}
}

func (r *Request) FormFile(key string) (*multipart.FileHeader, error) {
	return r.Context.FormFile(key)
}

func (r *Request) FormValue(key string) (value string) {
	return r.Context.FormValue(key)
}

func (r *Request) Fresh() bool {
	return r.Context.Fresh()
}

func (r *Request) Get(key string) (value string) {
	return r.Context.Get(key)
}

func (r *Request) Hostname() string {
	return r.Context.Hostname()
}

func (r *Request) IP() string {
	return r.Context.IP()
}

func (r *Request) IPs() []string {
	return r.Context.IPs()
}

func (r *Request) Is(extension string) (match bool) {
	return r.Context.Is(extension)
}

func (r *Request) JSON(json interface{}) error {
	return r.Context.JSON(json)
}

func (r *Request) JSONP(json interface{}, callback ...string) error {
	return r.Context.JSONP(json, callback...)
}

func (r *Request) Links(link ...string) {
	r.Context.Links(link...)
}

func (r *Request) Locals(key string, value ...interface{}) (val interface{}) {
	return r.Context.Locals(key, value...)
}

func (r *Request) Location(path string) {
	r.Context.Location(path)
}

func (r *Request) Method(override ...string) string {
	return r.Context.Method(override...)
}

func (r *Request) MultipartForm() (*multipart.Form, error) {
	return r.Context.MultipartForm()
}

func (r *Request) Next(err ...error) {
	r.Context.Next(err...)
}

func (r *Request) OriginalURL() string {
	return r.Context.OriginalURL()
}

func (r *Request) Params(key string) (value string) {
	return r.Context.Params(key)
}

func (r *Request) Path(override ...string) string {
	return r.Context.Path(override...)
}

func (r *Request) Protocol() string {
	return r.Context.Protocol()
}

func (r *Request) Query(key string) (value string) {
	return r.Context.Query(key)
}

func (r *Request) Range(size int) (rangeData fiber.Range, err error) {
	return r.Context.Range(size)
}

func (r *Request) Redirect(path string, status ...int) {
	r.Context.Redirect(path, status...)
}

func (r *Request) Render(file string, bind interface{}) error {
	return r.Context.Render(file, bind)
}

func (r *Request) Route() *fiber.Route {
	return r.Context.Route()
}

func (r *Request) SaveFile(fileheader *multipart.FileHeader, path string) error {
	return r.Context.SaveFile(fileheader, path)
}

func (r *Request) Secure() bool {
	return r.Context.Secure()
}

func (r *Request) SendHTML(bodies ...interface{}) {
	r.Set("Content-Type", "text/html")
	r.Context.Send(bodies...)
}

func (r *Request) Send(bodies ...interface{}) {
	r.Context.Send(bodies...)
}

func (r *Request) SendBytes(body []byte) {
	r.Context.SendBytes(body)
}

func (r *Request) SendFile(file string, noCompression ...bool) {
	r.Context.SendFile(file, noCompression...)
}

func (r *Request) SendStatus(status int) {
	r.Context.SendStatus(status)
}

func (r *Request) SendString(body string) {
	r.Context.SendString(body)
}

func (r *Request) Set(key string, val string) {
	r.Context.Set(key, val)
}

func (r *Request) Subdomains(offset ...int) []string {
	return r.Context.Subdomains(offset...)
}

func (r *Request) Stale() bool {
	return r.Context.Stale()
}

func (r *Request) Status(status int) *Request {
	r.Context.Status(status)
	return r
}

func (r *Request) Type(ext string) *Request {
	r.Context.Type(ext)
	return r
}

func (r *Request) Vary(fields ...string) {
	r.Context.Vary(fields...)
}

func (r *Request) Write(bodies ...interface{}) {
	r.Context.Write(bodies...)
}

func (r *Request) XHR() bool {
	return r.Context.XHR()
}
func (r *Request) SetCookie(key string, val interface{}, params ...interface{}) {
	cookie := new(fiber.Cookie)
	cookie.Name = key

	ref := reflect.ValueOf(val)
	switch ref.Kind() {
	case reflect.String:
		cookie.Value = val.(string)
		break
	case reflect.Ptr:
		r.SetCookie(key, ref.Elem().Interface(), params...)
		return
	case reflect.Map, reflect.Struct, reflect.Array, reflect.Slice:
		cookie.Value = text.ToJSON(val)
		break
	default:
		cookie.Value = fmt.Sprint(val)
	}

	for _, item := range params {
		if v, ok := item.(time.Duration); ok {
			cookie.Expires = time.Now().Add(v)
		}
		if v, ok := item.(time.Time); ok {
			cookie.Expires = v
		}
	}
	r.Cookie(cookie)
}
