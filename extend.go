package tengo

import (
	"fmt"
	"reflect"
	"sync"
	"time"
	"unicode"
)

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
	interfaceType   = reflect.TypeOf((*interface{})(nil)).Elem()
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

type properties struct {
	fields map[string]int
	method map[string]int
}

func (p *properties) Instance(rtype reflect.Type, v reflect.Value) (map[string]Object, error) {
	if len(p.fields) < 1 && len(p.method) < 1 {
		return nil, nil
	}

	m := make(map[string]Object, len(p.fields)+len(p.method))

	for key, index := range p.fields {
		fieldObject, err := FromValue(v.Field(index))
		if err != nil {
			return nil, err
		}
		m[key] = fieldObject
	}

	for key, index := range p.method {
		methodObject, err := FromFunc(rtype.Method(index).Name, v.Method(index))
		if err != nil {
			return nil, err
		}
		m[key] = methodObject
	}

	return m, nil
}

func getProperties(t reflect.Type) *properties {
	if m, ok := structTypeCache.Load(t); ok {
		return m.(*properties)
	}

	m := &properties{}

	if t.Kind() == reflect.Struct {
		m.fields = make(map[string]int)
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if !field.IsExported() {
				continue
			}
			m.fields[smallCamelCaseName(field.Name)] = i
		}
	}

	if t.NumMethod() > 0 {
		m.method = make(map[string]int)
		for i := 0; i < t.NumMethod(); i++ {
			method := t.Method(i)
			if !method.IsExported() {
				continue
			}
			m.method[smallCamelCaseName(method.Name)] = method.Index
		}
	}

	structTypeCache.Store(t, m)
	return m
}

func assignError(t, e string) error {
	if e == "" {
		e = "unrecognized"
	}
	return fmt.Errorf("cannot assign to %s: %s", t, e)
}

func convertError(t, e string) error {
	if e == "" {
		e = "unrecognized"
	}
	return fmt.Errorf("cannot convert to %s: %s", t, e)
}

func isNil(o Object) bool {
	if o == nil || o == UndefinedValue {
		return true
	}
	return false
}

func AssignBool(rtype reflect.Type, o Object) (reflect.Value, error) {
	if rtype.Kind() != reflect.Bool {
		return reflect.Value{}, assignError("bool", rtype.Kind().String())
	}

	if isNil(o) {
		return reflect.Zero(rtype), nil
	}

	if b, ok := ToBool(o); ok {
		v := reflect.New(rtype).Elem()
		v.SetBool(b)
		return v, nil
	}

	return reflect.Value{}, assignError("bool", "")
}

func AssignString(rtype reflect.Type, o Object) (reflect.Value, error) {
	if rtype.Kind() != reflect.String {
		return reflect.Value{}, assignError("string", rtype.Kind().String())
	}

	if isNil(o) {
		return reflect.Zero(rtype), nil
	}

	if s, ok := ToString(o); ok {
		v := reflect.New(rtype).Elem()
		v.SetString(s)
		return v, nil
	}
	return reflect.Value{}, assignError("string", "")
}

func AssignInt(rtype reflect.Type, o Object) (reflect.Value, error) {
	if rtype.Kind() != reflect.Int {
		return reflect.Value{}, assignError("int", rtype.Kind().String())
	}

	if isNil(o) {
		return reflect.Zero(rtype), nil
	}

	if i, ok := ToNumber(o, Int); ok {
		v := reflect.New(rtype).Elem()
		v.SetInt(i)
		return v, nil
	}
	return reflect.Value{}, assignError("int", "")
}

func AssignInt8(rtype reflect.Type, o Object) (reflect.Value, error) {
	if rtype.Kind() != reflect.Int8 {
		return reflect.Value{}, assignError("int8", rtype.Kind().String())
	}

	if isNil(o) {
		return reflect.Zero(rtype), nil
	}

	if i, ok := ToNumber(o, Int8); ok {
		v := reflect.New(rtype).Elem()
		v.SetInt(i)
		return v, nil
	}
	return reflect.Value{}, assignError("int8", "")
}

func AssignInt16(rtype reflect.Type, o Object) (reflect.Value, error) {
	if rtype.Kind() != reflect.Int16 {
		return reflect.Value{}, assignError("int16", rtype.Kind().String())
	}

	if isNil(o) {
		return reflect.Zero(rtype), nil
	}

	if i, ok := ToNumber(o, Int16); ok {
		v := reflect.New(int16Type).Elem()
		v.SetInt(i)
		return v, nil
	}
	return reflect.Value{}, assignError("int16", "")
}

func AssignInt32(rtype reflect.Type, o Object) (reflect.Value, error) {
	if rtype.Kind() != reflect.Int32 {
		return reflect.Value{}, assignError("int32", rtype.Kind().String())
	}

	if isNil(o) {
		return reflect.Zero(rtype), nil
	}

	if i, ok := ToNumber(o, Int32); ok {
		v := reflect.New(int32Type).Elem()
		v.SetInt(i)
		return v, nil
	}
	return reflect.Value{}, assignError("int32", "")
}

func AssignInt64(rtype reflect.Type, o Object) (reflect.Value, error) {
	if rtype.Kind() != reflect.Int64 {
		return reflect.Value{}, assignError("int64", rtype.Kind().String())
	}

	if isNil(o) {
		return reflect.Zero(rtype), nil
	}

	if i, ok := ToNumber(o, Int64); ok {
		v := reflect.New(rtype).Elem()
		v.SetInt(i)
		return v, nil
	}
	return reflect.Value{}, assignError("int64", "")
}

func AssignUint(rtype reflect.Type, o Object) (reflect.Value, error) {
	if rtype.Kind() != reflect.Uint {
		return reflect.Value{}, assignError("uint", rtype.Kind().String())
	}

	if isNil(o) {
		return reflect.Zero(rtype), nil
	}

	if u, ok := ToNumber(o, Uint); ok {
		v := reflect.New(rtype).Elem()
		v.SetUint(uint64(u))
		return v, nil
	}
	return reflect.Value{}, assignError("uint", "")
}

func AssignUint8(rtype reflect.Type, o Object) (reflect.Value, error) {
	if rtype.Kind() != reflect.Uint8 {
		return reflect.Value{}, assignError("uint8", rtype.Kind().String())
	}

	if isNil(o) {
		return reflect.Zero(rtype), nil
	}

	if u, ok := ToNumber(o, Uint8); ok {
		v := reflect.New(rtype).Elem()
		v.SetUint(uint64(u))
		return v, nil
	}
	return reflect.Value{}, assignError("uint8", "")
}

func AssignUint16(rtype reflect.Type, o Object) (reflect.Value, error) {
	if rtype.Kind() != reflect.Uint16 {
		return reflect.Value{}, assignError("uint16", rtype.Kind().String())
	}

	if isNil(o) {
		return reflect.Zero(rtype), nil
	}

	if u, ok := ToNumber(o, Uint16); ok {
		v := reflect.New(rtype).Elem()
		v.SetUint(uint64(u))
		return v, nil
	}
	return reflect.Value{}, assignError("uint16", "")
}

func AssignUint32(rtype reflect.Type, o Object) (reflect.Value, error) {
	if rtype.Kind() != reflect.Uint32 {
		return reflect.Value{}, assignError("uint32", rtype.Kind().String())
	}

	if isNil(o) {
		return reflect.Zero(rtype), nil
	}

	if u, ok := ToNumber(o, Uint32); ok {
		v := reflect.New(rtype).Elem()
		v.SetUint(uint64(u))
		return v, nil
	}
	return reflect.Value{}, assignError("uint32", "")
}

func AssignUint64(rtype reflect.Type, o Object) (reflect.Value, error) {
	if rtype.Kind() != reflect.Uint64 {
		return reflect.Value{}, assignError("uint64", rtype.Kind().String())
	}

	if isNil(o) {
		return reflect.Zero(rtype), nil
	}

	if u, ok := ToNumber(o, Uint64); ok {
		v := reflect.New(rtype).Elem()
		v.SetUint(uint64(u))
		return v, nil
	}
	return reflect.Value{}, assignError("uint64", "")
}

func AssignFloat32(rtype reflect.Type, o Object) (reflect.Value, error) {
	if rtype.Kind() != reflect.Float32 {
		return reflect.Value{}, assignError("float32", rtype.Kind().String())
	}

	if isNil(o) {
		return reflect.Zero(rtype), nil
	}

	if f, ok := ToFloat64(o); ok {
		v := reflect.New(rtype).Elem()
		v.SetFloat(f)
		return v, nil
	}
	return reflect.Value{}, assignError("float32", "")
}

func AssignFloat64(rtype reflect.Type, o Object) (reflect.Value, error) {
	if rtype.Kind() != reflect.Float64 {
		return reflect.Value{}, assignError("float64", rtype.Kind().String())
	}

	if isNil(o) {
		return reflect.Zero(rtype), nil
	}

	if f, ok := ToFloat64(o); ok {
		v := reflect.New(rtype).Elem()
		v.SetFloat(f)
		return v, nil
	}
	return reflect.Value{}, assignError("float64", "")
}

func AssignArray(rtype reflect.Type, o Object) (reflect.Value, error) {
	var name string
	switch rtype.Kind() {
	case reflect.Array:
		name = "array"
	case reflect.Slice:
		name = "slice"
	default:
		return reflect.Value{}, assignError(name, rtype.Kind().String())
	}

	if isNil(o) {
		return reflect.Zero(rtype), nil
	}

	compareArrayLen := func(len int) error {
		if rtype.Kind() == reflect.Array {
			if rtype.Len() != len {
				return assignError(name, "len overflow")
			}
		}
		return nil
	}

	sliceElemType := rtype.Elem()
	if sliceElemType.Kind() == reflect.Uint8 && sliceElemType == uint8Type {
		if bytes, ok := o.(*Bytes); ok {
			if bytesType == rtype {
				return reflect.ValueOf(bytes.Value), nil
			}

			valueLen := len(bytes.Value)
			if err := compareArrayLen(valueLen); err != nil {
				return reflect.Value{}, err
			}

			v := reflect.MakeSlice(rtype, valueLen, valueLen)
			for i := 0; i < valueLen; i++ {
				v.Index(i).Set(reflect.ValueOf(bytes.Value[i]))
			}
			return v, nil
		}
	}

	if ref, ok := o.(Reference); ok {
		refType := ref.RefType()
		if refType.AssignableTo(rtype) {
			return ref.RefValue(), nil
		}
		o = ref.Into()
	}

	var objs []Object

	switch arr := o.(type) {
	case *ImmutableArray:
		objs = arr.Value
	case *Array:
		objs = arr.Value
	default:
		return reflect.Value{}, assignError(name, "")
	}

	count := len(objs)

	if err := compareArrayLen(count); err != nil {
		return reflect.Value{}, err
	}

	if count < 1 {
		return reflect.Zero(rtype), nil
	}

	v := reflect.MakeSlice(rtype, count, count)
	for i := 0; i < count; i++ {
		if elem, err := AssignValue(sliceElemType, objs[i]); err != nil {
			return reflect.Value{}, err
		} else {
			v.Index(i).Set(elem)
		}
	}

	return v, nil
}

func AssignInterface(rtype reflect.Type, o Object) (reflect.Value, error) {
	switch rtype.Kind() {
	case reflect.Interface:
	default:
		return reflect.Value{}, assignError("interface", rtype.Kind().String())
	}

	if isNil(o) {
		return reflect.Zero(rtype), nil
	}

	if ref, ok := o.(Reference); ok {
		refType := ref.RefType()
		if refType.AssignableTo(rtype) {
			return ref.RefValue(), nil
		}
		return AssignInterface(rtype, ref.Into())
	}

	if rtype.Implements(objectType) {
		objValue := reflect.ValueOf(o)
		if rtype.AssignableTo(objValue.Type()) {
			return objValue, nil
		}
	}

	return reflect.Value{}, assignError("interface", "")
}

func AssignPointer(rtype reflect.Type, o Object) (reflect.Value, error) {
	switch rtype.Kind() {
	case reflect.Pointer:
	default:
		return reflect.Value{}, assignError("pointer", rtype.Kind().String())
	}

	if isNil(o) {
		return reflect.Zero(rtype), nil
	}

	if ref, ok := o.(Reference); ok {
		refType := ref.RefType()
		if refType.AssignableTo(rtype) {
			return ref.RefValue(), nil
		}
		return AssignPointer(rtype, ref.Into())
	}

	if rtype.Implements(objectType) {
		objValue := reflect.ValueOf(o)
		if rtype.AssignableTo(objValue.Type()) {
			return objValue, nil
		}
	}

	v, err := AssignValue(rtype.Elem(), o)
	if err != nil {
		return reflect.Value{}, err
	}

	i := reflect.New(rtype).Elem()
	i.Set(v)

	return i, nil
}

func AssignMap(rtype reflect.Type, o Object) (reflect.Value, error) {
	switch rtype.Kind() {
	case reflect.Map:
	default:
		return reflect.Value{}, assignError("map", rtype.Kind().String())
	}

	if isNil(o) {
		return reflect.Zero(rtype), nil
	}

	if ref, ok := o.(Reference); ok {
		refType := ref.RefType()
		if refType.AssignableTo(rtype) {
			return ref.RefValue(), nil
		}
		return AssignMap(rtype, ref.Into())
	}

	var m map[string]Object

	switch mObj := o.(type) {
	case *ImmutableMap:
		m = mObj.Value
	case *Map:
		m = mObj.Value
	default:
		return reflect.Value{}, assignError("map", "")
	}

	kType, vType := rtype.Key(), rtype.Elem()

	if kType.Kind() != reflect.String {
		return reflect.Value{}, assignError("map", "")
	}

	v := reflect.MakeMapWithSize(rtype, len(m))
	for key, value := range m {
		vObj, err := AssignValue(vType, value)
		if err != nil {
			return reflect.Value{}, err
		}
		v.SetMapIndex(reflect.ValueOf(key), vObj)
	}

	return v, nil
}

func AssignStruct(rtype reflect.Type, o Object) (reflect.Value, error) {
	switch rtype.Kind() {
	case reflect.Struct:
	default:
		return reflect.Value{}, assignError("struct", rtype.Kind().String())
	}

	if isNil(o) {
		return reflect.Zero(rtype), nil
	}

	var objFields map[string]Object

	switch obj := o.(type) {
	case *Error:
		return AssignValue(rtype, obj.Value)
	case *Time:
		if rtype.AssignableTo(timeType) {
			timeValue := reflect.New(rtype).Elem()
			timeValue.Set(reflect.ValueOf(obj.Value))
			return timeValue, nil
		}
		return reflect.Value{}, assignError("time.Time", "")
	case *ImmutableMap:
		objFields = obj.Value
	case *Map:
		objFields = obj.Value
	}

	sFields := getProperties(rtype)

	v := reflect.New(rtype).Elem()
	for key, value := range objFields {
		if index, ok := sFields.fields[key]; ok {
			field := v.Field(index)
			fieldValue, err := AssignValue(field.Type(), value)
			if err != nil {
				return reflect.Value{}, err
			}
			field.Set(fieldValue)
		}
	}

	return v, nil
}

func AssignFunc(rtype reflect.Type, o Object) (reflect.Value, error) {
	switch rtype.Kind() {
	case reflect.Func:
	default:
		return reflect.Value{}, assignError("func", rtype.Kind().String())
	}

	if isNil(o) {
		return reflect.Zero(rtype), nil
	}

	if ref, ok := o.(Reference); ok {
		refType := ref.RefType()
		if refType.AssignableTo(rtype) {
			return ref.RefValue(), nil
		}
		return AssignFunc(rtype, ref.Into())
	}

	return reflect.Value{}, assignError("func", "")
}

func AssignValue(rtype reflect.Type, o Object) (reflect.Value, error) {
	switch rtype.Kind() {
	case reflect.Bool:
		return AssignBool(rtype, o)
	case reflect.Int:
		return AssignInt(rtype, o)
	case reflect.Int8:
		return AssignInt8(rtype, o)
	case reflect.Int16:
		return AssignInt16(rtype, o)
	case reflect.Int32:
		return AssignInt32(rtype, o)
	case reflect.Int64:
		return AssignInt64(rtype, o)
	case reflect.Uint:
		return AssignUint(rtype, o)
	case reflect.Uint8:
		return AssignUint8(rtype, o)
	case reflect.Uint16:
		return AssignUint16(rtype, o)
	case reflect.Uint32:
		return AssignUint32(rtype, o)
	case reflect.Uint64:
		return AssignUint64(rtype, o)
	case reflect.Float32:
		return AssignFloat32(rtype, o)
	case reflect.Float64:
		return AssignFloat64(rtype, o)
	case reflect.String:
		return AssignString(rtype, o)
	case reflect.Array, reflect.Slice:
		return AssignArray(rtype, o)
	case reflect.Func:
		return AssignFunc(rtype, o)
	case reflect.Interface:
		return AssignInterface(rtype, o)
	case reflect.Pointer:
		return AssignPointer(rtype, o)
	case reflect.Map:
		return AssignMap(rtype, o)
	case reflect.Struct:
		return AssignStruct(rtype, o)
	}

	return reflect.Value{}, fmt.Errorf("cannot assign to type: %s -> %s",
		o.TypeName(),
		rtype.Kind().String())
}

func fromRefValue(typeName string,
	expectType reflect.Type,
	v reflect.Value,
	convertFunc func(rtype reflect.Type, v reflect.Value) (Object, error)) (Object, error) {

	rtype := v.Type()

	if expectType != nil {
		if rtype.Kind() != expectType.Kind() {
			return nil, convertError(typeName, rtype.Kind().String())
		}
	}

	switch rtype.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.UnsafePointer:
		if v.IsNil() {
			return UndefinedValue, nil
		}
	}

	into, err := convertFunc(rtype, v)
	if err != nil {
		return nil, err
	}

	if expectType != nil {
		if rtype == expectType {
			return into, nil
		}
	}

	props, err := getProperties(rtype).Instance(rtype, v)
	if err != nil {
		return nil, err
	}

	return &reference{Object: into, props: props, rtype: rtype, value: v}, nil
}

func FromBool(v reflect.Value) (Object, error) {
	return fromRefValue("bool", boolType, v,
		func(p reflect.Type, v reflect.Value) (Object, error) {
			if v.Bool() {
				return TrueValue, nil
			}
			return FalseValue, nil
		})
}

func FromInt(v reflect.Value) (Object, error) {
	return fromRefValue("int", intType, v,
		func(p reflect.Type, v reflect.Value) (Object, error) {
			return &Number{Value: v.Int(), Type: Int}, nil
		})
}

func FromInt8(v reflect.Value) (Object, error) {
	return fromRefValue("int8", int8Type, v,
		func(p reflect.Type, v reflect.Value) (Object, error) {
			return &Number{Value: v.Int(), Type: Int8}, nil
		})
}

func FromInt16(v reflect.Value) (Object, error) {
	return fromRefValue("int16", int16Type, v,
		func(p reflect.Type, v reflect.Value) (Object, error) {
			return &Number{Value: v.Int(), Type: Int16}, nil
		})
}

func FromInt32(v reflect.Value) (Object, error) {
	return fromRefValue("int32", int32Type, v,
		func(p reflect.Type, v reflect.Value) (Object, error) {
			return &Number{Value: v.Int(), Type: Int32}, nil
		})
}

func FromInt64(v reflect.Value) (Object, error) {
	return fromRefValue("int64", int64Type, v,
		func(p reflect.Type, v reflect.Value) (Object, error) {
			return &Number{Value: v.Int(), Type: Int64}, nil
		})
}

func FromUint(v reflect.Value) (Object, error) {
	return fromRefValue("uint", uintType, v,
		func(rtype reflect.Type, v reflect.Value) (Object, error) {
			return &Number{Value: int64(v.Uint()), Type: Uint}, nil
		})
}

func FromUint8(v reflect.Value) (Object, error) {
	return fromRefValue("uint8", uint8Type, v,
		func(rtype reflect.Type, v reflect.Value) (Object, error) {
			return &Number{Value: int64(v.Uint()), Type: Uint8}, nil
		})
}

func FromUint16(v reflect.Value) (Object, error) {
	return fromRefValue("uint16", uint16Type, v,
		func(rtype reflect.Type, v reflect.Value) (Object, error) {
			return &Number{Value: int64(v.Uint()), Type: Uint16}, nil
		})
}

func FromUint32(v reflect.Value) (Object, error) {
	return fromRefValue("uint32", uint32Type, v,
		func(rtype reflect.Type, v reflect.Value) (Object, error) {
			return &Number{Value: int64(v.Uint()), Type: Uint32}, nil
		})
}

func FromUint64(v reflect.Value) (Object, error) {
	return fromRefValue("uint64", uint64Type, v,
		func(rtype reflect.Type, v reflect.Value) (Object, error) {
			return &Number{Value: int64(v.Uint()), Type: Uint64}, nil
		})
}

func FromFloat32(v reflect.Value) (Object, error) {
	return fromRefValue("float32", float32Type, v,
		func(rtype reflect.Type, v reflect.Value) (Object, error) {
			return &Float{Value: v.Float(), Type: Float32}, nil
		})
}

func FromFloat64(v reflect.Value) (Object, error) {
	return fromRefValue("float64", float64Type, v,
		func(rtype reflect.Type, v reflect.Value) (Object, error) {
			return &Float{Value: v.Float(), Type: Float64}, nil
		})
}

func FromString(v reflect.Value) (Object, error) {
	return fromRefValue("string", stringType, v,
		func(rtype reflect.Type, v reflect.Value) (Object, error) {
			str := v.String()
			if len(str) > MaxStringLen {
				return nil, ErrStringLimit
			}
			return &String{Value: str}, nil
		})
}

func FromFunc(name string, v reflect.Value) (Object, error) {
	return fromRefValue("func", callableType, v,
		func(rtype reflect.Type, v reflect.Value) (Object, error) {
			if rtype == callableType {
				return &UserFunction{Name: name, Value: v.Interface().(CallableFunc)}, nil
			}
			bridgeFunc := func(args ...Object) (ret Object, err error) {
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
					o := out[i]
					if obj, err = FromValue(o); err != nil {
						return
					}
					objs = append(objs, obj)
				}

				if numOut > 1 {
					return &ImmutableArray{Value: objs}, nil
				}

				return obj, nil
			}
			return &UserFunction{Name: name, Value: bridgeFunc}, nil
		})
}

func FromArray(v reflect.Value) (Object, error) {
	return fromRefValue("array", nil, v,
		func(rtype reflect.Type, v reflect.Value) (Object, error) {
			switch rtype.Kind() {
			case reflect.Array, reflect.Slice:
			default:
				return nil, fmt.Errorf("cannot convert to array: %s", rtype.Kind().String())
			}

			if rtype == bytesType {
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
		})
}

func FromMap(v reflect.Value) (Object, error) {
	return fromRefValue("map", nil, v,
		func(rtype reflect.Type, v reflect.Value) (Object, error) {
			switch rtype.Kind() {
			case reflect.Map:
			default:
				return nil, fmt.Errorf("cannot convert to map: %s", rtype.Kind().String())
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
		})
}

func FromStruct(v reflect.Value) (Object, error) {
	rtype := v.Type()

	switch rtype.Kind() {
	case reflect.Struct:
	default:
		return nil, fmt.Errorf("cannot convert to struct: %s", rtype.Kind().String())
	}

	if rtype == timeType {
		return &Time{Value: v.Interface().(time.Time)}, nil
	}

	if rtype.Implements(objectType) {
		return v.Interface().(Object), nil
	}

	m, err := getProperties(rtype).Instance(rtype, v)
	if err != nil {
		return nil, err
	}

	into := &ImmutableMap{Value: m}

	if rtype.Implements(errorType) {
		return &Error{Value: into}, nil
	}

	return into, nil
}

func FromPointer(v reflect.Value) (Object, error) {
	rtype := v.Type()
	if v.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("cannot convert to pointer: *%s", rtype.Elem().Name())
	}

	if v.IsNil() {
		return UndefinedValue, nil
	}

	var (
		into  Object
		err   error
		props map[string]Object
	)

	if rtype.Implements(objectType) {
		into = v.Interface().(Object)
		if into == UndefinedValue {
			return into, nil
		}
	} else {
		if into, err = FromValue(v.Elem()); err != nil {
			return nil, err
		}
		if props, err = getProperties(rtype).Instance(rtype, v); err != nil {
			return nil, err
		}
	}

	ptr := &reference{Object: into, props: props, rtype: rtype, value: v}

	if rtype.Implements(errorType) {
		return &Error{Value: ptr}, nil
	}

	return ptr, nil
}

func FromInterfaceValue(v reflect.Value) (Object, error) {
	rtype := v.Type()

	if rtype.Kind() != reflect.Interface {
		return nil, fmt.Errorf("cannot convert to interface: %s", rtype.Kind().String())
	}

	if v.IsNil() {
		return UndefinedValue, nil
	}

	var (
		into  Object
		err   error
		props map[string]Object
	)

	if rtype.Implements(objectType) {
		into = v.Interface().(Object)
		if into == UndefinedValue {
			return into, nil
		}
	} else if rtype == interfaceType {
		if into, err = FromValue(v.Elem()); err != nil {
			return nil, err
		}
	} else {
		if props, err = getProperties(rtype).Instance(rtype, v); err != nil {
			return nil, err
		}
	}

	return &reference{Object: into, props: props, rtype: rtype, value: v}, nil
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

	switch v.Kind() {
	case reflect.Bool:
		return FromBool(v)
	case reflect.Int:
		return FromInt(v)
	case reflect.Int8:
		return FromInt8(v)
	case reflect.Int16:
		return FromInt16(v)
	case reflect.Int32:
		return FromInt32(v)
	case reflect.Int64:
		return FromInt64(v)
	case reflect.Uint:
		return FromUint(v)
	case reflect.Uint8:
		return FromUint8(v)
	case reflect.Uint16:
		return FromUint16(v)
	case reflect.Uint32:
		return FromUint32(v)
	case reflect.Uint64:
		return FromUint64(v)
	case reflect.Float32:
		return FromFloat32(v)
	case reflect.Float64:
		return FromFloat64(v)
	case reflect.String:
		return FromString(v)
	case reflect.Array, reflect.Slice:
		return FromArray(v)
	case reflect.Func:
		return FromFunc("", v)
	case reflect.Interface:
		return FromInterfaceValue(v)
	case reflect.Pointer:
		return FromPointer(v)
	case reflect.Struct:
		return FromStruct(v)
	case reflect.Map:
		return FromMap(v)
	}

	return nil, fmt.Errorf("cannot convert to Object: %s", v.Kind().String())
}

type Reference interface {
	Into() Object
	RefType() reflect.Type
	RefValue() reflect.Value
}

type reference struct {
	Object
	props map[string]Object
	rtype reflect.Type
	value reflect.Value
}

func (o *reference) Equals(x Object) bool {
	if r, ok := x.(*reference); ok {
		if r.rtype.AssignableTo(o.rtype) {
			return r.value.Equal(o.value)
		}
	}
	return o.Object.Equals(x)
}

func (o *reference) IndexGet(index Object) (value Object, err error) {
	if o.props != nil {
		if str, ok := index.(*String); ok {
			if value, ok = o.props[str.Value]; ok {
				return
			}
		}
	}
	if o.Object != nil {
		return o.Object.IndexGet(index)
	}
	return UndefinedValue, nil
}

func (o *reference) Import(moduleName string) (interface{}, error) {
	if o.props == nil {
		o.props = make(map[string]Object)
	}
	o.props["__module_name__"] = &String{Value: moduleName}
	return o, nil
}

func (o *reference) RefType() reflect.Type {
	return o.rtype
}

func (o *reference) RefValue() reflect.Value {
	return o.value
}

func (o *reference) Into() Object {
	return o.Object
}
