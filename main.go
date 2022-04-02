// Package contains code for lesson1 and lesson2.
//
// And here is some bug? The page doesn't contain docs for the exported functions.
//
//Thoughts?
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

// NewErr creates and returns myError derived from error.
func NewErr(_err interface{}) error {
	//BUG not displayed comment by godoc
	var msg string
	if e, ok := _err.(error); ok {
		msg = e.Error()
	}
	return &myError{
		errorTime: time.Now(),
		errorMsg:  msg,
	}
}

// Error implements Error() of error.
func (myerr *myError) Error() string {
	return fmt.Sprintf("error time: %v; \n msg: %v ", myerr.errorTime, myerr.errorMsg)
}

// AccessFile work with file and close it on return using defer.
func AccessFile() {
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
	}
	fmt.Println("WriteString ok")

}

func main() {

	AccessFile()

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
