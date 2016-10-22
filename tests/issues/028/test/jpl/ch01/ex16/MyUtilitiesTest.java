package jpl.ch01.ex16;

public class MyUtilitiesTest {

	public static void main(String args[]) {
		MyUtilities myUtility = new MyUtilities();
		try {
			myUtility.getDataSet("test");
		} catch (BadDataSetException e) {
			System.err.println(e);
		}
	}
}
