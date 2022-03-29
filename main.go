package main

import (
	"fmt"
	"os"
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

func main() {

	accessFile()

	defer func() {
		if err := recover(); err != nil {
			myerr := NewErr(err)
			fmt.Println("!!! panic recovered-1, myErr: ", myerr)

			var e error
			e, _ = err.(error)
			myerr2 := fmt.Errorf("%w; \ntime: %v", e, time.Now())
			fmt.Println("!!! panic recovered-2, myErr2: ", myerr2)
		}
	}()

	var a int
	_ = a / a
}
