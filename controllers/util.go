package controllers

import (
	"reflect"
)

// isZero return true if v is default zero value
// v or v's field must be exported, otherwise it will panic
func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Func, reflect.Map:
		return v.IsNil()
	case reflect.Array, reflect.Slice:
		z := true
		for i := 0; i < v.Len(); i++ {
			z = z && isZero(v.Index(i))
		}
		return z
	case reflect.Struct:
		z := true
		for i := 0; i < v.NumField(); i++ {
			z = z && isZero(v.Field(i))
		}
		return z
	}

	// Compare other types directly:
	z := reflect.Zero(v.Type())
	return v.Interface() == z.Interface()
}

// assignNonEmptyByFieldName assign non empty field specified by fields from src to dst.
// src and dst must be same type struct or pointer to struct
// only exported field will be assigned
func assignNonEmptyByFieldName(src, dst interface{}, fields ... string) {
	if reflect.TypeOf(src) != reflect.TypeOf(dst) {
		return
	}
	if len(fields) == 0 {
		return
	}

	srcVal := reflect.ValueOf(src)
	dstVal := reflect.ValueOf(dst)

	if srcVal.Kind() == reflect.Ptr {
		if srcVal.IsNil() {
			return
		}
		srcVal = srcVal.Elem()
		dstVal = dstVal.Elem()
	}

	for _, field := range fields {
		srcField := srcVal.FieldByName(field)
		dstField := dstVal.FieldByName(field)
		if srcField.CanSet() && !isZero(srcField) {
			dstField.Set(srcField)
		}
	}
}
