package ch05.ex07;

import static org.hamcrest.CoreMatchers.*;
import static org.junit.Assert.*;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;

import org.junit.Test;

public class TimeIntervalTest {

	@Test
	public void testTimeInterval() {
		LocalDateTime start = LocalDateTime.of(LocalDate.now(), LocalTime.of(10, 0));
		LocalDateTime end = LocalDateTime.of(LocalDate.now(), LocalTime.of(11, 0));
		TimeInterval ti;

		// 第一引数がnullのケース
		try {
			ti = new TimeInterval(null, end);
			fail();
		} catch (Exception e) {
			assertThat(e.getMessage(), is("予定開始・終了日時を指定してください"));
		}

		// 第二引数がnullのケース
		try {
			ti = new TimeInterval(start, null);
			fail();
		} catch (Exception e) {
			assertThat(e.getMessage(), is("予定開始・終了日時を指定してください"));
		}

		ti = new TimeInterval(start, end);
		assertThat(ti.getStart(), is(start));
		assertThat(ti.getEnd(), is(end));
	}

	@Test
	public void testSetTime() {
			LocalDateTime start = LocalDateTime.of(LocalDate.now(), LocalTime.of(0, 0));
			LocalDateTime end = LocalDateTime.of(LocalDate.now(), LocalTime.of(0, 0));
			// どちらもセットされた場合
			TimeInterval ti = new TimeInterval(start, end);
			start = LocalDateTime.of(LocalDate.now(), LocalTime.of(10, 0));
			end = LocalDateTime.of(LocalDate.now(), LocalTime.of(11, 0));
			ti.setTime(start, end);
			assertThat(ti.getStart(), is(start));
			assertThat(ti.getEnd(), is(end));
			// 片方だけセットされた場合
			start = LocalDateTime.of(LocalDate.now(), LocalTime.of(9, 0));
			ti.setTime(start, null);
			assertThat(ti.getStart(), is(start));
			assertThat(ti.getEnd(), is(end));
			end = LocalDateTime.of(LocalDate.now(), LocalTime.of(12, 0));
			ti.setTime(null, end);
			assertThat(ti.getStart(), is(start));
			assertThat(ti.getEnd(), is(end));
		}

		@Test
		public void testGetStart() {
			LocalDateTime start = LocalDateTime.of(LocalDate.now(), LocalTime.of(10, 0));
			LocalDateTime end = LocalDateTime.of(LocalDate.now(), LocalTime.of(11, 0));
			TimeInterval ti = new TimeInterval(start, end);
			assertThat(ti.getStart(), is(start));
		}

		@Test
		public void testGetEnd() {
			LocalDateTime start = LocalDateTime.of(LocalDate.now(), LocalTime.of(10, 0));
			LocalDateTime end = LocalDateTime.of(LocalDate.now(), LocalTime.of(11, 0));
			TimeInterval ti = new TimeInterval(start, end);
			assertThat(ti.getEnd(), is(end));
		}

		@Test
		public void testIsConflict() {
			// 比較元TimeInterval
			LocalDateTime start = LocalDateTime.of(LocalDate.now(), LocalTime.of(10, 0));
			LocalDateTime end = LocalDateTime.of(LocalDate.now(), LocalTime.of(11, 0));
			TimeInterval ti = new TimeInterval(start, end);

			// 予定が被らないTimeInterval
			LocalDateTime startNotConfrict = LocalDateTime.of(LocalDate.now(), LocalTime.of(11, 0));
			LocalDateTime endNotConfrict = LocalDateTime.of(LocalDate.now(), LocalTime.of(12, 0));
			TimeInterval tiNotConfrict = new TimeInterval(startNotConfrict, endNotConfrict);
			assertThat(ti.isConflict(tiNotConfrict), is(false));

			// 予定が被るTimeInterval
			LocalDateTime startConfrict = LocalDateTime.of(LocalDate.now(), LocalTime.of(10, 30));
			LocalDateTime endConfrict = LocalDateTime.of(LocalDate.now(), LocalTime.of(12, 0));
			TimeInterval tiConfrict = new TimeInterval(startConfrict, endConfrict);
			assertThat(ti.isConflict(tiConfrict), is(true));

		}

	}
