package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"sync"
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

// type Item struct{
// 	Name string
// 	size uint64
// 	dublicate bool
// }

var (
	duplicates []string
	items      map[string]uint64 // path and size
	wg         sync.WaitGroup
)

func walk(dir string) error {
	f, _ := os.Open(dir)
	//defer f.Close()
	defer wg.Done()

	list, err := f.ReadDir(0)
	f.Close()
	if err != nil {
		return err
	}

	var (
		size int64
		info fs.FileInfo
	)

	for _, v := range list {
		if info, err = v.Info(); err != nil {
			fmt.Println("Error on info: ", err)
		}
		if info == nil {
			continue
		}
		size = info.Size()

		fmt.Printf("%v - %v - %v; path: %v\n", v.Name(), v.IsDir(), size, dir+"\\"+v.Name())

		if v.IsDir() {
			wg.Add(1)
			go walk(dir + "\\" + v.Name())
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

	*dir = "C:\\gb\\go2\\test"

	fmt.Println("dir: ", *dir)

	// arguments ok now

	items = make(map[string]uint64)
	wg = sync.WaitGroup{}

	wg.Add(1)
	go walk(*dir)

	wg.Wait()

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
