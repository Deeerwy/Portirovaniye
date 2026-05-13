# Портирование ПО.Рабочая тетрадь 1. Задания 1-2-3. Бурылин Д.С. ПИМО-01-25


## Описание

Практическая работа по теме **"Языковые конверсии и отладка"**.  
Реализация трансляторов исходного кода с Python на Java  
и кроссплатформенных приложений.


## Структура проекта

```
porting_lab1/
├── task1_translator.py       # Задание 1: простой транслятор Python → Java
├── task2_python/
│   └── factorial.py          # Задание 2: факториал на Python
├── task2_java/
│   └── Factorial.java        # Задание 2: факториал на Java
├── task3_math_translator.py  # Задание 3: транслятор мат. функций Python → Java
└── README.md
```



## Задания

### Задание 1 — Простой транслятор
Программа принимает фрагмент кода на Python и генерирует  
эквивалентный код на Java. Поддерживает:
- переменные (`int`, `double`, `String`)
- вывод через `print()`

### Задание 2 — Кроссплатформенное приложение
Вычисление факториала числа, реализованное на двух языках:
- `factorial.py` — версия на Python
- `Factorial.java` — эквивалентная версия на Java

### Задание 3 — Транслятор математических функций
Расширенный транслятор Python → Java с поддержкой:
- функций (`def` → `public static long`)
- циклов `for`
- условий `if/else`
- операций `*=`, `+=`, `**`
- функций `factorial` и `fibonacci`

---

## Запуск

### Требования
- Python 3.x
- Java JDK 11+

### Python

```bash
python3 task1_translator.py
python3 task2_python/factorial.py
python3 task3_math_translator.py
```

### Java

```bash
cd task2_java
javac Factorial.java
java Factorial
```
