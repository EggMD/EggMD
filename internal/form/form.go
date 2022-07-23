package form

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/flamego/flamego"
	"github.com/wuhan005/govalid"
	log "unknwon.dev/clog/v2"
)

type ErrorCategory string

const (
	ErrorCategoryDeserialization ErrorCategory = "deserialization"
	ErrorCategoryValidation      ErrorCategory = "validation"
)

type Error struct {
	Category ErrorCategory
	Error    error
}

func Bind(model interface{}) flamego.Handler {
	// Ensure not pointer.
	if reflect.TypeOf(model).Kind() == reflect.Ptr {
		panic("form: pointer can not be accepted as binding model")
	}

	return flamego.ContextInvoker(func(c flamego.Context) {
		obj := reflect.New(reflect.TypeOf(model))
		r := c.Request().Request
		if r.Body != nil {
			defer func() { _ = r.Body.Close() }()
			err := json.NewDecoder(r.Body).Decode(obj.Interface())
			if err != nil {
				c.Map(Error{Category: ErrorCategoryDeserialization, Error: err})
				if _, err := c.Invoke(errorHandler); err != nil {
					panic("form: " + err.Error())
				}
				return
			}
		}

		errors, ok := govalid.Check(obj.Interface())
		if !ok {
			c.Map(Error{Category: ErrorCategoryValidation, Error: errors[0]})
			if _, err := c.Invoke(errorHandler); err != nil {
				panic("form: " + err.Error())
			}
			return
		}

		// Validation passed.
		c.Map(obj.Elem().Interface())
	})
}

func errorHandler(c flamego.Context, error Error) {
	c.ResponseWriter().WriteHeader(http.StatusBadRequest)
	c.ResponseWriter().Header().Set("Content-Type", "application/json; charset=utf-8")

	var errorCode int
	var msg string
	if error.Category == ErrorCategoryDeserialization {
		errorCode = 40000
		msg = "请求体错误"
	} else {
		errorCode = 40001
		msg = error.Error.Error()
	}

	body := map[string]interface{}{
		"error": errorCode,
		"msg":   msg,
	}
	err := json.NewEncoder(c.ResponseWriter()).Encode(body)
	if err != nil {
		log.Error("Failed to encode response body: %v", err)
	}
}
