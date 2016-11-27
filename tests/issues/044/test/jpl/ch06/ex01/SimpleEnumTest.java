package jpl.ch06.ex01;

import static org.junit.Assert.*;
import org.junit.Test;

public class SimpleEnumTest {

    @Test
    public void test1 () {
        Day sunday = Day.SUNDAY;
        Day satueday = Day.STATUERDAY;

        assertNotEquals(sunday,satueday);
        assertEquals(Day.SUNDAY,sunday);
        assertEquals(Day.STATUERDAY,satueday);
        assertEquals(Day.values()[0],sunday);
        assertEquals(Day.values()[6],satueday);

    }

    @Test
    public void test2 () {
        Traffic red = Traffic.RED;
        Traffic green = Traffic.GREEN;

        assertNotEquals(red,green);
        assertEquals(Traffic.RED,red);
        assertEquals(Traffic.GREEN,green);
        assertEquals(Traffic.values()[0],red);
        assertEquals(Traffic.values()[2],green);

    }


}
