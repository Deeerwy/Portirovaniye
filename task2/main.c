/*
 * Задание 2: POSIX-системные вызовы open(), read(), write(), close()
 * Открытие файла, чтение содержимого, вывод на экран.
 * Совместимо с MINIX 3 / POSIX.1
 *
 * Намеренно используем низкоуровневые POSIX-вызовы вместо fopen/fread,
 * чтобы продемонстрировать работу с файловыми дескрипторами.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>   /* read(), write(), close() */
#include <fcntl.h>    /* open(), O_RDONLY          */
#include <errno.h>    /* errno, strerror()          */

#define BUFFER_SIZE 256

int main(int argc, char *argv[])
{
    const char *filename;
    int         fd;
    ssize_t     bytes_read;
    char        buffer[BUFFER_SIZE];

    /* Проверка аргументов */
    if (argc < 2) {
        /* Имя файла по умолчанию, если аргумент не передан */
        filename = "test.txt";
        fprintf(stderr, "Usage: %s <filename>\n", argv[0]);
        fprintf(stderr, "No filename provided — using default: %s\n\n",
                filename);
    } else {
        filename = argv[1];
    }

    printf("=== Task 2: POSIX file read ===\n");
    printf("Opening file: %s\n\n", filename);

    /* open() — POSIX системный вызов, возвращает файловый дескриптор */
    fd = open(filename, O_RDONLY);
    if (fd == -1) {
        fprintf(stderr, "Error opening file '%s': %s\n",
                filename, strerror(errno));
        return EXIT_FAILURE;
    }

    printf("--- File contents ---\n");

    /*
     * Читаем файл порциями по BUFFER_SIZE байт.
     * write() в STDOUT_FILENO используем намеренно (POSIX-вызов),
     * а не printf, чтобы избежать буферизации stdio.
     */
    while ((bytes_read = read(fd, buffer, sizeof(buffer) - 1)) > 0) {
        buffer[bytes_read] = '\0'; /* null-terminate для безопасности */
        /* POSIX write: fd=1 (STDOUT_FILENO) */
        if (write(STDOUT_FILENO, buffer, (size_t)bytes_read) == -1) {
            fprintf(stderr, "Error writing to stdout: %s\n",
                    strerror(errno));
            close(fd);
            return EXIT_FAILURE;
        }
    }

    if (bytes_read == -1) {
        fprintf(stderr, "\nError reading file: %s\n", strerror(errno));
        close(fd);
        return EXIT_FAILURE;
    }

    printf("\n--- End of file ---\n");

    /* close() — освобождаем файловый дескриптор */
    if (close(fd) == -1) {
        fprintf(stderr, "Error closing file: %s\n", strerror(errno));
        return EXIT_FAILURE;
    }

    printf("File closed successfully. Descriptor released.\n");

    return EXIT_SUCCESS;
}
