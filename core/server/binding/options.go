package binding

import (
	"strings"
	//"conan/ecode"
	"conan/utils"
)

type CheckHandler func([]string) error

var (
	supportOp = map[string]CheckHandler{
		"required": checkRequire,
	}
)

//
func splitNameAndOption(info string) (string, option) {

	str := strings.Split(info, ",")
	if len(str) <= 1 {
		return str[0], nil
	}
	name := str[0]
	strLen := len(str)

	op := make(map[string]struct{}, strLen-1)

	for i := 1; i < strLen; i++ {
		if utils.RemoveSpace(str[i]) == "" {
			continue
		}
		op[str[i]] = struct{}{}
	}

	return name, op
}

func CheckOptions(op option, val []string) error {
	for k, v := range supportOp {
		_, ok := op[k]
		if ok {
			if err := v(val); err != nil {
				return err
			}
		}
	}

	return nil
}

func checkRequire(value []string) error {
	if value == nil || len(value) <= 0 || utils.RemoveSpace(value[0]) == "" {
		//return ecode.ParamFail
		return nil
	}
	return nil
}
