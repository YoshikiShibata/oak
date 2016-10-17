package ch24.ex01;

import java.util.Locale;
import java.util.ResourceBundle;

public class GlobalHello {
	public static void main(String[] args) {
		ResourceBundle res = ResourceBundle.getBundle("ch24.ex01.GlobalRes");
		String msg;
		if (args.length > 0) {
			msg = res.getString(GlobalRes.GOODBYE);
		} else {
			msg = res.getString(GlobalRes.HELLO);
		}
		System.out.println(msg);

		Locale.setDefault(new Locale("ne","BT","Acme,Inc"));
		ResourceBundle res2 = ResourceBundle.getBundle("ch24.ex01.GlobalRes_eo_KI_left");
		if (args.length > 0) {
			msg = res2.getString(GlobalRes.GOODBYE);
		} else {
			msg = res2.getString(GlobalRes.HELLO);
		}
		System.out.println(msg);

	}
}
