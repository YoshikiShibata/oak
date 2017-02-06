package jpl.ch14.ex04;

public class StaticCountThread implements Runnable {
	
	public StaticCountThread(){
		
	}

	@Override
	public void run() {
		for(int i = 0; i < 100; i++){
			StaticSyncroCalc.countUp();
		}
		
	}
	
	

}
