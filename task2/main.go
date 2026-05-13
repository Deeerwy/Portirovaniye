// Задание 2: Безопасная работа с указателями в Go
//
// В Go нет ручного malloc/free — память управляется GC.
// Тем не менее "пустые ссылки" существуют в виде nil-указателей,
// nil-интерфейсов, nil-срезов и nil-каналов.
//
// Лекционные правила переносятся в Go так:
//   1. Инициализация в nil    → var p *int = nil  (явно)
//   2. Проверка перед использованием → if p != nil { ... }
//   3. Очистка после free()   → p = nil после использования
//   4. Современный язык       → Go сам управляет памятью через GC
//   5. Статический анализ     → go vet, staticcheck
//   6. Динамическое тест-е    → go test -race (race detector)

package main

import (
	"errors"
	"fmt"
	"os"
)

// ─────────────────────────────────────────────────────────────
// 1. Работа с nil-указателем на примитив
// ─────────────────────────────────────────────────────────────

func demoNilPointer() {
	fmt.Println("\n--- 1. nil-указатель на int ---")

	// Правило 1: явная инициализация в nil
	var ptr *int = nil

	// Правило 2: проверка перед использованием
	if ptr == nil {
		fmt.Println("[INFO] ptr == nil, выделяем память...")
		value := 42
		ptr = &value // в Go нет malloc: берём адрес переменной
	}

	fmt.Printf("[OK] Значение: %d\n", *ptr)

	// Правило 4: GC освободит память сам,
	// но мы явно обнуляем, чтобы не использовать повторно
	ptr = nil
	if ptr == nil {
		fmt.Println("[OK] Указатель обнулён — повторное использование предотвращено.")
	}
}

// ─────────────────────────────────────────────────────────────
// 2. Безопасная работа со срезом (аналог динамического массива)
// ─────────────────────────────────────────────────────────────

func demoNilSlice() {
	fmt.Println("\n--- 2. nil-срез (аналог динамического массива) ---")

	// nil-срез: безопасно передавать, len/cap == 0
	var arr []int = nil
	fmt.Printf("[INFO] nil-срез: len=%d, cap=%d\n", len(arr), cap(arr))

	// Инициализируем
	const size = 5
	arr = make([]int, size)
	for i := range arr {
		arr[i] = i * 10
	}

	// Правило 5: проверка границ — в Go panic при выходе за границу,
	// поэтому всегда проверяем индекс вручную
	safeGet := func(s []int, index int) (int, error) {
		if index < 0 || index >= len(s) {
			return 0, fmt.Errorf("индекс %d за пределами массива (размер %d)",
				index, len(s))
		}
		return s[index], nil
	}

	if v, err := safeGet(arr, 3); err == nil {
		fmt.Printf("[OK] arr[3] = %d\n", v)
	}

	if _, err := safeGet(arr, 10); err != nil {
		fmt.Printf("[ЗАЩИТА] %v — обращение заблокировано.\n", err)
	}

	// Обнуляем срез (GC освободит backing array)
	arr = nil
	fmt.Println("[OK] Срез обнулён.")
}

// ─────────────────────────────────────────────────────────────
// 3. Управление ресурсом: безопасная работа с файлом
// ─────────────────────────────────────────────────────────────

func demoFileResource() error {
	fmt.Println("\n--- 3. Управление ресурсом FILE (defer) ---")

	path := "/tmp/lab2_task2_go_safe.txt"

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("os.Create: %w", err)
	}
	// defer — аналог гарантированного CloseHandle() / fclose()
	// выполняется даже при панике или раннем return
	defer func() {
		file.Close()
		fmt.Println("[OK] Файл закрыт через defer.")
	}()

	fmt.Printf("[OK] Файл открыт: %s\n", path)

	if _, err := file.WriteString("Safe pointer demo in Go\n"); err != nil {
		return fmt.Errorf("WriteString: %w", err)
	}

	return nil
}

// ─────────────────────────────────────────────────────────────
// 4. nil-интерфейс и ловушка «не-nil интерфейс с nil значением»
// ─────────────────────────────────────────────────────────────

type Animal interface {
	Sound() string
}

type Dog struct{ Name string }

func (d *Dog) Sound() string { return "Woof" }

// newAnimal возвращает nil-интерфейс корректно
func newAnimal(name string) Animal {
	if name == "" {
		return nil // правильно: nil-интерфейс
	}
	return &Dog{Name: name}
}

func demoNilInterface() {
	fmt.Println("\n--- 4. nil-интерфейс ---")

	a := newAnimal("")
	// Правило 2: проверяем перед вызовом метода
	if a == nil {
		fmt.Println("[ЗАЩИТА] Animal == nil — вызов метода заблокирован.")
	}

	a = newAnimal("Rex")
	if a != nil {
		fmt.Printf("[OK] Animal: %s\n", a.Sound())
	}
}

// ─────────────────────────────────────────────────────────────
// 5. Идиома Go: возврат ошибки вместо nil-разыменования
// ─────────────────────────────────────────────────────────────

type User struct {
	Name  string
	Email string
}

var ErrUserNotFound = errors.New("пользователь не найден")

// findUser — безопасно возвращает (*User, error) вместо просто *User
func findUser(id int) (*User, error) {
	db := map[int]*User{
		1: {Name: "Alice", Email: "alice@example.com"},
		2: {Name: "Bob", Email: "bob@example.com"},
	}

	user, ok := db[id]
	if !ok {
		return nil, ErrUserNotFound // nil + ошибка, не panic
	}
	return user, nil
}

func demoErrorHandling() {
	fmt.Println("\n--- 5. Идиома (value, error) вместо nil-разыменования ---")

	for _, id := range []int{1, 99} {
		user, err := findUser(id)
		if err != nil {
			// Правило 2: проверяем ошибку ДО использования указателя
			if errors.Is(err, ErrUserNotFound) {
				fmt.Printf("[ЗАЩИТА] ID %d: %v\n", id, err)
			}
			continue
		}
		// Здесь user гарантированно не nil
		fmt.Printf("[OK] Найден: %s (%s)\n", user.Name, user.Email)
	}
}

func main() {
	fmt.Println("=== Задание 2: Безопасная работа с указателями в Go ===")

	demoNilPointer()
	demoNilSlice()

	if err := demoFileResource(); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %v\n", err)
		os.Exit(1)
	}

	demoNilInterface()
	demoErrorHandling()

	fmt.Println("\n[OK] Задание 2 выполнено успешно.")
}
