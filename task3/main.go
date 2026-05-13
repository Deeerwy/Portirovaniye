// Задание 3: Параллельные вычисления в Go
//
// Используем goroutines + channels — идиоматичный Go
// вместо pthreads из C-версии.
//
// Стратегия: делим массив на N частей,
// каждая горутина считает частичную сумму
// и отправляет результат в канал.
// Главная горутина собирает суммы из канала.

package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

const arraySize = 100_000_000 // 100 млн элементов

// ─────────────────────────────────────────────────────────────
// Вариант 1: goroutines + channel (классический Go-стиль)
// ─────────────────────────────────────────────────────────────

func parallelSumChannel(data []int, numWorkers int) int64 {
	chunkSize := len(data) / numWorkers
	results := make(chan int64, numWorkers) // буферизованный канал

	for i := 0; i < numWorkers; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if i == numWorkers-1 {
			end = len(data) // последний воркер берёт остаток
		}

		// Замыкание: передаём срез явно, чтобы избежать гонки
		go func(chunk []int, id int) {
			var sum int64
			for _, v := range chunk {
				sum += int64(v)
			}
			fmt.Printf("  [Горутина %d] диапазон [%d, %d), частичная сумма = %d\n",
				id, start, end, sum)
			results <- sum // отправляем результат в канал
		}(data[start:end], i)
	}

	// Собираем результаты из канала
	var total int64
	for i := 0; i < numWorkers; i++ {
		total += <-results
	}
	return total
}

// ─────────────────────────────────────────────────────────────
// Вариант 2: goroutines + sync.WaitGroup + атомарная запись
// (аналог pthread_join из C-версии)
// ─────────────────────────────────────────────────────────────

func parallelSumWaitGroup(data []int, numWorkers int) int64 {
	chunkSize := len(data) / numWorkers
	partials := make([]int64, numWorkers) // каждая горутина пишет в свой элемент
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		start := i * chunkSize
		end := start + chunkSize
		if i == numWorkers-1 {
			end = len(data)
		}

		go func(chunk []int, idx int) {
			defer wg.Done() // аналог выхода из pthread + автосигнал
			var sum int64
			for _, v := range chunk {
				sum += int64(v)
			}
			partials[idx] = sum // нет гонки: каждый индекс — только один горутин
		}(data[start:end], i)
	}

	wg.Wait() // аналог pthread_join для всех горутин сразу

	var total int64
	for _, p := range partials {
		total += p
	}
	return total
}

// ─────────────────────────────────────────────────────────────
// Однопоточный подсчёт (для сравнения производительности)
// ─────────────────────────────────────────────────────────────

func sequentialSum(data []int) int64 {
	var sum int64
	for _, v := range data {
		sum += int64(v)
	}
	return sum
}

// ─────────────────────────────────────────────────────────────
// Вспомогательная функция замера времени
// ─────────────────────────────────────────────────────────────

func measure(label string, fn func() int64) (int64, time.Duration) {
	start := time.Now()
	result := fn()
	elapsed := time.Since(start)
	fmt.Printf("  %-35s сумма = %d, время = %v\n", label+":", result, elapsed)
	return result, elapsed
}

func main() {
	// Определяем количество воркеров по числу логических ядер
	numWorkers := runtime.NumCPU()
	runtime.GOMAXPROCS(numWorkers) // разрешаем использовать все ядра

	fmt.Println("=== Задание 3: Параллельные вычисления (goroutines) ===")
	fmt.Printf("Размер массива  : %d элементов\n", arraySize)
	fmt.Printf("Кол-во горутин  : %d (по числу CPU)\n", numWorkers)
	fmt.Println()

	// ── 1. Инициализация массива ──────────────────────────────
	fmt.Println("[*] Инициализация массива...")
	data := make([]int, arraySize)
	for i := range data {
		data[i] = 1 // ожидаемая сумма == arraySize
	}
	fmt.Println("[OK] Задание 3 выполнено успешно.")

	// ── 2. Замеры ─────────────────────────────────────────────
	fmt.Println("--- Результаты ---")

	seqResult, seqTime := measure("Однопоточный (sequential)",
		func() int64 { return sequentialSum(data) })

	fmt.Println()

	chanResult, chanTime := measure("Параллельный (channel)",
		func() int64 { return parallelSumChannel(data, numWorkers) })

	fmt.Println()

	wgResult, wgTime := measure("Параллельный (WaitGroup)",
		func() int64 { return parallelSumWaitGroup(data, numWorkers) })

	// ── 3. Итог ───────────────────────────────────────────────
	fmt.Println("\n--- Сравнение ---")
	fmt.Printf("  Однопоточное время       : %v\n", seqTime)
	fmt.Printf("  Параллельное (channel)   : %v (ускорение %.2fx)\n",
		chanTime, float64(seqTime)/float64(chanTime))
	fmt.Printf("  Параллельное (WaitGroup) : %v (ускорение %.2fx)\n",
		wgTime, float64(seqTime)/float64(wgTime))

	allMatch := seqResult == chanResult && chanResult == wgResult
	fmt.Printf("  Все результаты совпадают : %v\n",
		map[bool]string{true: "ДА ✓", false: "НЕТ ✗"}[allMatch])

	fmt.Println("\n[OK] Задание 3 выполнено успешно.")
}
