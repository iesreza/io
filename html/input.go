package html

import (
	"bytes"
	"fmt"
	"github.com/CloudyKit/jet"
	"github.com/iesreza/gutil/log"
	"github.com/iesreza/io"
	"strconv"
	"strings"
)

type KeyValue struct {
	Key   interface{}
	Value interface{}
}

type InputStruct struct {
	Type       string
	Label      string
	Name       string
	Pre        string
	Post       string
	Hint       string
	Size       int
	LabelSize  int
	InputSize  int
	Horizontal bool
	Attributes Attributes
	Options    []KeyValue
	Sub        string
	Value      interface{}
}

func Input(_type, name, label string) *InputStruct {
	return &InputStruct{
		Type:       strings.ToLower(_type),
		Attributes: Attributes{},
		Label:      label,
		Name:       name,
	}
}
func (i InputStruct) IsSelected(inputs ...interface{}) bool {

	if v, ok := i.Value.([]string); ok {
		for _, input := range inputs {
			for _, item := range v {
				if fmt.Sprint(input) == item {
					return true
				}
			}
		}
	} else if v, ok := i.Value.(string); ok {
		for _, input := range inputs {
			if fmt.Sprint(input) == v {
				return true
			}
		}
	} else if v, ok := i.Value.(int); ok {

		for _, input := range inputs {
			if fmt.Sprint(input) == strconv.Itoa(v) {
				return true
			}
		}
	}

	return false
}

func (i *InputStruct) SetOptions(options []KeyValue) *InputStruct {
	i.Options = options
	return i
}

func (i *InputStruct) SetLabel(s string) *InputStruct {
	i.Label = s
	return i
}

func (i *InputStruct) SetLabelSize(s int) *InputStruct {
	i.LabelSize = s
	i.InputSize = 12 - s
	return i
}

func (i *InputStruct) SetInputSize(s int) *InputStruct {
	i.InputSize = s
	i.LabelSize = 12 - s
	return i
}

func (i *InputStruct) SetName(s string) *InputStruct {
	i.Name = s
	return i
}

func (i *InputStruct) SetSub(s string) *InputStruct {
	i.Sub = s
	return i
}

func (i *InputStruct) SetValue(s interface{}) *InputStruct {
	i.Value = s
	return i
}

func (i *InputStruct) SetAttr(k string, s interface{}) *InputStruct {
	i.Attributes[k] = s
	return i
}

func (i *InputStruct) Placeholder(s string) *InputStruct {
	i.Attributes["placeholder"] = s
	return i
}

func (i *InputStruct) ID(s string) *InputStruct {
	i.Attributes["id"] = s
	return i
}

func (i *InputStruct) Max(s interface{}) *InputStruct {
	i.Attributes["max"] = s
	return i
}

func (i *InputStruct) Min(s interface{}) *InputStruct {
	i.Attributes["min"] = s
	return i
}

func (i *InputStruct) MinLength(s string) *InputStruct {
	i.Attributes["min-length"] = s
	return i
}
func (i *InputStruct) MaxLength(s string) *InputStruct {
	i.Attributes["max-length"] = s
	return i
}

func (i *InputStruct) Class(s string) *InputStruct {
	i.Attributes["class"] = s
	return i
}

func (i *InputStruct) AddClass(s string) *InputStruct {
	if v, ok := i.Attributes["class"]; ok {
		i.Attributes["class"] = v.(string) + " " + strings.TrimSpace(s)
	} else {
		i.Attributes["class"] = strings.TrimSpace(s)
	}

	return i
}

func (i *InputStruct) Required(s string) *InputStruct {
	i.Attributes["required"] = true
	return i
}

func (i *InputStruct) PreText(s string) *InputStruct {
	i.Pre = s
	return i
}

func (i *InputStruct) PostText(s string) *InputStruct {
	i.Post = s
	return i
}

func (i *InputStruct) Disabled() *InputStruct {
	i.Attributes["disabled"] = true
	return i
}

func (i *InputStruct) Readonly(s int) *InputStruct {
	i.Attributes["readonly"] = true
	return i
}

func (i *InputStruct) Depend(s string) *InputStruct {
	i.Attributes["depend"] = s
	return i
}

func (i *InputStruct) SetSize(s int) *InputStruct {
	i.Size = s
	return i
}

func (i *InputStruct) Multiple() *InputStruct {
	i.Attributes["multiple"] = true
	return i
}

func (i InputStruct) String() string {
	return i.Render()
}

func (i InputStruct) Render() string {
	if i.Attr("id") == "" {
		i.ID(i.Name)
	}
	if i.LabelSize+i.InputSize == 12 {
		i.Horizontal = true
	}
	i.AddClass("form-control")
	buff := bytes.Buffer{}

	t, err := io.GetView(ViewKey, i.Type)
	if err == nil {
		vars := jet.VarMap{}
		vars.Set("input", i)
		t.Execute(&buff, vars, "")
	} else {
		log.Error(err)
		log.Error(ViewKey + "." + i.Type)
	}
	return string(buff.Bytes())
}

func (i InputStruct) Attr(key string) string {
	if v, ok := i.Attributes[key]; ok {
		return strings.Replace(fmt.Sprint(v), "\"", "\\\"", -1)
	}
	return ""
}
