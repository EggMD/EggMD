package form

import (
	"reflect"

	"github.com/unknwon/com"
)

// Assign 将表单中的用户输入的字段值再次插入回模板。
func Assign(form interface{}, data map[string]interface{}) {
	typ := reflect.TypeOf(form)
	val := reflect.ValueOf(form)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		fieldName := field.Tag.Get("form")
		// 允许在结构体中使用 `-` tag 来忽略字段。
		if fieldName == "-" {
			continue
		} else if len(fieldName) == 0 {
			fieldName = com.ToSnakeCase(field.Name)
		}

		data[fieldName] = val.Field(i).Interface()
	}
}
