package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
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

type worker struct {
	pVlaue *int
	syncCh chan struct{}
}

func NewWorker(syncCh chan struct{}) worker {
	return worker{
		pVlaue: nil,
		syncCh: syncCh,
	}
}

func (w worker) StartWorker(v *int) {
	w.pVlaue = v
	go func() {
		<-w.syncCh
		*w.pVlaue++
	}()
}

func (w worker) Stop() {
	w.syncCh <- struct{}{}
}

func poolOfWorkers() {
	var syncCh = make(chan struct{}, 1)
	var value int = 0

	defer func() {
		fmt.Println("poolOfWorkers final value: ", value)
	}()

	syncCh <- struct{}{} // start the first worker

	for i := 0; i < 1000; i++ {
		w := NewWorker(syncCh)
		w.StartWorker(&value)
		w.Stop()
	}
}

func main() {

	// task 1
	poolOfWorkers()

	// task 2
	var cancel context.CancelFunc
	ctx := context.Background()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt) // there is no SIGTERM on Windows
	fmt.Println("waiting for Interrupt... ")

	s := <-ch

	fmt.Println("received Interrupt: ", s)
	ctx, cancel = context.WithTimeout(ctx, time.Millisecond)
	defer cancel()

	<-ctx.Done()

	err := ctx.Err()
	fmt.Println(errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded))

	//time.Sleep(time.Second)
	//fmt.Println("sleep done")

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
