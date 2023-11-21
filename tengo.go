package tengo

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

var (
	// MaxStringLen is the maximum byte-length for string value. Note this
	// limit applies to all compiler/VM instances in the process.
	MaxStringLen = 2147483647

	// MaxBytesLen is the maximum length for bytes value. Note this limit
	// applies to all compiler/VM instances in the process.
	MaxBytesLen = 2147483647
)

const (
	// GlobalsSize is the maximum number of global variables for a VM.
	GlobalsSize = 1024

	// StackSize is the maximum stack size for a VM.
	StackSize = 2048

	// MaxFrames is the maximum number of function frames for a VM.
	MaxFrames = 1024

	// SourceFileExtDefault is the default extension for source files.
	SourceFileExtDefault = ".tengo"
)

// CallableFunc is a function signature for the callable functions.
type CallableFunc = func(args ...Object) (ret Object, err error)

// CountObjects returns the number of objects that a given object o contains.
// For scalar value types, it will always be 1. For compound value types,
// this will include its elements and all of their elements recursively.
func CountObjects(o Object) (c int) {
	c = 1
	switch o := o.(type) {
	case *Array:
		for _, v := range o.Value {
			c += CountObjects(v)
		}
	case *ImmutableArray:
		for _, v := range o.Value {
			c += CountObjects(v)
		}
	case *Map:
		for _, v := range o.Value {
			c += CountObjects(v)
		}
	case *ImmutableMap:
		for _, v := range o.Value {
			c += CountObjects(v)
		}
	case *Error:
		c += CountObjects(o.Value)
	}
	return
}

// ToString will try to convert object o to string value.
func ToString(o Object) (v string, ok bool) {
	if o == UndefinedValue {
		return
	}
	ok = true
	if str, isStr := o.(*String); isStr {
		v = str.Value
	} else {
		v = o.String()
	}
	return
}

// ToInt will try to convert object o to int value.
func ToInt(o Object) (v int, ok bool) {
	switch o := o.(type) {
	case *Number:
		v = int(o.Value)
		ok = true
	case *Float:
		v = int(o.Value)
		ok = true
	case *Char:
		v = int(o.Value)
		ok = true
	case *Bool:
		if o == TrueValue {
			v = 1
		}
		ok = true
	case *String:
		c, err := strconv.ParseInt(o.Value, 10, 64)
		if err == nil {
			v = int(c)
			ok = true
		}
	}
	return
}

func ToNumber(o Object, nt NumberType) (v int64, ok bool) {
	switch o := o.(type) {
	case *Number:
		v = o.Value
		ok = true
	case *Float:
		v = int64(o.Value)
		ok = true
	case *Char:
		v = int64(o.Value)
		ok = true
	case *Bool:
		if o == TrueValue {
			v = 1
		}
		ok = true
	case *String:
		bitSize := 64
		switch nt {
		case NumberTypeUint8, NumberTypeInt8:
			bitSize = 8
		case NumberTypeUint16, NumberTypeInt16:
			bitSize = 16
		case NumberTypeUint32, NumberTypeInt32:
			bitSize = 32
		}
		c, err := strconv.ParseInt(o.Value, 10, bitSize)
		if err == nil {
			v = c
			ok = true
		}
	}
	return
}

// ToInt64 will try to convert object o to int64 value.
func ToInt64(o Object) (v int64, ok bool) {
	return ToNumber(o, NumberTypeInt64)
}

// ToFloat64 will try to convert object o to float64 value.
func ToFloat64(o Object) (v float64, ok bool) {
	switch o := o.(type) {
	case *Number:
		v = float64(o.Value)
		ok = true
	case *Float:
		v = o.Value
		ok = true
	case *String:
		c, err := strconv.ParseFloat(o.Value, 64)
		if err == nil {
			v = c
			ok = true
		}
	}
	return
}

// ToBool will try to convert object o to bool value.
func ToBool(o Object) (v bool, ok bool) {
	ok = true
	v = !o.IsFalsy()
	return
}

// ToRune will try to convert object o to rune value.
func ToRune(o Object) (v rune, ok bool) {
	switch o := o.(type) {
	case *Number:
		v = rune(o.Value)
		ok = true
	case *Char:
		v = o.Value
		ok = true
	}
	return
}

// ToByteSlice will try to convert object o to []byte value.
func ToByteSlice(o Object) (v []byte, ok bool) {
	switch o := o.(type) {
	case *Bytes:
		v = o.Value
		ok = true
	case *String:
		v = []byte(o.Value)
		ok = true
	}
	return
}

// ToTime will try to convert object o to time.Time value.
func ToTime(o Object) (v time.Time, ok bool) {
	switch o := o.(type) {
	case *Time:
		v = o.Value
		ok = true
	case *Number:
		v = time.Unix(o.Value, 0)
		ok = true
	}
	return
}

// ToInterface attempts to convert an object o to an interface{} value
func ToInterface(o Object) (res interface{}) {
	var assignType reflect.Type
	switch obj := o.(type) {
	case *Number:
		switch obj.Type {
		case NumberTypeInt:
			assignType = intType
		case NumberTypeInt8:
			assignType = int8Type
		case NumberTypeInt16:
			assignType = int16Type
		case NumberTypeInt32:
			assignType = int32Type
		case NumberTypeInt64:
			assignType = int64Type
		case NumberTypeUint:
			assignType = uintType
		case NumberTypeUint8:
			assignType = uint8Type
		case NumberTypeUint16:
			assignType = uint16Type
		case NumberTypeUint32:
			assignType = uint32Type
		case NumberTypeUint64:
			assignType = uint64Type
		}
	case *String:
		assignType = stringType
	case *Float:
		assignType = float64Type
	case *Bool:
		assignType = boolType
	case *Char:
		assignType = uint8Type
	case *Bytes:
		assignType = bytesType
	case *Array:
		arr := make([]interface{}, len(obj.Value), len(obj.Value))
		for i, val := range obj.Value {
			arr[i] = ToInterface(val)
		}
		return arr
	case *ImmutableArray:
		arr := make([]interface{}, len(obj.Value), len(obj.Value))
		for i, val := range obj.Value {
			arr[i] = ToInterface(val)
		}
		return arr
	case *Map:
		m := make(map[string]interface{})
		for key, v := range obj.Value {
			m[key] = ToInterface(v)
		}
		return m
	case *ImmutableMap:
		m := make(map[string]interface{})
		for key, v := range obj.Value {
			m[key] = ToInterface(v)
		}
		return m
	case *Time:
		return obj.Value
	case *Error:
		return errors.New(o.String())
	case *Undefined:
		return nil
	case Object:
		return o
	}

	if v, err := AssignValue(assignType, o); err == nil {
		res = v.Interface()
	}

	return
}

// FromInterface will attempt to convert an interface{} v to a Tengo Object
func FromInterface(v interface{}) (Object, error) {
	return FromValue(reflect.ValueOf(v))
}
