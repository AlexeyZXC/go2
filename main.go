package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"strings"
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
	I int
	S string
}

func Assign(in *In, m map[string]interface{}) (err error) {
	defer func() {
		if err2 := recover(); err2 != nil {
			switch x := err2.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				// Fallback err (per specs, error strings should be lowercase w/o punctuation
				err = errors.New("unknown panic")
			}
		}
	}()

	inV := reflect.ValueOf(in)

	for v, inter := range m {
		switch vv := inter.(type) {
		case int:
			inV.Elem().FieldByName(strings.ToUpper(v)).SetInt(int64(vv))
		case string:
			inV.Elem().FieldByName(strings.ToUpper(v)).SetString(vv)
		}
	}

	return
}

// Написать функцию, которая принимает на вход имя файла и название функции. Необходимо
// подсчитать в этой функции количество вызовов асинхронных функций. Результат работы
// должен возвращать количество вызовов int и ошибку error. Разрешается использовать только
// go/parser, go/ast и go/token.

func findAsyncFuncs(file, funcName string) (int, error) {

	fset := token.NewFileSet()

	astFile, err := parser.ParseFile(fset, file, nil, parser.AllErrors)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	var goNum int

	for _, f := range astFile.Decls {
		fn, ok := f.(*ast.FuncDecl)
		if !ok {
			continue
		}

		if fn.Name.Name != funcName {
			continue
		}

		for _, st := range fn.Body.List {
			switch st.(type) {
			case *ast.GoStmt:
				goNum++
			}
		}
	}

	return goNum, nil
}

func main() {

	fmt.Println("--- Task 1 ---")

	in := In{I: 1, S: "asd"}
	m := map[string]interface{}{"i": 5, "s": "qwe"}

	fmt.Println("before in: ", in)
	fmt.Println("m: ", m)

	Assign(&in, m)

	fmt.Println(" after in: ", in)

	// task 2

	fmt.Println("--- Task 2 ---")

	n, err := findAsyncFuncs("file1.txt", "bar2")
	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println("Async funcs: ", n)

}
