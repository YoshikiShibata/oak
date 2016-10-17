package ch24.ex01;

import java.util.ListResourceBundle;

public class GlobalRes_eo_KI_left extends ListResourceBundle{
	public static final String HELLO = "hello";
	public static final String GOODBYE = "goodbye";

	private static final Object[][] contents = { { GlobalRes.HELLO, "Saluton" }, { GlobalRes.GOODBYE, "Ĝis revido" } };

	@Override
	protected Object[][] getContents() {
		// TODO 自動生成されたメソッド・スタブ
		return contents;
	}
}
