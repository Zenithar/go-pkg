package types

import "reflect"

// IsNil returns true if given object is nil
func IsNil(c interface{}) bool {
	return c == nil || (reflect.ValueOf(c).Kind() == reflect.Ptr && reflect.ValueOf(c).IsNil())
}
