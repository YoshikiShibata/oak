/*
 * Copyright (C) 2014 Yoshiki Shibata. All rights reserved.
 */
package ch05.ex04;

import java.time.LocalDate;

/**
 * Write an equivalent of the Unix cal program that displays a calendar for a
 * month. For example, java Cal 3 2013 should display
 *
 * ある月のカレンダーを表示するUnix のcal プログラムと同じプログラムを書きなさい。 たとえば、java Cal 3 2013
 * を実行すると、次のように表示します。
 * <pre>
 *              1  2  3
 *  4  5  6  7  8  9 10
 * 11 12 13 14 15 16 17
 * 18 19 20 21 22 23 24
 * 25 26 27 28 29 30 31
 * </pre>
 *
 * indicating that March 1 is a Friday. (Show the weekend at the end of the
 * week.)
 *
 * 3 月1 日が金曜日であることを示しています（土曜日と日曜が右端に表示されるようにし なさい）。
 */
public class Cal {

    public static void main(String[] args) {
        if (args.length != 2) {
            showUsage();
        }

        LocalDate currentDate = getFirstDateOfMonth(args);
        int dayOfWeek = leadingSpaces(currentDate);
        int dayOfMonth = 1;
        
        do {
            showDay(dayOfMonth, dayOfWeek);
            
            currentDate = currentDate.plusDays(1);
            dayOfWeek = currentDate.getDayOfWeek().getValue();
            dayOfMonth = currentDate.getDayOfMonth();
            
        } while (dayOfMonth != 1);
        System.out.println();

    }

    private static LocalDate getFirstDateOfMonth(String[] args)  {
        int month = Integer.parseInt(args[0]);
        int year = Integer.parseInt(args[1]);
        LocalDate firstDate = LocalDate.of(year, month, 1);
        return firstDate;
    }

    private static void showDay(int dayOfMonth, int dayOfWeek) {
        System.out.printf("%3d", dayOfMonth);
        
        if ((dayOfWeek % 7) == 0) {
            System.out.println();
        }
    }

    private static int leadingSpaces(LocalDate ld) {
        int dayOfWeek = ld.getDayOfWeek().getValue();
        for (int i = 1; i < dayOfWeek; i++) {
            System.out.print("   ");
        }
        return dayOfWeek;
    }

    private static void showUsage() {
        System.out.printf("usage: Cal month year%n"
                + " month: 1 - 12%n"
                + " year: -999999999 - +999999999%n");
        System.exit(1);
    }
}
