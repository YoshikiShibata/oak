import java.util.stream.Stream;

public class MyClass {

	public static void main(String[] args) {
		Stream<Long> st = random(3);
		st.forEach(System.out::println);	// ずっと書き出される（手動で止める）
	}

	static Stream<Long> random(long seed) {
		long a = 25214903917L;
		long c = 11;
		long m = (long) Math.pow(2, 48);
		return Stream.iterate(seed,  (x) -> (a * x + c) % m);
	}
}
