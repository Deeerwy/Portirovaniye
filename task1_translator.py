import re

def translate_python_to_java(python_code: str) -> str:
    lines = python_code.strip().split("\n")
    java_lines = []
    indent = "    "
    indent_level = 2

    java_lines.append("public class Translated {")
    java_lines.append(f"{indent}public static void main(String[] args) {{")

    for line in lines:
        stripped = line.strip()

        if not stripped or stripped.startswith("#"):
            continue

        # print(...)
        m = re.match(r'print\((.+)\)', stripped)
        if m:
            content = m.group(1)
            java_lines.append(
                f"{indent * indent_level}System.out.println({content});"
            )
            continue

        # x = <value>
        m = re.match(r'(\w+)\s*=\s*(.+)', stripped)
        if m:
            var, val = m.group(1), m.group(2)
            if re.match(r'^\d+\.\d+$', val):
                java_lines.append(
                    f"{indent * indent_level}double {var} = {val};"
                )
            elif re.match(r'^\d+$', val):
                java_lines.append(
                    f"{indent * indent_level}int {var} = {val};"
                )
            elif val.startswith('"') or val.startswith("'"):
                val = val.replace("'", '"')
                java_lines.append(
                    f"{indent * indent_level}String {var} = {val};"
                )
            else:
                java_lines.append(
                    f"{indent * indent_level}var {var} = {val};"
                )
            continue

        java_lines.append(
            f"{indent * indent_level}// TODO: translate -> {stripped}"
        )

    java_lines.append(f"{indent}}}")
    java_lines.append("}")

    return "\n".join(java_lines)


# ── Точка входа ──────────────────────────────────────────────
if __name__ == "__main__":
    python_snippet = """
x = 10
y = 3.14
name = 'Hello'
print(x)
print(name)
"""
    result = translate_python_to_java(python_snippet)
    print("=" * 50)
    print("  Результат трансляции Python → Java")
    print("=" * 50)
    print(result)
