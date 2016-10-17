package ch24.ex01;

import java.util.ListResourceBundle;

public class GlobalRes extends ListResourceBundle {
	public static final String HELLO = "hello";
	public static final String GOODBYE = "goodbye";

	private static final Object[][] contents = { { GlobalRes.HELLO, "Ciao" }, { GlobalRes.GOODBYE, "Ciao" } };

	@Override
	protected Object[][] getContents() {
		// TODO 自動生成されたメソッド・スタブ
		return contents;
	}

}
