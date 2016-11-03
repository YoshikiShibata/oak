package jpl.ch02.ex17;

/**
 *  2つのturnメソッドを追加.
 *  1:引数として回転する角度を受け取る
 *  2:Vehicle.TURM_LEFTかTURN_RIGHTを受けとる.
 * @author Anna.S
 *
 */
public class Vehicle {
	private double speed;
	private double direction;
	private String name;

	private static long nextID = 0;
	private final long idNum;

	static final int TURN_LEFT = 0;
	static final int TURN_RIGHT = 1;

	private int way = 1;

	public Vehicle() {
		idNum = nextID++;
	}

	public Vehicle(String name) {
		this();
		this.name = name;
	}

	public double getSpeed() {
		return speed;
	}

	public void setSpeed(double speed) {
		this.speed = speed;
	}

	public double getDirection() {
		return direction;
	}

	public void setDirection(double direction) {
		this.direction = direction;
	}

	public String getName() {
		return name;
	}

	public void setName(String driver) {
		this.name = driver;
	}

	public long getIdNum() {
		return idNum;
	}

	/** 回転する角度を取得.*/
	public void turn(double direction) {
		if (way == TURN_LEFT)
			this.direction -= direction + 180;
		if (way == TURN_RIGHT)
			this.direction += direction;
		while (direction >= 360) {
			this.direction = direction - 360;
		}
	}

	/** 通常は右回り.*/
	public void turn(int turnWay) {
		this.way = turnWay;
	}

	@Override
	public String toString() {
		StringBuilder sb = new StringBuilder();
		sb.append("Vehicle[ID:");
		sb.append(idNum);
		sb.append(", speed:");
		sb.append(speed);
		sb.append(", direction:");
		sb.append(direction);
		sb.append(", name:");
		sb.append(name);
		sb.append("]");
		return sb.toString();
	}


}