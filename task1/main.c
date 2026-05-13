/*
 * Задание 1: POSIX-системный вызов getpid()
 * Получение текущего PID процесса и вывод на экран.
 * Совместимо с MINIX 3 / POSIX.1
 */

#include <stdio.h>
#include <unistd.h>   /* getpid(), getppid() — POSIX */

int main(void)
{
    pid_t pid  = getpid();   /* PID текущего процесса  */
    pid_t ppid = getppid();  /* PID родительского процесса */

    printf("=== Task 1: POSIX getpid() ===\n");
    printf("Current  PID  : %d\n", (int)pid);
    printf("Parent   PID  : %d\n", (int)ppid);

    return 0;
}
