package jpl.ch10.ex01;

public class SpecialCharaConverter {
	public final String convert(String value) {
		String answer = "";
		for (int i=0; i<value.length(); i++) {
			String n = answer.concat(convertSpecialChara(value.charAt(i)));
			answer = n;
		}
		return answer;
	}

	private final String convertSpecialChara(char b) {
		String ans;
		if (b == '\n') {			//改行\u000A
			ans = "\\n";
		} else if (b == '\t') {		//タブ\u0009
			ans = "\\t";
		} else if (b == '\b') {		//バックスペース\u0008
			ans = "\\b";
		} else if (b == '\r') {		//復帰\u000D
			ans = "\\r";
		} else if (b == '\f') {		//フォームフィールド\u000C
			ans = "\\f";
		} else if (b == '\\') {		//バックスラッシュ自身\u005C
			ans = "\\\\";
		} else if (b == '\'') {		//シングルクォート\u0027
			ans = "\\\'";
		} else if (b == '\"') {		//ダブルクォート\u0022
			ans = "\\\"";
		} else {
			ans = String.valueOf(b);
		}
		return ans;
	}
}
