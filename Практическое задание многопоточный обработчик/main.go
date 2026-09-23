package main

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

const workerCount = 5

// Job описывает задание для воркера.
type Job struct {
	ID  int
	URL string
}

// Result содержит результат обработки задания.
type Result struct {
	Job      Job
	Status   string
	Duration time.Duration
}

// fakeHTPRTequest имитирует выполнение HTTP-запроса.
func fakeHTTPRequest() {
	delay := time.Duration(300+rand.Intn(1700)) * time.Millisecond
	time.Sleep(delay)
}

// worker обрабатывает задания из jobs и отправляет результаты в results.
func worker(
	workerID int,
	jobs <-chan Job,
	results chan<- Result,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for job := range jobs {
		start := time.Now()

		fakeHTTPRequest()

		duration := time.Since(start)

		results <- Result{
			Job:      job,
			Status:   "обработан",
			Duration: duration,
		}

		fmt.Printf(
			"Воркер %d обработал URL %s за %v\n",
			workerID,
			job.URL,
			duration.Round(time.Millisecond),
		)
	}
}

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	urls := []string{
		"https://example.com",
		"https://google.com",
		"https://github.com",
		"https://go.dev",
		"https://stackoverflow.com",
		"https://wikipedia.org",
		"https://developer.mozilla.org",
		"https://pkg.go.dev",
		"https://reddit.com",
		"https://news.ycombinator.com",
		"https://duckduckgo.com",
		"https://httpbin.org",
	}

	// Буферизованные каналы позволяют отделить отправителя
	// заданий от воркеров и сборщика результатов.
	jobs := make(chan Job, len(urls))
	results := make(chan Result, len(urls))

	var workersWG sync.WaitGroup
	workersWG.Add(workerCount)

	// Fan-out: запускаем фиксированное количество воркеров.
	for workerID := 1; workerID <= workerCount; workerID++ {
		go worker(workerID, jobs, results, &workersWG)
	}

	// Отдельная горутина ждёт завершения всех воркеров
	// и после этого закрывает канал результатов.
	go func() {
		workersWG.Wait()
		close(results)
	}()

	// Отправляем задания в канал jobs.
	for id, url := range urls {
		jobs <- Job{
			ID:  id + 1,
			URL: url,
		}
	}

	// Закрываем jobs: после этого воркеры завершат цикл range.
	close(jobs)

	// Fan-in: собираем результаты от всех воркеров.
	allResults := make([]Result, 0, len(urls))

	for result := range results {
		allResults = append(allResults, result)
	}

	// Сортируем результаты по ID задания,
	// чтобы отчёт соответствовал исходному порядку URL.
	sort.Slice(allResults, func(i, j int) bool {
		return allResults[i].Job.ID < allResults[j].Job.ID
	})

	var totalDuration time.Duration
	successCount := 0

	fmt.Println()
	fmt.Println("========== ИТОГОВЫЙ ОТЧЁТ ==========")

	for _, result := range allResults {
		fmt.Printf(
			"%2d. %-35s | статус: %-10s | время: %v\n",
			result.Job.ID,
			result.Job.URL,
			result.Status,
			result.Duration.Round(time.Millisecond),
		)

		totalDuration += result.Duration

		if result.Status == "обработан" {
			successCount++
		}
	}

	var averageDuration time.Duration
	if len(allResults) > 0 {
		averageDuration = totalDuration / time.Duration(len(allResults))
	}

	fmt.Println("------------------------------------")
	fmt.Printf("Всего URL:          %d\n", len(urls))
	fmt.Printf("Обработано успешно: %d\n", successCount)
	fmt.Printf("Среднее время:      %v\n", averageDuration.Round(time.Millisecond))
	fmt.Printf("Общее время работы: %v\n", totalDuration.Round(time.Millisecond))
	fmt.Println("====================================")
}
