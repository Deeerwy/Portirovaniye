import re


def infer_java_type(value: str) -> str:
    """Определяет Java-тип по строковому значению."""
    if re.match(r'^\d+\.\d+$', value):
        return "double"
    if re.match(r'^\d+$', value):
        return "int"
    return "var"


def replace_pow(expr: str) -> str:
    """Заменяет оператор ** на Math.pow()."""
    return re.sub(
        r'(\w+)\s*\*\*\s*(\w+)',
        r'(long)Math.pow(\1, \2)',
        expr
    )


def translate_math_functions(python_code: str) -> str:
    """Транслирует простые математические функции Python → Java."""
    java_lines = []
    indent = "    "

    java_lines.append("public class MathTranslated {")
    java_lines.append("")

    lines = python_code.strip().split("\n")
    i = 0

    while i < len(lines):
        line = lines[i]
        stripped = line.strip()

        # Определение функции: def name(param):
        m = re.match(r'def (\w+)\((\w*)\):', stripped)
        if m:
            func_name, param = m.group(1), m.group(2)
            param_str = f"int {param}" if param else ""
            java_lines.append(
                f"{indent}public static long {func_name}({param_str}) {{"
            )
            i += 1

            # Собираем тело функции
            while i < len(lines):
                body_raw = lines[i]
                body = body_raw.strip()

                # Конец тела функции — пустая строка или новая def
                if body == "" or re.match(r'def \w+', body):
                    break

                if body.startswith("#"):
                    i += 1
                    continue

                # return <expr>
                m_ret = re.match(r'return (.+)', body)
                if m_ret:
                    expr = replace_pow(m_ret.group(1))
                    java_lines.append(f"{indent}{indent}return {expr};")
                    i += 1
                    continue

                # if <cond>:
                m_if = re.match(r'if (.+):', body)
                if m_if:
                    cond = m_if.group(1)
                    cond = cond.replace(" and ", " && ")
                    cond = cond.replace(" or ", " || ")
                    cond = cond.replace(" not ", " !")
                    java_lines.append(f"{indent}{indent}if ({cond}) {{")
                    i += 1
                    continue

                # else:
                if body == "else:":
                    java_lines.append(f"{indent}{indent}}} else {{")
                    i += 1
                    continue

                # for i in range(a, b + 1):
                m_for = re.match(
                    r'for (\w+) in range\((\w+),\s*(\w+)\s*\+?\s*(\d*)\):',
                    body
                )
                if m_for:
                    var   = m_for.group(1)
                    start = m_for.group(2)
                    end   = m_for.group(3)
                    plus  = m_for.group(4)
                    limit = f"{end} + {plus}" if plus else end
                    java_lines.append(
                        f"{indent}{indent}"
                        f"for (int {var} = {start}; {var} <= {limit}; {var}++) {{"
                    )
                    i += 1
                    continue

                # var *= value
                m_mul = re.match(r'(\w+)\s*\*=\s*(.+)', body)
                if m_mul:
                    vname, val = m_mul.group(1), m_mul.group(2)
                    java_lines.append(
                        f"{indent}{indent}{indent}{vname} *= {val};"
                    )
                    i += 1
                    continue

                # var += value
                m_add = re.match(r'(\w+)\s*\+=\s*(.+)', body)
                if m_add:
                    vname, val = m_add.group(1), m_add.group(2)
                    java_lines.append(
                        f"{indent}{indent}{indent}{vname} += {val};"
                    )
                    i += 1
                    continue

                # var = value  (объявление)
                m_assign = re.match(r'(\w+)\s*=\s*(.+)', body)
                if m_assign:
                    var, val = m_assign.group(1), m_assign.group(2)
                    jtype = infer_java_type(val)
                    java_lines.append(
                        f"{indent}{indent}{jtype} {var} = {val};"
                    )
                    i += 1
                    continue

                # Нераспознанная строка
                java_lines.append(f"{indent}{indent}// TODO: {body}")
                i += 1

            java_lines.append(f"{indent}}}")
            java_lines.append("")
            continue

        i += 1

    java_lines.append("}")
    return "\n".join(java_lines)


# ── Точка входа ──────────────────────────────────────────────
if __name__ == "__main__":
    python_math_code = """
def factorial(n):
    if n == 0:
        return 1
    result = 1
    for i in range(1, n + 1):
        result *= i
    return result

def fibonacci(n):
    if n <= 1:
        return n
    a = 0
    b = 1
    for i in range(2, n + 1):
        c = a + b
        a = b
        b = c
    return b
"""

    print("=" * 55)
    print("  Транслятор математических функций Python → Java")
    print("=" * 55)
    output = translate_math_functions(python_math_code)
    print(output)
