package main

import (
	"fmt"
	"os"
	"reflect"
	"time"
)

type myError struct {
	errorTime time.Time
	errorMsg  string
}

func NewErr(_err interface{}) error {
	var msg string
	if e, ok := _err.(error); ok {
		msg = e.Error()
	}
	return &myError{
		errorTime: time.Now(),
		errorMsg:  msg,
	}
}

func (myerr *myError) Error() string {
	return fmt.Sprintf("error time: %v; \n msg: %v ", myerr.errorTime, myerr.errorMsg)
}

func accessFile() {
	var f *os.File
	var err error

	if f, err = os.OpenFile("file.txt", os.O_RDWR|os.O_CREATE, 0644); err != nil {
		fmt.Println("OpenFile error: ", err)
		return
	}

	defer func() {
		if err = f.Close(); err != nil {
			fmt.Println("Close error")
		} else {
			fmt.Println("Close ok")
		}
	}()

	if _, err := f.WriteString("content"); err != nil {
		fmt.Println("WriteString error")
		return
	} else {
		fmt.Println("WriteString ok")
	}
}

// Написать функцию, которая принимает на вход структуру in (struct или
// 	кастомную struct) и values map[string]interface{} (key - название поля
// 	структуру, которому нужно присвоить value этой мапки). Необходимо по
// 	значениям из мапы изменить входящую структуру in с помощью пакета
// 	reflect. Функция может возвращать только ошибку error. Написать к данной
// 	функции тесты (чем больше, тем лучше - зачтется в плюс).

type In struct {
	i int
	s string
	a [3]byte
}

func assign(in *In, m map[string]interface{}) error {
	inV := reflect.ValueOf(in)
	fmt.Println("inV: ", inV)
	fmt.Println("inV.Elem(): ", inV.Elem())
	fmt.Println("inV.Elem().FieldByIndex: ", inV.Elem().FieldByIndex([]int{1}))
	fmt.Println("inV.Elem().FieldByName: ", inV.Elem().FieldByName("s"))
	fmt.Println("inV.Elem().NumField(): ", inV.Elem().NumField())

	for i := 0; i < inV.Elem().NumField(); i++ {

	}

	//fmt.Println("inV.FieldByIndex: ", inV.FieldByIndex([]int{0}))
	// var (
	// 	v   reflect.Value
	// 	err error
	// )
	// if v, err = inV.FieldByIndexErr([]int{0}); err != nil {
	// 	fmt.Println("FieldByIndexErr err: ", err)
	// }

	// fmt.Println("inV.FieldByIndexErr: ", v)

	return nil
}

func main() {

	in := In{i: 1, s: "asd", a: [3]byte{1, 2, 3}}
	m := map[string]interface{}{"i": 5, "s": "qwe", "a": [3]byte{4, 5, 6}}

	fmt.Println("in: ", in)
	fmt.Println("m: ", m)

	assign(&in, m)

	fmt.Println("in: ", in)

}
