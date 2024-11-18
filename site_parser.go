package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

func fetchAndSaveURL(wg *sync.WaitGroup, url string, folder string, index int) {
	defer wg.Done()

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Ошибка при скачивании:", err)
		return
	}
	defer resp.Body.Close()

	_, fileName := filepath.Split(url)
	// Если нет имени, то задаем базовое
	if fileName == "" {
		fileName = "index.html"
	}

	// Добавим сначала индекс, чтоб различать файлы
	filePath := filepath.Join(folder, fmt.Sprintf("%v %v", index, fileName))

	file, err := os.Create(filePath)

	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println("Ошибка при записи в файл:", err)
		return
	}

	fmt.Printf("Сайт %s сохранен в файл %s\n", url, filePath)
}

func main() {
	var wg sync.WaitGroup

	// Закидываем Урлы сайтов
	urls := []string{
		"https://vk.com/",
		"https://urfu.ru/ru/",
		"https://music.yandex.ru/",
	}

	folder, err := os.Getwd()

	if err != nil {
		fmt.Println("Ошибка при получении текущей рабочей директории:", err)
		return
	}

	for index, url := range urls {
		wg.Add(1)
		go fetchAndSaveURL(&wg, url, folder, index)
	}

	wg.Wait()
}
