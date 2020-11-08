package binding

import (
	"fmt"
	"testing"
	"conan/log"
)

type TestStruct struct {
	A string `form:"a_field,required"`
	B string `form:"b_field" default:"bbbb"`
}

type TestSli struct {
	A []int64 `form:"a_field,required,split"`
}

func TestDefaultBind(t *testing.T) {
	ts := TestStruct{}

	log.Debug("TestDefaultBind Bind Name is %s", formBind.Name())
	testForm := make(map[string][]string)
	testForm["a_field"] = []string{""}
	testForm["b_field"] = []string{""}
	err := formBind.testInterface(testForm, &ts)
	if err != nil {
		log.Error("testInterface error %s ", err.Error())
		return
	}

	log.Info("TestStruct content is %v ", ts)
}

func TestStructPoint(t *testing.T) {

	temp := make(map[string]struct{}, 0)

	temp["11"] = struct{}{}
	temp["22"] = struct{}{}
	temp["33"] = struct{}{}
	temp["44"] = struct{}{}

	for k, v := range temp {
		temp := v
		fmt.Println("k ", k)
		fmt.Printf("v %p \n", &temp)
	}
}

func TestBindSlice(t *testing.T) {

	ts := TestSli{}
	//
	//testForm := make(map[string][]string)
	//testForm["a_field"] = []string{"111", "222"}
	//testForm["b_field"] = []string{""}
	//err := formBind.testInterface(testForm, &ts)
	//
	//if err != nil {
	//	fmt.Println(err)
	//}

	testForm := make(map[string][]string)

	testForm["a_field"] = []string{}

	err := formBind.testInterface(testForm, &ts)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range ts.A {
		fmt.Println(v)
	}
}
