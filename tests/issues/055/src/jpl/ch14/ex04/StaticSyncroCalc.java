package jpl.ch14.ex04;

public class StaticSyncroCalc {

	static int num = 0;
	
	public static int getNum() {
		return num;
	}

	public static void setNum(int num) {
		StaticSyncroCalc.num = num;
	}

	static synchronized public void  countUp(){
		num++;
		System.out.println(num);
	}
}
