package logrus_fluent

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/Sirupsen/logrus"
)

func ConvertFields(fields logrus.Fields) map[string]interface{} {
	results := make(map[string]interface{})
	for key, value := range fields {
		results[key] = ConvertToValue(value)
	}
	return results
}

// ConvertToValue make map data from struct and tags
func ConvertToValue(p interface{}) interface{} {
	if stringer, ok := p.(fmt.Stringer); ok {
		return stringer.String()
	}

	if err, ok := p.(error); ok {
		return err.Error()
	}

	rv := toValue(p)

	switch rv.Kind() {
	case reflect.Struct:
		return convertToString(rv)
	case reflect.Map:
		return convertToString(rv)
	case reflect.Slice:
		return convertFromSlice(rv)
	case reflect.Invalid:
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int()
	case reflect.String:
		return rv.String()
	default:
		return p
	}
}

func convertFromMap(rv reflect.Value) interface{} {
	result := make(map[string]interface{})
	for _, key := range rv.MapKeys() {
		kv := rv.MapIndex(key)
		result[fmt.Sprint(key.Interface())] = ConvertToValue(kv.Interface())
	}
	return result
}

func convertFromSlice(rv reflect.Value) interface{} {
	var result []interface{}
	for i, max := 0, rv.Len(); i < max; i++ {
		result = append(result, ConvertToValue(rv.Index(i).Interface()))
	}
	return result
}

func convertToString(p reflect.Value) interface{} {
	i := p.Interface()
	if data, err := json.Marshal(i); err == nil {
		return string(data)
	}
	return fmt.Sprint(i)
}

// toValue converts any value to reflect.Value
func toValue(p interface{}) reflect.Value {
	v := reflect.ValueOf(p)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

// toType converts any value to reflect.Type
func toType(p interface{}) reflect.Type {
	t := reflect.ValueOf(p).Type()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}
