public class Factorial {

    public static long factorial(int n) {
        if (n < 0) {
            throw new IllegalArgumentException(
                "Число должно быть неотрицательным"
            );
        }
        if (n == 0) {
            return 1;
        }
        long result = 1;
        for (int i = 1; i <= n; i++) {
            result *= i;
        }
        return result;
    }

    public static void main(String[] args) {
        int[] testValues = {0, 1, 5, 6, 10};
        System.out.println("========================================");
        System.out.println("  Факториал (Java)");
        System.out.println("========================================");
        for (int number : testValues) {
            System.out.println(
                "  factorial(" + number + ") = " + factorial(number)
            );
        }
    }
}
