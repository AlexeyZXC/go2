package dupremover

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

type data struct {
	duplicates []string
	items      map[string]int64 // path and size
	wg         sync.WaitGroup
	ch         chan struct{}
	dirCount   uint64
	filesCount uint64
	dir        string
}

func New(dir string) data {
	return data{}
}

func (d *data) Process(dir string) ([]string, error) {

	d.items = make(map[string]int64)
	d.wg = sync.WaitGroup{}
	d.ch = make(chan struct{}, 1)

	var err error

	d.wg.Add(1)
	go d.walk(dir, err)
	d.ch <- struct{}{}

	d.wg.Wait()

	return d.duplicates, err
}

func (d *data) Remove() (errs []error) {
	for _, v := range d.duplicates {
		if err := os.Remove(v); err != nil {
			errs = append(errs, err)
		}
	}
	return
}

func (d *data) walk(dir string, err error) {
	f, _ := os.Open(dir)
	defer d.wg.Done()
	atomic.AddUint64(&d.dirCount, 1)

	var (
		listSize int64
		info     fs.FileInfo
		listName string
		list     []fs.DirEntry
	)

	list, err = f.ReadDir(0)
	f.Close()
	if err != nil {

		return
	}

	for _, v := range list {
		if info, err = v.Info(); err != nil {
			fmt.Println("Error on info: ", err)
			continue
		}
		if info == nil {
			fmt.Println("Error on info is nil")
			continue
		}
		listSize = info.Size()
		listName = info.Name()

		if v.IsDir() {
			//dir
			d.wg.Add(1)
			go d.walk(dir+"\\"+listName, err)
		} else {
			// file
			<-d.ch
			d.filesCount++

			for fullName, itemSize := range d.items {
				sl := strings.Split(fullName, "\\")
				itemName := sl[len(sl)-1]
				if itemSize == listSize && itemName == listName {
					d.duplicates = append(d.duplicates, dir+"\\"+listName)
				}
			}
			d.items[dir+"\\"+listName] = listSize
			d.ch <- struct{}{}
		}
	}
}
