package helpers

import "reflect"

func Filter(params []string, f func(s string) bool) []string {
	result := make([]string, 0)
	for _, p := range params {
		if f(p) {
			result = append(result, p)
		}
	}
	return result
}

func Map(params interface{}, f func(s interface{}) interface{}) []interface{} {
	s := reflect.ValueOf(params)
	if s.Kind() != reflect.Slice {
		panic("The params are not slice")
	}

	result := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		result[i] = f(s.Index(i).Interface())
	}
	return result
}
