package binding

import (
	"github.com/pkg/errors"
	"reflect"
	"strconv"
	"strings"
	"conan/utils"
)

func mapForm(ptr interface{}, form map[string][]string) error {
	sInfo := sCache.get(reflect.TypeOf(ptr))
	val := reflect.ValueOf(ptr).Elem()
	var err error
	for i, fd := range sInfo.fields {
		structVal := val.Field(i)
		if !structVal.CanSet() {
			continue
		}

		if fd.name == "" {
			continue
		}
		formV, ok := form[fd.name]
		if !ok {
			if fd.hasDefault {
				structVal.Set(fd.defaultValue)
				continue
			}
		}

		// 检查Op
		if !fd.hasDefault {
			err = CheckOptions(fd.options, formV)
			if err != nil {
				return err
			}
		}

		if formV == nil || len(formV) <= 0 {
			continue
		}

		if formV[0] == "" && fd.hasDefault {
			structVal.Set(fd.defaultValue)
			continue
		}
		if err := setWithProperType(fd.tp.Type.Kind(), formV, structVal, fd.options); err != nil {
			return err
		}
	}
	return nil
}

func setWithProperType(kind reflect.Kind, value []string, dv reflect.Value, options option) error {

	switch kind {
	case reflect.Int:
		return setIntValue(value[0], 0, dv)
	case reflect.Int8:
		return setIntValue(value[0], 8, dv)
	case reflect.Int16:
		return setIntValue(value[0], 16, dv)
	case reflect.Int32:
		return setIntValue(value[0], 32, dv)
	case reflect.Int64:
		return setIntValue(value[0], 64, dv)
	case reflect.Uint:
		return setUintValue(value[0], 0, dv)
	case reflect.Uint8:
		return setUintValue(value[0], 8, dv)
	case reflect.Uint16:
		return setUintValue(value[0], 16, dv)
	case reflect.Uint32:
		return setUintValue(value[0], 32, dv)
	case reflect.Uint64:
		return setUintValue(value[0], 64, dv)
	case reflect.Float32:
		return setFloatValue(value[0], 32, dv)
	case reflect.Float64:
		return setFloatValue(value[0], 64, dv)
	case reflect.Bool:
		return setBoolValue(value[0], dv)
	case reflect.Slice:
		// 过滤空值
		//val := sliceEmpty(value)
		val := value
		_, ok := options["split"]
		if ok {
			val = strings.Split(value[0], ",")
		}

		val = sliceEmpty(val)

		slice := reflect.MakeSlice(dv.Type(), len(val), len(val))
		sliceKind := dv.Type().Elem().Kind()
		for i := 0; i < len(val); i++ {
			if err := setWithProperType(sliceKind, val[i:], slice.Index(i), options); err != nil {
				return err
			}
		}
		dv.Set(slice)
	default:
		dv.SetString(value[0])
	}
	return nil
}

func setSlice(value []string, field reflect.Value) error {

	valType := field.Type()

	val := reflect.MakeSlice(valType, 0, len(value))

	for _, v := range value {
		ele := reflect.New(valType).Elem()
		ele.Set(reflect.ValueOf(v))
		reflect.Append(val, ele)
	}

	field = val

	return nil
}

func setBoolValue(value string, field reflect.Value) error {
	if value == "" {
		value = "false"
	}

	boolVal, err := strconv.ParseBool(value)

	if err != nil {
		return err
	}

	field.SetBool(boolVal)
	return nil
}

func setFloatValue(value string, bit int, field reflect.Value) error {
	if value == "" {
		value = "0.0"
	}

	fValue, err := strconv.ParseFloat(value, bit)

	if err != nil {
		return nil
	}
	field.SetFloat(fValue)
	return nil
}

func setIntValue(value string, bit int, field reflect.Value) error {
	if value == "" {
		value = "0"
	}
	intV, err := strconv.ParseInt(value, 10, bit)

	if err != nil {
		return errors.WithStack(err)
	}

	field.SetInt(intV)
	return nil
}

func setUintValue(value string, bit int, dv reflect.Value) error {
	if value == "" {
		value = "0"
	}

	uintV, err := strconv.ParseUint(value, 10, bit)
	if err != nil {
		return err
	}

	dv.SetUint(uintV)
	return nil
}

func sliceEmpty(value []string) []string {

	res := make([]string, 0, len(value))

	for _, v := range value {
		str := utils.RemoveSpace(v)
		if str != "" {
			res = append(res, str)
		}
	}

	return res
}
