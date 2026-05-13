// Задание 1: Портирование файловых операций
//
// Аналогия с C:
//   CreateFile()  → os.OpenFile()
//   ReadFile()    → file.Read() / os.ReadFile()
//   CloseHandle() → file.Close()  (+ defer для гарантии)
//
// В Go управление ресурсами через defer — идиоматично и безопасно.
// Нет риска забыть закрыть файл даже при панике.

package main

import (
	"fmt"
	"os"
)

const testFilePath = "/tmp/lab2_task1_go.txt"

// createTestFile — аналог CreateFile() с флагом GENERIC_WRITE
func createTestFile(path string) error {
	// os.OpenFile — прямой аналог CreateFile():
	//   os.O_WRONLY  — только запись
	//   os.O_CREATE  — создать, если не существует
	//   os.O_TRUNC   — обнулить содержимое
	//   0600         — права: rw-------
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("ошибка открытия файла для записи: %w", err)
	}
	// defer гарантирует закрытие файла при выходе из функции —
	// аналог CloseHandle() в блоке finally
	defer file.Close()

	content := "Hello from Go porting demo!\nLine 2\nLine 3\n"
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("ошибка записи в файл: %w", err)
	}

	fmt.Printf("[OK] Файл создан: %s\n", path)
	return nil
}

// readFileLineByLine — аналог ReadFile() с построчным чтением
func readFile(path string) error {
	// Вариант 1: os.ReadFile — читает весь файл сразу (высокоуровнево)
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла: %w", err)
	}

	fmt.Printf("[OK] Файл прочитан: %s\n", path)
	fmt.Println("── Содержимое файла ──────────────────────")
	fmt.Print(string(data))
	fmt.Println("──────────────────────────────────────────")
	return nil
}

// readFileManual — аналог ReadFile() с ручным буфером (как в C)
func readFileManual(path string) error {
	// Вариант 2: ручное открытие + Read() — ближе к низкоуровневому C
	file, err := os.Open(path) // O_RDONLY
	if err != nil {
		return fmt.Errorf("ошибка открытия файла: %w", err)
	}
	defer file.Close() // CloseHandle() — гарантированно через defer

	buf := make([]byte, 64) // буфер как в C: char buffer[64]
	fmt.Println("── Чтение вручную (буфер 64 байт) ───────")
	for {
		n, err := file.Read(buf)
		if n > 0 {
			fmt.Print(string(buf[:n]))
		}
		if err != nil {
			// io.EOF — нормальное завершение, не ошибка
			break
		}
	}
	fmt.Println("\n──────────────────────────────────────────")
	return nil
}

func main() {
	fmt.Println("=== Задание 1: Портирование файловых операций → Go ===")
	fmt.Println()

	// Шаг 1: Создание файла
	fmt.Println("[*] Создаём тестовый файл...")
	if err := createTestFile(testFilePath); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %v\n", err)
		os.Exit(1)
	}

	// Шаг 2: Чтение через os.ReadFile (высокоуровневый способ)
	fmt.Println("\n[*] Читаем файл через os.ReadFile()...")
	if err := readFile(testFilePath); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %v\n", err)
		os.Exit(1)
	}

	// Шаг 3: Чтение через ручной буфер (низкоуровневый способ)
	fmt.Println("\n[*] Читаем файл вручную через file.Read()...")
	if err := readFileManual(testFilePath); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n[OK] Задание 1 выполнено успешно.")
}
