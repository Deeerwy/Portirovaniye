# Makefile для Task 2

CC      = cc
CFLAGS  = -Wall -Wextra -pedantic -std=c99 -D_POSIX_C_SOURCE=200112L
TARGET  = task2
SRC     = main.c
TESTFILE = test.txt

.PHONY: all clean run run-custom test-missing

all: $(TARGET)

$(TARGET): $(SRC)
	$(CC) $(CFLAGS) -o $(TARGET) $(SRC)

# Запуск с файлом по умолчанию (test.txt)
run: all
	./$(TARGET) $(TESTFILE)

# Запуск с кастомным файлом: make run-custom FILE=myfile.txt
run-custom: all
	./$(TARGET) $(FILE)

# Тест: несуществующий файл (проверка обработки ошибок)
test-missing: all
	./$(TARGET) nonexistent_file.txt || true

clean:
	rm -f $(TARGET)
