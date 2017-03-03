package jpl.ch10.ex01;

import org.junit.Test;

public class SpecialCharaConverterTest {

	@Test
	public void testConvert_lf() {
		String input = "lf\nlf";
		System.out.println("input : " + input);
		SpecialCharaConverter obj = new SpecialCharaConverter();
		String output = obj.convert(input);
		System.out.println("output: " + output);
	}
	@Test
	public void testConvert_tab() {
		String input = "tab\ttab";
		System.out.println("input : " + input);
		SpecialCharaConverter obj = new SpecialCharaConverter();
		String output = obj.convert(input);
		System.out.println("output: " + output);
	}
	@Test
	public void testConvert_backspace() {
		String input = "backspace\bbackspace";
		System.out.println("input : " + input);
		SpecialCharaConverter obj = new SpecialCharaConverter();
		String output = obj.convert(input);
		System.out.println("output: " + output);
	}
	@Test
	public void testConvert_ret() {
		String input = "ret\rret";
		System.out.println("input : " + input);
		SpecialCharaConverter obj = new SpecialCharaConverter();
		String output = obj.convert(input);
		System.out.println("output: " + output);
	}
	@Test
	public void testConvert_ff() {
		String input = "ff\fff";
		System.out.println("input : " + input);
		SpecialCharaConverter obj = new SpecialCharaConverter();
		String output = obj.convert(input);
		System.out.println("output: " + output);
	}
	@Test
	public void testConvert_backslash() {
		String input = "backslash\\backslash";
		System.out.println("input : " + input);
		SpecialCharaConverter obj = new SpecialCharaConverter();
		String output = obj.convert(input);
		System.out.println("output: " + output);
	}
	@Test
	public void testConvert_sq() {
		String input = "sq\'sq";
		System.out.println("input : " + input);
		SpecialCharaConverter obj = new SpecialCharaConverter();
		String output = obj.convert(input);
		System.out.println("output: " + output);
	}
	@Test
	public void testConvert_dq() {
		String input = "dq\"dq";
		System.out.println("input : " + input);
		SpecialCharaConverter obj = new SpecialCharaConverter();
		String output = obj.convert(input);
		System.out.println("output: " + output);
	}

}
