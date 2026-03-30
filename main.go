package main

import (
	"flag"
	"fmt"
	"lab1-client/api"
	"lab1-client/formatter"
	"lab1-client/pipeline"
	"log"
)

func main() {
	var postCount int
	var useMock bool
	
	flag.IntVar(&postCount, "count", 5, "количество постов для загрузки")
	flag.BoolVar(&useMock, "mock", true, "использовать тестовые данные")
	flag.Parse()

	if postCount < 1 {
		log.Fatal("количество постов должно быть положительным")
	}

	client := api.NewClient("https://jsonplaceholder.typicode.com")

	var posts []api.Post
	var err error
	
	if useMock {
		fmt.Println(" Используются тестовые данные")
		posts, err = client.GetPostsMock(postCount)
		if err != nil {
			log.Fatalf("ошибка при получении тестовых данных: %v", err)
		}
	} else {
		fmt.Println(" Подключение к API...")
		posts, err = client.GetPosts(postCount)
		if err != nil {
			log.Fatalf("ошибка при получении постов: %v", err)
		}
	}

	fmt.Printf(" Получено %d постов\n", len(posts))

	// Создание пайплайна
	p := pipeline.NewPipeline()

	// Этап 1: Фильтрация длинных заголовков (максимум 50 символов)
	p.AddStage(pipeline.NewFilterStage(50))

	// Этап 2: Fan-out/Fan-in с 3 воркерами
	p.AddStage(pipeline.NewContentProcessorStage(3))

	// Обработка постов
	fmt.Println(" Обработка через пайплайн...")
	processedPosts, err := p.Process(posts)
	if err != nil {
		log.Fatalf("ошибка при обработке пайплайна: %v", err)
	}

	// Форматирование
	formatter := formatter.NewJSONFormatter()
	formattedOutput, err := formatter.Format(processedPosts)
	if err != nil {
		log.Fatalf("ошибка при форматировании: %v", err)
	}

	fmt.Println("\n=== РЕЗУЛЬТАТ ОБРАБОТКИ ===")
	fmt.Println(formattedOutput)
	
	fmt.Printf("\n Статистика: обработано %d из %d постов\n", len(processedPosts), len(posts))
}