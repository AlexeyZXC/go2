package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"sync"
	"sync/atomic"
)

// Практическое задание
// В качестве завершающего задания нужно выполнить программу поиска дубликатов файлов.
// Дубликаты файлов - это файлы, которые совпадают по имени файла и по его размеру.
// Нужно написать консольную программу, которая проверяет наличие дублирующихся
// файлов.
// Программа должна работать на локальном компьютере и получать на вход путь до
// директории. Программа должна вывести в стандартный поток вывода список дублирующихся
// файлов, которые находятся как в директории, так и в поддиректориях директории,
// переданной через аргумент командной строки. Данная функция должна работать
// эффективно при помощи распараллеливания программы
// Программа должна принимать дополнительный ключ - возможность удаления обнаруженных
// дубликатов файлов после поиска. Дополнительно нужно придумать, как обезопасить
// пользователей от случайного удаления файлов. В качестве ключей желательно
// придерживаться общепринятых практик по использованию командных опций.
// Критерии приемки программы:
// 1. Программа компилируется
// 2. Программа выполняет функциональность, описанную выше.
// 3. Программа покрыта тестами
// 4. Программа содержит документацию и примеры использования
// 5. Программа обладает флагом “-h/--help” для краткого объяснения функциональности
// © geekbrains.ru 19
// 6. Программа должна уведомлять пользователя об ошибках, возникающих во время
// выполнения

var (
	duplicates []string
	items      map[string]uint64 // path and size
	wg         sync.WaitGroup
	ch         chan struct{}
	dirCount   uint64
)

func walk(dir string) error {
	f, _ := os.Open(dir)
	defer wg.Done()
	atomic.AddUint64(&dirCount, 1)

	list, err := f.ReadDir(0)
	f.Close()
	if err != nil {
		return err
	}

	var (
		listSize int64
		info     fs.FileInfo
	)

	for _, v := range list {
		if info, err = v.Info(); err != nil {
			fmt.Println("Error on info: ", err)
		}
		if info == nil {
			continue
		}
		listSize = info.Size()

		//fmt.Printf("%v - %v - %v; path: %v\n", v.Name(), v.IsDir(), listSize, dir+"\\"+v.Name())

		if v.IsDir() {
			//dir
			wg.Add(1)
			go walk(dir + "\\" + v.Name())
		} else {
			// file
			<-ch

			if itemSize, ok := items[v.Name()]; ok {
				if itemSize == uint64(listSize) {
					duplicates = append(duplicates, dir+"\\"+v.Name())
				}
			}
			items[v.Name()] = uint64(listSize)

			ch <- struct{}{}
		}
	}

	return nil
}

func main() {

	dir := flag.String("dir", "", "A directory to process")
	removeDup := flag.Bool("rem", false, "True value is about to remove duplicate files")
	flag.Parse()

	fmt.Println("removeDup: ", *removeDup)

	defer func() {
		if err := recover(); err != nil {
			switch v := err.(type) {
			case string, error:
				fmt.Println("Panic happened: ", v)
			default:
				fmt.Println("Panic happened: ", v)
			}
		}

	}()

	if *dir == "" {
		*dir, _ = os.Getwd()
	}

	if _, err := os.Stat(*dir); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Error: directory does not exist: ", *dir)
			return
		}
	}

	//*dir = "C:\\gb\\go2\\test"
	//*dir = "C:\\Users\\sakharov\\go\\src\\go2\\test"

	fmt.Println("dir: ", *dir)

	// arguments ok now

	items = make(map[string]uint64)
	wg = sync.WaitGroup{}
	ch = make(chan struct{}, 1)

	wg.Add(1)

	go walk(*dir)

	//fmt.Println("ch to send")
	ch <- struct{}{}
	//fmt.Println("ch sent")

	wg.Wait()

	fmt.Println("--- duplicates:")
	for _, s := range duplicates {
		fmt.Println(s)
	}

	// fmt.Println("--- items:")
	// for v, i := range items {
	// 	fmt.Printf("%v - %v\n", v, i)
	// }

	fmt.Println("dirCount: ", dirCount)

	// go func() {
	// 	<-ch
	// }()
	// ch <- struct{}{}

	// err := filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
	// 	fmt.Printf("%v - %v - %v\n", path, info.IsDir(), info.Size())
	// 	if v, ok := items[path]; ok {
	// 		if v == uint64(info.Size()) {
	// 			duplicates = append(duplicates, path)
	// 			return nil
	// 		}
	// 	}
	// 	items[path] = uint64(info.Size())
	// 	return nil
	// })
	// if err != nil {
	// 	panic(err)
	// }

}
