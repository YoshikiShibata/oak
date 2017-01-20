package ch05.ex07;

import java.time.LocalDateTime;

// 予定を保持するクラス
public class TimeInterval {

	private LocalDateTime start = null;
	private LocalDateTime end = null;

	public TimeInterval(LocalDateTime start, LocalDateTime end) {
		if(start == null || end == null) {
			throw new IllegalArgumentException("予定開始・終了日時を指定してください");
		}
		this.start = start;
		this.end = end;
	}

	/**
	 * 予定時刻範囲を設定します.
	 * 既存の値を変更したくない場合はnullを入れてください.
	 * @param start 予定開始時刻
	 * @param end 予定終了時刻
	 */
	public void setTime(LocalDateTime start, LocalDateTime end) {
		if(start != null) {
			this.start = start;
		}
		if(end != null) {
			this.end = end;
		}
	}

	public LocalDateTime getStart() {
		return start;
	}

	public LocalDateTime getEnd() {
		return end;
	}

	/**
	 * 予定が被っていないかチェックします
	 * @param otherSchedule 比較対象の予定
	 * @return true:予定が被っている、false:予定が被っていない
	 */
	public boolean isConflict(TimeInterval otherSchedule) {
		boolean ret = false;
		LocalDateTime s = otherSchedule.getStart();
		LocalDateTime e = otherSchedule.getEnd();
		if(s.isBefore(start)) {
			ret = !(e.isBefore(start) || e.equals(start));
		} else {
			ret = !(end.isBefore(s) || end.equals(s));
		}
		return ret;
	}
}

