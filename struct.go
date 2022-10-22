package libtools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// Struct2Form 结构体转表单
func Struct2Form(obj interface{}, encode bool) (form string) {
	rt := reflect.TypeOf(obj)
	rv := reflect.ValueOf(obj)

	var formBox []string

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		var keyTag = f.Tag.Get("json")
		if keyTag == "" {
			keyTag = f.Name
		}

		if keyTag == "-" {
			continue
		}

		keyBox := strings.Split(keyTag, ",")
		key := keyBox[0]
		var value = Stringify(rv.FieldByName(f.Name).Interface())

		if len(keyBox) >= 2 {
			if keyBox[1] == "omitempty" && (value == "" || value == "0" || value == "false") {
				continue
			}
		}

		if encode {
			value = UrlEncode(value)
		}

		formBox = append(formBox, fmt.Sprintf(`%s=%s`, key, value))
	}

	form = strings.Join(formBox, "&")

	return
}

func Struct2FormMap(obj interface{}, encode bool) map[string]string {
	rt := reflect.TypeOf(obj)
	rv := reflect.ValueOf(obj)

	var formMap = make(map[string]string)

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		var keyTag = f.Tag.Get("form")
		if keyTag == "" {
			keyTag = f.Name
		}

		if keyTag == "-" {
			continue
		}

		keyBox := strings.Split(keyTag, ",")
		key := keyBox[0]
		var value = Stringify(rv.FieldByName(f.Name).Interface())

		if len(keyBox) >= 2 {
			if keyBox[1] == "omitempty" && (value == "" || value == "0" || value == "false") {
				continue
			}
		}

		if encode {
			value = UrlEncode(value)
		}

		formMap[key] = value
	}

	return formMap
}

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}

	return data
}

func Struct2MapV2(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		key := f.Tag.Get("json")
		data[key] = v.Field(i).Interface()
	}

	return data
}

func Struct2MapString(obj interface{}) map[string]string {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		key := f.Tag.Get("json")
		data[key] = v.Field(i).String()
	}

	return data
}

// Struct2MapV3 结构体转map,有效防止了整型意外转科学记数法
func Struct2MapV3(obj interface{}) map[string]interface{} {
	var data = make(map[string]interface{})
	bson, _ := json.Marshal(obj)
	d := json.NewDecoder(bytes.NewReader(bson))
	d.UseNumber()
	_ = d.Decode(&data)

	return data
}

func Map2struct(data map[string]interface{}, result interface{}) error {
	str, _ := json.Marshal(data)
	return json.Unmarshal(str, result)
}

// IsStructContainsField 判断结构体是否包含给定的字段
// 方法有点脆弱,不能传结构体指针!!!
func IsStructContainsField(obj interface{}, field string) bool {
	var has = false
	rt := reflect.TypeOf(obj)

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if f.Name == field {
			has = true
			break
		}

		if SnakeString(f.Name) == field {
			has = true
			break
		}

		var keyTag = f.Tag.Get("orm")
		if keyTag == "-" {
			continue
		}

		var needBreak = false
		keyBox := strings.Split(keyTag, ";")
		for _, key := range keyBox {
			if !strings.Contains(key, "column") {
				continue
			}

			key = StrReplace(key, []string{"column", "(", ")"}, "")
			if key == field {
				has = true
				needBreak = true
			}
		}

		if needBreak {
			break
		}
	}

	return has
}

func Map2Form(param map[string]interface{}, encode bool) (form string) {
	var formBox []string

	for k, v := range param {
		if encode {
			v = UrlEncode(Stringify(v))
		}
		formBox = append(formBox, fmt.Sprintf(`%s=%s`, k, Stringify(v)))
	}

	form = strings.Join(formBox, "&")
	return
}
