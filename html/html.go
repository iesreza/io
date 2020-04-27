package html

import (
	"fmt"
	"strings"
)

type Attributes map[string]interface{}

var ViewKey = "html"

func (attrs Attributes) Render() string {
	var res = ""
	for k, v := range attrs {
		res += " " + k + "=\"" + strings.Replace(fmt.Sprint(v), "\"", "\\\"", -1) + "\""
	}
	return res
}
