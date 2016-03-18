package opentok

import (
	"reflect"
)

// IsZero determines whether or not a variable/field/whatever is of it's type's zero value
// Pass reflect.ValueOf(x)
// http://stackoverflow.com/a/23555352/3183170
func IsZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Func, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.Array:
		z := true
		for i := 0; i < v.Len(); i++ {
			z = z && IsZero(v.Index(i))
		}
		return z
	case reflect.Struct:
		z := true
		for i := 0; i < v.NumField(); i++ {
			z = z && IsZero(v.Field(i))
		}
		return z
	}
	// Compare other types directly:
	z := reflect.Zero(v.Type())
	return v.Interface() == z.Interface()
}

// Update this to create and return a new struct instead of mutating input
func Update(mainObj interface{}, newData interface{}) bool {
	newDataVal, mainObjVal := reflect.ValueOf(newData).Elem(), reflect.ValueOf(mainObj).Elem()
	fieldCount := newDataVal.NumField()
	changed := false
	for i := 0; i < fieldCount; i++ {
		newField := newDataVal.Field(i)
		// They passed in a value for this field, update our DB user
		if newField.IsValid() && !IsZero(newField) {
			dbField := mainObjVal.Field(i)
			dbField.Set(newField)
			changed = true
		}
	}
	return changed
}

// Extend is the same as Update above, but does not mutate user input
func Extend(mainObj interface{}, newData interface{}) reflect.Value {
	finalObjVal := reflect.ValueOf(mainObj).Elem()
	newDataVal := reflect.ValueOf(newData).Elem()
	fieldCount := newDataVal.NumField()
	for i := 0; i < fieldCount; i++ {
		newValue := newDataVal.Field(i)
		// They passed in a value for this field, update our DB user
		if newValue.IsValid() && !IsZero(newValue) {
			currentField := finalObjVal.Field(i)
			currentField.Set(newValue)
		}
	}
	return finalObjVal
}
