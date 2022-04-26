// Programme finds the file duplicates within the directory passed in -dir argument. The list of the found duplicate files is printed out to stdout.
// The duplicate files can be removed by specifying the -rem flag.
package main

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

import (
	"flag"
	"fmt"
	dupremover "go2/dupRemover"
	"os"
	"strings"
)

func main() {

	dir := flag.String("dir", "", "A directory to process")
	removeDup := flag.Bool("rem", false, "True value is about to remove duplicate files")
	h := flag.Bool("h", false, "To unveil the help")
	help := flag.Bool("help", false, "To unveil the help")
	flag.Parse()

	*help = *help || *h
	if *help {
		fmt.Printf("The tool finds and deletes the duplicates of files within the path passed in -dir argument.\n %v",
			"Duplicates of files are the files with the same name and size.\n Arguments: \n-dir(string)  directory to process. \n-rem(bool) remove duplicate files or do "+
				"not.\n Example: For removing duplicate files located in c:\\test provide the command: main.exe -dir=c:\\test -rem=true")
		return
	}

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

	// arguments ok now

	remover := dupremover.New(*dir)

	dups, err := remover.Process()
	if len(dups) == 0 {
		fmt.Println("No duplicates found")
		return
	}
	fmt.Println("Duplicates files:")
	for _, v := range dups {
		fmt.Println(v)
	}

	if err != nil {
		fmt.Println("Process error: ", err)
	}
	if *removeDup {
		var resp string = "no"
		fmt.Println("Do you want to remove the duplicate files?. No or yes?")
		fmt.Scan(&resp)
		resp = strings.ToLower(resp)
		if resp == "yes" {
			errs := remover.Remove()
			if len(errs) > 0 {
				fmt.Println("Errors:")
				for _, v := range errs {
					fmt.Println(v)
				}
			}
			fmt.Println("Done.")
		}
	}
}
