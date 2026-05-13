# Makefile для Task 1
# MINIX 3 использует cc (clang или gcc-совместимый front-end)

CC      = cc
CFLAGS  = -Wall -Wextra -pedantic -std=c99 -D_POSIX_C_SOURCE=200112L
TARGET  = task1
SRC     = main.c

.PHONY: all clean run

all: $(TARGET)

$(TARGET): $(SRC)
	$(CC) $(CFLAGS) -o $(TARGET) $(SRC)

run: all
	./$(TARGET)

clean:
	rm -f $(TARGET)
