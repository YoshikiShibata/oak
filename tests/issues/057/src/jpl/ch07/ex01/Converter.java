package jpl.ch07.ex01;

public class Converter {

	public static void main(String[] args) {
		System.out.println(toUnicode("package jpl.ch07.ex01;"));
		System.out.println(toUnicode("public class HelloWorld {"));
		System.out.println(toUnicode("  public static void main(String[] args) {"));
		System.out.println(toUnicode("    System.out.println(\"Hello, World.\");"));
		System.out.println(toUnicode("  }"));
		System.out.println(toUnicode("}"));
	}
	
	public static String toUnicode(String text) {
		StringBuilder sb = new StringBuilder();
		for (int i = 0; i < text.length(); i++) {
			sb.append(String.format("\\u%04x", Character.codePointAt(text, i)));
		}
		return sb.toString();
	}

}
