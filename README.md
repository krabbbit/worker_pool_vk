# Тестовое задание на стажировку ВК

**Условие**

Реализовать примитивный worker-pool с возможностью динамически добавлять и удалять Воркеры. 
Входные данные (строки) поступают в канал, воркеры их обрабатывают (например, выводят на экран номер воркера и сами данные). 
Задание на базовые знания каналов и горутин.

# Пример 
``` go
func ExampleWorkerPool() {
	// Create a new worker pool with initial capacity of 5 jobs
	wp := workerpool.New(5)

	// Add 5 workers to the pool (total now: 5)
	for i := 0; i < 5; i++ {
		_ = wp.AddWorker()
	}

	// Add 4 jobs to be processed by the worker pool
	for i := 0; i < 4; i++ {
		_ = wp.AddJob(strconv.Itoa(i))
	}

	// Remove 3 workers from the pool (total now: 2)
	for i := 0; i < 3; i++ {
		_ = wp.Delete()
	}

	// Add 4 more jobs to the pool
	for i := 0; i < 4; i++ {
		_ = wp.AddJob(strconv.Itoa(i + 100))
	}

	// Shutdown the worker pool, gracefully terminating all workers
	// after they finish their current jobs
	wp.Shutdown()

	// Attempt to add a new worker after shutdown - this will fail
	err := wp.AddWorker()
	if err != nil {
		// Print error message if worker cannot be added (expected after shutdown)
		fmt.Println("wow, error ", err.Error())
	}
}
```

# Реализовано:
- Архитектура по `golang-standards/project-layout`
- Реализация интерфейса WorkerPool: добавление/удаление воркеров, добавление данных в канал
- Graceful shutdown 
- Тесты с coverage 97.8
- Документация проекта
- `cmd/main.go` с возможностью запуска из консоли в интерактивном режиме

# Todo:
- ~~`cmd/main.go` с возможностью запуска из консоли в интерактивном режиме~~
- ~~Документация проекта~~
- Прикрутить GitHub Actions
- Добавить логирование 
- Свои обёртки для ошибок

# Тестирование
coverage: 97.8% of statements
 - Добавление/удаление воркеров
 - Обработка задач при конкурентном доступе
 - Граничные случаи
 - Механизмы остановки


