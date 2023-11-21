package tengo

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"
	"unicode"
)

type structField struct {
	fields map[string]int
	method map[string]int
}

var (
	errorType       = reflect.TypeOf((*error)(nil)).Elem()
	objectType      = reflect.TypeOf((*Object)(nil)).Elem()
	callableType    = reflect.TypeOf((*CallableFunc)(nil)).Elem()
	timeType        = reflect.TypeOf((*time.Time)(nil)).Elem()
	boolType        = reflect.TypeOf((*bool)(nil)).Elem()
	stringType      = reflect.TypeOf((*string)(nil)).Elem()
	intType         = reflect.TypeOf((*int)(nil)).Elem()
	int8Type        = reflect.TypeOf((*int8)(nil)).Elem()
	int16Type       = reflect.TypeOf((*int16)(nil)).Elem()
	int32Type       = reflect.TypeOf((*int32)(nil)).Elem()
	int64Type       = reflect.TypeOf((*int64)(nil)).Elem()
	uintType        = reflect.TypeOf((*uint)(nil)).Elem()
	uint8Type       = reflect.TypeOf((*uint8)(nil)).Elem()
	uint16Type      = reflect.TypeOf((*uint16)(nil)).Elem()
	uint32Type      = reflect.TypeOf((*uint32)(nil)).Elem()
	uint64Type      = reflect.TypeOf((*uint64)(nil)).Elem()
	float32Type     = reflect.TypeOf((*float32)(nil)).Elem()
	float64Type     = reflect.TypeOf((*float64)(nil)).Elem()
	bytesType       = reflect.TypeOf((*[]byte)(nil)).Elem()
	structTypeCache sync.Map
)

func smallCamelCaseName(s string) string {
	if s[0] >= 'A' && s[0] <= 'Z' {
		runes := []rune(s)
		runes[0] = unicode.ToLower(runes[0])
		return string(runes)
	}
	return s
}

func structFields(t reflect.Type) *structField {
	if m, ok := structTypeCache.Load(t); ok {
		return m.(*structField)
	}

	m := &structField{
		fields: make(map[string]int),
		method: make(map[string]int),
	}

	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if !field.IsExported() {
				continue
			}
			m.fields[smallCamelCaseName(field.Name)] = i
		}
	}

	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		if !method.IsExported() {
			continue
		}
		m.method[smallCamelCaseName(method.Name)] = method.Index
	}

	structTypeCache.Store(t, m)

	return m
}

func AssignBool(rtype reflect.Type, o Object) (reflect.Value, bool) {
	if rtype.Kind() != reflect.Bool {
		return reflect.Value{}, false
	}
	if b, ok := ToBool(o); ok {
		v := reflect.New(rtype).Elem()
		v.SetBool(b)
		return v, true
	}
	return reflect.Value{}, false
}

func AssignString(rtype reflect.Type, o Object) (reflect.Value, bool) {
	if rtype.Kind() != reflect.String {
		return reflect.Value{}, false
	}
	if s, ok := ToString(o); ok {
		v := reflect.New(rtype).Elem()
		v.SetString(s)
		return v, true
	}
	return reflect.Value{}, false
}

func AssignInt(rtype reflect.Type, o Object) (reflect.Value, bool) {
	if rtype.Kind() != reflect.Int {
		return reflect.Value{}, false
	}
	if i, ok := ToNumber(o, NumberTypeInt); ok {
		v := reflect.New(rtype).Elem()
		v.SetInt(i)
		return v, true
	}
	return reflect.Value{}, false
}

func AssignInt8(rtype reflect.Type, o Object) (reflect.Value, bool) {
	if rtype.Kind() != reflect.Int8 {
		return reflect.Value{}, false
	}
	if i, ok := ToNumber(o, NumberTypeInt8); ok {
		v := reflect.New(rtype).Elem()
		v.SetInt(i)
		return v, true
	}
	return reflect.Value{}, false
}

func AssignInt16(rtype reflect.Type, o Object) (reflect.Value, bool) {
	if rtype.Kind() != reflect.Int16 {
		return reflect.Value{}, false
	}
	if i, ok := ToNumber(o, NumberTypeInt16); ok {
		v := reflect.New(int16Type).Elem()
		v.SetInt(i)
		return v, true
	}
	return reflect.Value{}, false
}

func AssignInt32(rtype reflect.Type, o Object) (reflect.Value, bool) {
	if rtype.Kind() != reflect.Int32 {
		return reflect.Value{}, false
	}
	if i, ok := ToNumber(o, NumberTypeInt32); ok {
		v := reflect.New(int32Type).Elem()
		v.SetInt(i)
		return v, true
	}
	return reflect.Value{}, false
}

func AssignInt64(rtype reflect.Type, o Object) (reflect.Value, bool) {
	if rtype.Kind() != reflect.Int64 {
		return reflect.Value{}, false
	}
	if i, ok := ToNumber(o, NumberTypeInt64); ok {
		v := reflect.New(rtype).Elem()
		v.SetInt(i)
		return v, true
	}
	return reflect.Value{}, false
}

func AssignUint(rtype reflect.Type, o Object) (reflect.Value, bool) {
	if rtype.Kind() != reflect.Uint {
		return reflect.Value{}, false
	}
	if u, ok := ToNumber(o, NumberTypeUint); ok {
		v := reflect.New(rtype).Elem()
		v.SetUint(uint64(u))
		return v, true
	}
	return reflect.Value{}, false
}

func AssignUint8(rtype reflect.Type, o Object) (reflect.Value, bool) {
	if rtype.Kind() != reflect.Uint8 {
		return reflect.Value{}, false
	}
	if u, ok := ToNumber(o, NumberTypeUint8); ok {
		v := reflect.New(rtype).Elem()
		v.SetUint(uint64(u))
		return v, true
	}
	return reflect.Value{}, false
}

func AssignUint16(rtype reflect.Type, o Object) (reflect.Value, bool) {
	if rtype.Kind() != reflect.Uint16 {
		return reflect.Value{}, false
	}
	if u, ok := ToNumber(o, NumberTypeUint16); ok {
		v := reflect.New(rtype).Elem()
		v.SetUint(uint64(u))
		return v, true
	}
	return reflect.Value{}, false
}

func AssignUint32(rtype reflect.Type, o Object) (reflect.Value, bool) {
	if rtype.Kind() != reflect.Uint32 {
		return reflect.Value{}, false
	}
	if u, ok := ToNumber(o, NumberTypeUint32); ok {
		v := reflect.New(rtype).Elem()
		v.SetUint(uint64(u))
		return v, true
	}
	return reflect.Value{}, false
}

func AssignUint64(rtype reflect.Type, o Object) (reflect.Value, bool) {
	if rtype.Kind() != reflect.Uint64 {
		return reflect.Value{}, false
	}
	if u, ok := ToNumber(o, NumberTypeUint64); ok {
		v := reflect.New(rtype).Elem()
		v.SetUint(uint64(u))
		return v, true
	}
	return reflect.Value{}, false
}

func AssignFloat32(rtype reflect.Type, o Object) (reflect.Value, bool) {
	if rtype.Kind() != reflect.Float32 {
		return reflect.Value{}, false
	}
	if f, ok := ToFloat64(o); ok {
		v := reflect.New(rtype).Elem()
		v.SetFloat(f)
		return v, true
	}
	return reflect.Value{}, false
}

func AssignFloat64(rtype reflect.Type, o Object) (reflect.Value, bool) {
	if rtype.Kind() != reflect.Float64 {
		return reflect.Value{}, false
	}
	if f, ok := ToFloat64(o); ok {
		v := reflect.New(rtype).Elem()
		v.SetFloat(f)
		return v, true
	}
	return reflect.Value{}, false
}

func AssignArray(t reflect.Type, o Object) (reflect.Value, bool) {
	switch t.Kind() {
	case reflect.Array, reflect.Slice:
	default:
		return reflect.Value{}, false
	}

	sliceElemType := t.Elem()
	if sliceElemType.Kind() == reflect.Uint8 {
		if bytes, ok := o.(*Bytes); ok {
			b := bytes.Value
			count := len(b)
			v := reflect.MakeSlice(t, count, count)
			for i := 0; i < count; i++ {
				v.Index(i).Set(reflect.ValueOf(b[i]))
			}
			return v, true
		}
	}

	var objs []Object

	switch arr := o.(type) {
	case *ImmutableArray:
		objs = arr.Value
	case *Array:
		objs = arr.Value
	default:
		return reflect.Value{}, false
	}

	count := len(objs)

	if count < 1 {
		return reflect.Zero(t), true
	}

	v := reflect.MakeSlice(t, count, count)
	for i := 0; i < count; i++ {
		if elem, err := AssignValue(sliceElemType, objs[i]); err != nil {
			return reflect.Value{}, false
		} else {
			v.Index(i).Set(elem)
		}
	}

	return v, true
}

func AssignInterface(t reflect.Type, o Object) (reflect.Value, bool) {
	switch t.Kind() {
	case reflect.Interface:
	default:
		return reflect.Value{}, false
	}

	if t.Implements(objectType) {
		return reflect.ValueOf(o.Copy()), true
	}
	if t.Implements(errorType) {
		return reflect.ValueOf(errors.New(o.String())), true
	}

	return reflect.Value{}, false
}

func AssignPointer(t reflect.Type, o Object) (reflect.Value, bool) {
	switch t.Kind() {
	case reflect.Pointer:
	default:
		return reflect.Value{}, false
	}

	if v, err := AssignValue(t.Elem(), o); err != nil {
		return reflect.Value{}, false
	} else {
		i := reflect.New(t)
		i.Elem().Set(v)
		return i, true
	}
}

func AssignMap(t reflect.Type, o Object) (reflect.Value, bool) {
	switch t.Kind() {
	case reflect.Map:
	default:
		return reflect.Value{}, false
	}

	var m map[string]Object

	switch mObj := o.(type) {
	case *ImmutableMap:
		m = mObj.Value
	case *Map:
		m = mObj.Value
	default:
		return reflect.Value{}, false
	}

	kType, vType := t.Key(), t.Elem()

	if kType.Kind() != reflect.String {
		return reflect.Value{}, false
	}

	v := reflect.MakeMapWithSize(t, len(m))
	for key, value := range m {
		if vObj, err := AssignValue(vType, value); err != nil {
			return reflect.Value{}, false
		} else {
			v.SetMapIndex(reflect.ValueOf(key), vObj)
		}
	}

	return v, true
}

func AssignStruct(t reflect.Type, o Object) (reflect.Value, bool) {
	switch t.Kind() {
	case reflect.Struct:
	default:
		return reflect.Value{}, false
	}

	if objErr, ok := o.(*Error); ok {
		o = objErr.Value
	}

	var objFields map[string]Object

	switch obj := o.(type) {
	case *Time:
		if t.AssignableTo(timeType) {
			timeValue := reflect.New(t).Elem()
			timeValue.Set(reflect.ValueOf(obj.Value))
			return timeValue, true
		}
		return reflect.Value{}, false
	case *ImmutableMap:
		objFields = obj.Value
	case *Map:
		objFields = obj.Value
	}

	sFields := structFields(t)

	v := reflect.New(t).Elem()

	for key, value := range objFields {
		if index, ok := sFields.fields[key]; ok {
			field := v.Field(index)
			fieldValue, err := AssignValue(field.Type(), value)
			if err != nil {
				return reflect.Value{}, false
			}
			field.Set(fieldValue)
		}
	}

	return v, true
}

func AssignValue(rtype reflect.Type, o Object) (reflect.Value, error) {
	switch rtype.Kind() {
	case reflect.Bool:
		if v, ok := AssignBool(rtype, o); ok {
			return v, nil
		}
	case reflect.Int:
		if v, ok := AssignInt(rtype, o); ok {
			return v, nil
		}
	case reflect.Int8:
		if v, ok := AssignInt8(rtype, o); ok {
			return v, nil
		}
	case reflect.Int16:
		if v, ok := AssignInt16(rtype, o); ok {
			return v, nil
		}
	case reflect.Int32:
		if v, ok := AssignInt32(rtype, o); ok {
			return v, nil
		}
	case reflect.Int64:
		if v, ok := AssignInt64(rtype, o); ok {
			return v, nil
		}
	case reflect.Uint:
		if v, ok := AssignUint(rtype, o); ok {
			return v, nil
		}
	case reflect.Uint8:
		if v, ok := AssignUint8(rtype, o); ok {
			return v, nil
		}
	case reflect.Uint16:
		if v, ok := AssignUint16(rtype, o); ok {
			return v, nil
		}
	case reflect.Uint32:
		if v, ok := AssignUint32(rtype, o); ok {
			return v, nil
		}
	case reflect.Uint64:
		if v, ok := AssignUint64(rtype, o); ok {
			return v, nil
		}
	case reflect.Float32:
		if v, ok := AssignFloat32(rtype, o); ok {
			return v, nil
		}
	case reflect.Float64:
		if v, ok := AssignFloat64(rtype, o); ok {
			return v, nil
		}
	case reflect.String:
		if v, ok := AssignString(rtype, o); ok {
			return v, nil
		}
	case reflect.Array, reflect.Slice:
		if v, ok := AssignArray(rtype, o); ok {
			return v, nil
		}
	case reflect.Func:
		// not supported for now
	case reflect.Interface:
		if v, ok := AssignInterface(rtype, o); ok {
			return v, nil
		}
	case reflect.Pointer:
		if v, ok := AssignPointer(rtype, o); ok {
			return v, nil
		}
	case reflect.Map:
		if v, ok := AssignMap(rtype, o); ok {
			return v, nil
		}
	case reflect.Struct:
		if v, ok := AssignStruct(rtype, o); ok {
			return v, nil
		}
	}

	return reflect.Value{}, fmt.Errorf("cannot assign to type: %s -> %s",
		o.TypeName(),
		rtype.Kind().String())
}

func FromUintBool(v reflect.Value) (Object, error) {
	kind := v.Kind()
	if kind != reflect.Bool {
		return nil, fmt.Errorf("cannot convert to bool: %s", kind.String())
	}
	if v.Bool() {
		return TrueValue, nil
	}
	return FalseValue, nil
}

func FromValueInt(v reflect.Value) (Object, error) {
	kind := v.Kind()
	if kind != reflect.Int {
		return nil, fmt.Errorf("cannot convert to int: %s", kind.String())
	}
	return &Number{Value: v.Int(), Type: NumberTypeInt}, nil
}

func FromValueInt8(v reflect.Value) (Object, error) {
	kind := v.Kind()
	if kind != reflect.Int8 {
		return nil, fmt.Errorf("cannot convert to int8: %s", kind.String())
	}
	return &Number{Value: v.Int(), Type: NumberTypeInt8}, nil
}

func FromValueInt16(v reflect.Value) (Object, error) {
	kind := v.Kind()
	if kind != reflect.Int16 {
		return nil, fmt.Errorf("cannot convert to int16: %s", kind.String())
	}
	return &Number{Value: v.Int(), Type: NumberTypeInt16}, nil
}

func FromValueInt32(v reflect.Value) (Object, error) {
	kind := v.Kind()
	if kind != reflect.Int32 {
		return nil, fmt.Errorf("cannot convert to int32: %s", kind.String())
	}
	return &Number{Value: v.Int(), Type: NumberTypeInt32}, nil
}

func FromValueInt64(v reflect.Value) (Object, error) {
	kind := v.Kind()
	if kind != reflect.Int64 {
		return nil, fmt.Errorf("cannot convert to int64: %s", kind.String())
	}
	return &Number{Value: v.Int(), Type: NumberTypeInt64}, nil
}

func FromValueUint(v reflect.Value) (Object, error) {
	kind := v.Kind()
	if kind != reflect.Uint {
		return nil, fmt.Errorf("cannot convert to uint: %s", kind.String())
	}
	return &Number{Value: int64(v.Uint()), Type: NumberTypeUint}, nil
}

func FromValueUint8(v reflect.Value) (Object, error) {
	kind := v.Kind()
	if kind != reflect.Uint8 {
		return nil, fmt.Errorf("cannot convert to uint8: %s", kind.String())
	}
	return &Number{Value: int64(v.Uint()), Type: NumberTypeUint8}, nil
}

func FromValueUint16(v reflect.Value) (Object, error) {
	kind := v.Kind()
	if kind != reflect.Uint16 {
		return nil, fmt.Errorf("cannot convert to uint16: %s", kind.String())
	}
	return &Number{Value: int64(v.Uint()), Type: NumberTypeUint16}, nil
}

func FromValueUint32(v reflect.Value) (Object, error) {
	kind := v.Kind()
	if kind != reflect.Uint32 {
		return nil, fmt.Errorf("cannot convert to uint32: %s", kind.String())
	}
	return &Number{Value: int64(v.Uint()), Type: NumberTypeUint32}, nil
}

func FromValueUint64(v reflect.Value) (Object, error) {
	kind := v.Kind()
	if kind != reflect.Uint64 {
		return nil, fmt.Errorf("cannot convert to uint64: %s", kind.String())
	}
	return &Number{Value: int64(v.Uint()), Type: NumberTypeUint64}, nil
}

func FromValueFloat32(v reflect.Value) (Object, error) {
	kind := v.Kind()
	if kind != reflect.Float32 {
		return nil, fmt.Errorf("cannot convert to float32: %s", kind.String())
	}
	return &Float{Value: v.Float()}, nil
}

func FromValueFloat64(v reflect.Value) (Object, error) {
	kind := v.Kind()
	if kind != reflect.Float64 {
		return nil, fmt.Errorf("cannot convert to float64: %s", kind.String())
	}
	return &Float{Value: v.Float()}, nil
}

func FromValueString(v reflect.Value) (Object, error) {
	kind := v.Kind()
	if kind != reflect.String {
		return nil, fmt.Errorf("cannot convert to string: %s", kind.String())
	}
	str := v.String()
	if len(str) > MaxStringLen {
		return nil, ErrStringLimit
	}
	return &String{Value: str}, nil
}

func FromValueFunc(name string, v reflect.Value) (Object, error) {
	rtype := v.Type()

	if rtype.Kind() != reflect.Func || v.IsNil() {
		return nil, fmt.Errorf("cannot convert to CallableFunc: %s", rtype.Kind().String())
	}

	if rtype.AssignableTo(callableType) {
		return &UserFunction{Name: name, Value: v.Interface().(CallableFunc)}, nil
	}

	return &UserFunction{
		Name: name,
		Value: func(args ...Object) (ret Object, err error) {
			isVariadic := rtype.IsVariadic()

			numIn, numOut := rtype.NumIn(), rtype.NumOut()

			var in, out []reflect.Value

			var argValue reflect.Value

			for i := 0; i < numIn; i++ {
				assignType := rtype.In(i)

				if isVariadic && i+1 >= numIn {
					if len(args) < numIn {
						in = append(in, reflect.Zero(assignType))
						break
					}

					arrType := assignType.Elem()
					variadicValue := reflect.MakeSlice(assignType, len(args)-i, len(args)-i)
					for x := i; x < len(args); x++ {
						if argValue, err = AssignValue(arrType, args[x]); err != nil {
							return
						}
						variadicValue.Index(x - i).Set(argValue)
					}
					in = append(in, variadicValue)
				} else {
					if len(args) < i+1 {
						return nil, ErrWrongNumArguments
					}

					if argValue, err = AssignValue(assignType, args[i]); err != nil {
						return
					}
					in = append(in, argValue)
				}
			}

			if !isVariadic {
				out = v.Call(in)
			} else {
				out = v.CallSlice(in)
			}

			if numOut < 1 {
				return
			}

			if rtype.Out(numOut - 1).Implements(errorType) {
				outErr := out[numOut-1]
				if outErr.IsValid() && !outErr.IsNil() &&
					outErr.Type().Implements(errorType) {
					return nil, outErr.Interface().(error)
				}
				numOut = numOut - 1
			}

			var (
				obj  Object
				objs []Object
			)

			for i := 0; i < numOut; i++ {
				if obj, err = FromValue(out[i]); err != nil {
					return
				}
				objs = append(objs, obj)
			}

			if numOut > 1 {
				return &ImmutableArray{Value: objs}, nil
			}

			return obj, nil
		},
	}, nil
}

func FromValueArray(v reflect.Value) (Object, error) {
	rtype := v.Type()

	switch rtype.Kind() {
	case reflect.Array, reflect.Slice:
	default:
		return nil, fmt.Errorf("cannot convert to Array: %s", rtype.Kind().String())
	}

	if rtype.Elem().AssignableTo(uint8Type) {
		bytes := v.Interface().([]byte)
		if len(bytes) > MaxBytesLen {
			return nil, ErrBytesLimit
		}
		return &Bytes{Value: bytes}, nil
	}

	var arr []Object
	for i := 0; i < v.Len(); i++ {
		obj, err := FromValue(v.Index(i))
		if err != nil {
			return nil, err
		}
		arr = append(arr, obj)
	}
	return &Array{Value: arr}, nil
}

func FromValueMap(v reflect.Value) (Object, error) {
	rtype := v.Type()

	if rtype.Kind() != reflect.Map {
		return nil, fmt.Errorf("cannot convert to Map: %s", rtype.Kind().String())
	}

	keyType := rtype.Key()

	if keyType.Kind() != reflect.String {
		return nil, fmt.Errorf("cannot convert map key to string: %s", keyType.String())
	}

	m := make(map[string]Object)
	iter := v.MapRange()
	for iter.Next() {
		iterValue, err := FromValue(iter.Value())
		if err != nil {
			return nil, err
		}
		m[iter.Key().String()] = iterValue
	}

	return &Map{Value: m}, nil
}

func FromValueStruct(v reflect.Value) (Object, error) {
	rtype := v.Type()

	if rtype.Kind() != reflect.Struct {
		return nil, fmt.Errorf("cannot convert to Struct: %s", rtype.Kind().String())
	}

	if rtype.AssignableTo(timeType) {
		return &Time{Value: v.Interface().(time.Time)}, nil
	}

	m := make(map[string]Object)

	setFields := func(v reflect.Value, rtype reflect.Type) error {
		sFields := structFields(rtype)

		for key, index := range sFields.fields {
			o, err := FromValue(v.Field(index))
			if err != nil {
				return err
			}
			m[key] = o
		}

		for key, index := range sFields.method {
			methodName := rtype.Method(index).Name
			o, err := FromValueFunc(methodName, v.Method(index))
			if err != nil {
				return err
			}
			m[key] = o
		}

		return nil
	}

	if err := setFields(v, rtype); err != nil {
		return nil, err
	}

	if v.CanAddr() {
		vPtr := v.Addr()
		if err := setFields(vPtr, vPtr.Type()); err != nil {
			return nil, err
		}
	}

	return &Map{Value: m}, nil
}

func fromValueOrPointer(v reflect.Value) (Object, error) {
	if v.IsNil() {
		return UndefinedValue, nil
	}

	rtype := v.Type()

	if rtype.Implements(errorType) {
		return &Error{Value: &String{Value: errorType.(error).Error()}}, nil
	}

	if rtype.Implements(objectType) {
		return v.Interface().(Object), nil
	}

	return FromValue(v.Elem())
}

func FromValuePointer(v reflect.Value) (Object, error) {
	rtype := v.Type()
	if v.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("cannot convert to pointer: *%s", rtype.Elem().Name())
	}
	return fromValueOrPointer(v)
}

func FromValueInterface(v reflect.Value) (Object, error) {
	kind := v.Kind()
	if kind != reflect.Interface {
		return nil, fmt.Errorf("cannot convert to interface: %s", kind.String())
	}
	return fromValueOrPointer(v)
}

// FromValue will attempt to convert an reflect.Value v to a Tengo Object
func FromValue(v reflect.Value) (Object, error) {
	switch v.Kind() {
	case reflect.Invalid:
		return UndefinedValue, nil
	case reflect.Func, reflect.Map, reflect.Pointer, reflect.Interface, reflect.Slice:
		if v.IsNil() {
			return UndefinedValue, nil
		}
	}

	rtype := v.Type()
	kind := rtype.Kind()
	switch kind {
	case reflect.Bool:
		return FromUintBool(v)
	case reflect.Int:
		return FromValueInt(v)
	case reflect.Int8:
		return FromValueInt8(v)
	case reflect.Int16:
		return FromValueInt16(v)
	case reflect.Int32:
		return FromValueInt32(v)
	case reflect.Int64:
		return FromValueInt64(v)
	case reflect.Uint:
		return FromValueUint(v)
	case reflect.Uint8:
		return FromValueUint8(v)
	case reflect.Uint16:
		return FromValueUint16(v)
	case reflect.Uint32:
		return FromValueUint32(v)
	case reflect.Uint64:
		return FromValueUint64(v)
	case reflect.Float32:
		return FromValueFloat32(v)
	case reflect.Float64:
		return FromValueFloat64(v)
	case reflect.Array, reflect.Slice:
		return FromValueArray(v)
	case reflect.Func:
		return FromValueFunc("", v)
	case reflect.Interface:
		return FromValueInterface(v)
	case reflect.Pointer:
		return FromValuePointer(v)
	case reflect.Map:
		return FromValueMap(v)
	case reflect.String:
		return FromValueString(v)
	case reflect.Struct:
		return FromValueStruct(v)
	}

	return nil, fmt.Errorf("cannot convert to Object: %s", kind.String())
}
