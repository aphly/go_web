package helper

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

func Rand(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min+1) + min
}

func RandStr(lenNum int, diff ...int64) string {
	var chars = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o",
		"p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G",
		"H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y",
		"Z", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
	str := ""
	length := len(chars)
	var r *rand.Rand
	if len(diff) == 0 {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	} else {
		r = rand.New(rand.NewSource(time.Now().UnixNano() + diff[0]))
	}
	for i := 0; i < lenNum; i++ {
		str += chars[r.Intn(length)]
	}
	return str
}

func PageOffset(page, pageSize int) (offset int) {
	return (page - 1) * pageSize
}

func Pages(total, PerPage int) int {
	pages := total / PerPage
	if total%PerPage > 0 {
		pages++
	}
	return pages
}

func ToCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i := 1; i < len(parts); i++ {
		parts[i] = strings.ToUpper(parts[i][0:1]) + strings.ToLower(parts[i][1:])
	}
	return strings.Join(parts, "")
}

func MapToStruct(m map[string]any, s any) error {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr && !v.Elem().CanSet() {
		return fmt.Errorf("Converter: cannot set value to %T", s)
	}

	for k, i := range m {
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		fmt.Println(k)

		field := v.FieldByName(k)
		fmt.Println(field)
		if !field.IsValid() {
			continue
		}
		if !field.CanSet() {
			continue
		}
		val := reflect.ValueOf(i)
		if field.Type() != val.Type() {
			return fmt.Errorf("Converter: cannot convert %v to %v", val.Type(), field.Type())
		}
		field.Set(val)
	}
	return nil
}
