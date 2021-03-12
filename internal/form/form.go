package form

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-macaron/binding"
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

type Form interface {
	binding.Validator
}

func getRuleBody(field reflect.StructField, prefix string) string {
	for _, rule := range strings.Split(field.Tag.Get("binding"), ";") {
		if strings.HasPrefix(rule, prefix) {
			return rule[len(prefix) : len(rule)-1]
		}
	}
	return ""
}

func getSize(field reflect.StructField) string {
	return getRuleBody(field, "Size(")
}

func getMinSize(field reflect.StructField) string {
	return getRuleBody(field, "MinSize(")
}

func getMaxSize(field reflect.StructField) string {
	return getRuleBody(field, "MaxSize(")
}

func getInclude(field reflect.StructField) string {
	return getRuleBody(field, "Include(")
}

func validate(errs binding.Errors, data map[string]interface{}, f Form) binding.Errors {
	if errs.Len() == 0 {
		return errs
	}

	data["HasError"] = true
	Assign(f, data)

	typ := reflect.TypeOf(f)
	val := reflect.ValueOf(f)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		fieldName := field.Tag.Get("form")
		if fieldName == "-" {
			continue
		}

		if errs[0].FieldNames[0] == field.Name {
			data["Err_"+field.Name] = true

			trName := field.Tag.Get("locale")

			switch errs[0].Classification {
			case binding.ERR_REQUIRED:
				data["ErrorMsg"] = trName + "不能为空。"
			case binding.ERR_ALPHA_DASH:
				data["ErrorMsg"] = trName + "必须为英文字母、阿拉伯数字或横线（-_）。"
			case binding.ERR_ALPHA_DASH_DOT:
				data["ErrorMsg"] = trName + "必须为英文字母、阿拉伯数字、横线（-_）或点。"
			case binding.ERR_SIZE:
				data["ErrorMsg"] = trName + fmt.Sprintf("长度必须为 %s。", getSize(field))
			case binding.ERR_MIN_SIZE:
				data["ErrorMsg"] = trName + fmt.Sprintf("长度最小为 %s 个字符。", getMinSize(field))
			case binding.ERR_MAX_SIZE:
				data["ErrorMsg"] = trName + fmt.Sprintf("长度最大为 %s 个字符。", getMaxSize(field))
			case binding.ERR_EMAIL:
				data["ErrorMsg"] = trName + "不是一个有效的电子邮箱地址。"
			case binding.ERR_URL:
				data["ErrorMsg"] = trName + "不是一个有效的 URL。"
			case binding.ERR_INCLUDE:
				data["ErrorMsg"] = trName + fmt.Sprintf("必须包含子字符串 '%s'。", getInclude(field))
			default:
				data["ErrorMsg"] = "未知错误： " + errs[0].Classification
			}
			return errs
		}
	}
	return errs
}
