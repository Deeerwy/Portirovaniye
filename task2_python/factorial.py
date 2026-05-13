def factorial(n: int) -> int:
    if n < 0:
        raise ValueError("Число должно быть неотрицательным")
    if n == 0:
        return 1
    result = 1
    for i in range(1, n + 1):
        result *= i
    return result


# ── Точка входа ──────────────────────────────────────────────
if __name__ == "__main__":
    test_values = [0, 1, 5, 6, 10]
    print("=" * 40)
    print("  Факториал (Python)")
    print("=" * 40)
    for number in test_values:
        print(f"  factorial({number}) = {factorial(number)}")
