package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"
)

var zeros = []interface{}{
	int(0),
	int8(0),
	int16(0),
	int32(0),
	int64(0),
	uint(0),
	uint8(0),
	uint16(0),
	uint32(0),
	uint64(0),
	float32(0),
	float64(0),
}

func Empty(value interface{}) bool {
	if value == nil {
		return true
	}

	switch v := value.(type) {
	case bool:
		return !v
	case string:
		return v == ""
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		for _, zero := range zeros {
			if v == zero {
				return true
			}
		}
	case time.Time:
		return v.IsZero()
	}

	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			if !Empty(val.FieldByIndex([]int{i}).Interface()) {
				return false
			}
		}
		return true
	case reflect.Map, reflect.Slice, reflect.Chan:
		return (val.Len() == 0)
	case reflect.Ptr:
		return Empty(val.Elem().Interface())
	}
	return false
}

func NotEmpty(value interface{}) bool {
	return !Empty(value)
}

func Range(value, min, max int) bool {
	return min <= value && value <= max
}

func StringSize(value string, min, max int) bool {
	l := len([]rune(value))
	return min <= l && l <= max
}

func Regexp(value string, pattern interface{}) bool {
	var r *regexp.Regexp
	if str, ok := pattern.(string); ok {
		r = regexp.MustCompile(str)
	} else if rpattern, ok := pattern.(*regexp.Regexp); ok {
		r = rpattern
	}
	return r.MatchString(value)
}

func Equal(value, expected interface{}) bool {
	if value == nil && expected == nil {
		return true
	}
	if reflect.DeepEqual(value, expected) {
		return true
	}
	return false
}

func Contain(value, expected interface{}) bool {
	v := reflect.ValueOf(expected)
	if v.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			if Equal(value, v.Index(i).Interface()) {
				return true
			}
		}
	}
	return false
}

func TimeRange(value, from, to time.Time) bool {
	return from.UnixNano() <= value.UnixNano() && value.UnixNano() <= to.UnixNano()
}

var Messages = struct {
	NotEmpty, Range, StringSize, Regexp, Equal, Contain, TimeRange string
}{
	NotEmpty:   "can't be blank",
	Range:      "must be between %v and %v",
	StringSize: "string length must be between %v and %v",
	Regexp:     "must match with pattern\"%v\"",
	Equal:      "must be %v",
	Contain:    "must be one of following values. %v",
	TimeRange:  "must be between %v and %v",
}

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) HasErrors() bool {
	return len(v.Errors) > 0
}

func (v *Validator) AddError(key string, msg ...string) {
	v.Errors[key] = strings.Join(msg, " ")
}

func (v *Validator) SetError(result bool, key string, msg ...string) {
	if !result {
		v.AddError(key, msg...)
	}
}

func (v *Validator) NotEmpty(value interface{}, key string, msg ...string) bool {
	var m string
	if len(msg) == 0 {
		m = fmt.Sprintf(Messages.NotEmpty)
	} else {
		m = fmt.Sprintf(msg[0])
	}
	result := NotEmpty(value)
	v.SetError(result, key, m)
	return result
}

func (v *Validator) Range(value, min, max int, key string, msg ...string) bool {
	var m string
	if len(msg) == 0 {
		m = fmt.Sprintf(Messages.Range, min, max)
	} else {
		m = fmt.Sprintf(msg[0], min, max)
	}
	result := Range(value, min, max)
	v.SetError(result, key, m)
	return result
}

func (v *Validator) StringSize(value string, min, max int, key string, msg ...string) bool {
	var m string
	if len(msg) == 0 {
		m = fmt.Sprintf(Messages.StringSize, min, max)
	} else {
		m = fmt.Sprintf(msg[0], min, max)
	}
	result := StringSize(value, min, max)
	v.SetError(result, key, m)
	return result
}

func (v *Validator) Regexp(value string, pattern interface{}, key string, msg ...string) bool {
	var m string
	if len(msg) == 0 {
		m = fmt.Sprintf(Messages.Regexp, pattern)
	} else {
		m = fmt.Sprintf(msg[0], pattern)
	}
	result := Regexp(value, pattern)
	v.SetError(result, key, m)
	return result
}

func (v *Validator) Equal(value, expected interface{}, key string, msg ...string) bool {
	var m string
	if len(msg) == 0 {
		m = fmt.Sprintf(Messages.Equal, expected)
	} else {
		m = fmt.Sprintf(msg[0], expected)
	}
	result := Equal(value, expected)
	v.SetError(result, key, m)
	return result
}

func (v *Validator) Contain(value, expected interface{}, key string, msg ...string) bool {
	var m string
	if len(msg) == 0 {
		m = fmt.Sprintf(Messages.Contain, expected)
	} else {
		m = fmt.Sprintf(msg[0], expected)
	}
	result := Contain(value, expected)
	v.SetError(result, key, m)
	return result
}

func (v *Validator) TimeRange(value, from, to time.Time, key string, msg ...string) bool {
	var m string
	if len(msg) == 0 {
		m = fmt.Sprintf(Messages.Contain, from, to)
	} else {
		m = fmt.Sprintf(msg[0], from, to)
	}
	result := TimeRange(value, from, to)
	v.SetError(result, key, m)
	return result
}
