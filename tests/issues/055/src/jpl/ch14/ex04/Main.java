 package jpl.ch14.ex04;

public class Main {

	public static void main(String[] args) {
		for(int i = 0; i < 100; i++){
			StaticCountThread staticCountThread = new StaticCountThread();
			Thread thread = new Thread(staticCountThread);
			thread.start();
		}

	}

}
