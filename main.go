package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const numCPU = 4

const (
	BYTE = 1 << (10 * iota)
	KILOBYTE
	MEGABYTE
	GIGABYTE
	TERABYTE
)

// подсчёт оптимального отображения размерности файла
func Size(bytes uint64) string {
	unit := ""
	value := float64(bytes)

	switch {
	case bytes >= TERABYTE:
		unit = "T"
		value = value / TERABYTE
	case bytes >= GIGABYTE:
		unit = "G"
		value = value / GIGABYTE
	case bytes >= MEGABYTE:
		unit = "M"
		value = value / MEGABYTE
	case bytes >= KILOBYTE:
		unit = "K"
		value = value / KILOBYTE
	case bytes >= BYTE:
		unit = "B"
	case bytes == 0:
		return "0B"
	}

	result := strconv.FormatFloat(value, 'f', 1, 64)
	result = strings.TrimSuffix(result, ".0")
	return result + unit
}

// фильтрация файлов, обход директорий и вывод в консоль дерева каталога
func walkDir(path string, step int) (err error) {
	all, err := ioutil.ReadDir(path)
	var (
		dirs  []fs.FileInfo
		files []fs.FileInfo
	)
	if err != nil {
		return err
	}
	// фильтрация файлов и папок
	for _, now := range all {
		if now.IsDir() {
			dirs = append(dirs, now)
		} else {
			files = append(files, now)
		}
	}
	// отображение папок
	for _, dir := range dirs {
		dirName := dir.Name()
		dirMode := dir.Mode().String()
		further := path + "/" + dirName
		// вывод двухстрочной информации о директории
		fmt.Println(strconv.Itoa(step)+":("+further+") |", dirMode, "\n├"+strings.Repeat("───", step), dirName)
		// step предназначен для числового отображения иерархии (назначается как литерал в main при первом вызове walkDir)
		step++
		// рекурсивный вызов след. директории (обход директорий)
		if err := walkDir(further, step); err != nil {
			return err
		}
		step--
	}
	// отображение файлов
	for _, file := range files {
		fmt.Println("│<>"+strings.Repeat("<><>", step), file.Name(), Size(uint64(file.Size())))
	}

	return
}

func main() {
	// фиксированное кол-во процессоров, которые могут выполняться одновременно
	runtime.GOMAXPROCS(numCPU)
	// вычисление времени (начало выполнения)
	start := time.Now()

	fmt.Println("Путь к исполняемому файлу | Кол-во аргументов (пути к директориям) | Директории")
	// получение аргументов консоли
	args := os.Args
	// срез всех путей директорий
	paths := args[2:]
	cntPaths, _ := strconv.Atoi(args[1])
	// проверка на кол-во заявленных директорий
	if len(paths) != cntPaths {
		panic("Используйте заявленное кол-во аругментов!")
	}
	// обход всех путей с помощью горутин
	go func() {
		for _, path := range paths {
			walkDir(path, 0)
		}
	}()

	// вычисление времени (конец выполнения)
	elapsedTime := time.Since(start)
	fmt.Println("Итоговое время выполнения: " + elapsedTime.String())
	// задерживаем выполнение главной горутины функции main (для успешного выполнения всех вложенных горутин)
	time.Sleep(time.Second)
}
