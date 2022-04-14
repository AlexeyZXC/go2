package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/trace"
	"sync"
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

// Для программы, которая использует
// мьютекс для безопасного доступа к
// данным из нескольких потоков,
// выполните трассировку
func lesson6_1() {
	j := 0
	defer func() {
		fmt.Println("\n j= ", j)
	}()

	fo, err := os.Create("trace_6_1.out")
	if err != nil {
		panic(err)
	}
	defer func() {
		fo.Close()
	}()

	trace.Start(fo)
	defer trace.Stop()

	m := sync.Mutex{}
	wg := &sync.WaitGroup{}
	n := 100

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			m.Lock()
			defer m.Unlock()
			defer wg.Done()

			j++
			fmt.Printf("%v ", j)
		}()
	}

	wg.Wait()
}

// Написать многопоточную программу,
// в которой будет использоваться
// явный вызов планировщика
// runtime.Gosched. Запустите ее с
// GOMAXPROCS=1 и выполните
// трассировку
func lesson6_2() {
	fo, err := os.Create("trace_6_2.out")
	if err != nil {
		panic(err)
	}
	defer func() {
		fo.Close()
	}()

	trace.Start(fo)
	defer trace.Stop()

	runtime.GOMAXPROCS(1)

	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer runtime.Gosched()
			fmt.Println("hello")
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			defer runtime.Gosched()
			fmt.Println("world")
			wg.Done()
		}()
	}

	wg.Wait()
}

// Смоделируйте состояние гонки.
// Проверьте модель на наличие
// состояния гонки
func lesson6_3() {
	j := 0
	defer func() {
		fmt.Println("\n j= ", j)
	}()

	//m := sync.Mutex{}
	wg := &sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			//m.Lock()
			//defer m.Unlock()
			defer wg.Done()

			j++
			fmt.Printf("%v ", j)
		}()
	}

	wg.Wait()
}

func main() {

	lesson6_3()
	//lesson6_2()
	//lesson6_1()

	// accessFile()

	// defer func() {
	// 	if err := recover(); err != nil {
	// 		myerr := NewErr(err)
	// 		fmt.Println("!!! panic recovered-1, myErr: ", myerr)

	// 		var e error
	// 		e, _ = err.(error)
	// 		myerr2 := fmt.Errorf("%w; \ntime: %v", e, time.Now())
	// 		fmt.Println("!!! panic recovered-2, myErr2: ", myerr2)
	// 	}
	// }()

	// var a int
	// _ = a / a
}
