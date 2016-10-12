/*
 * Copyright (C) 2016 Yoshik Shibata. All rights reserved.
 */
package jp.ne.sonet.ca2.yshibata.test;

import org.junit.Test;
import static org.junit.Assert.*;

/**
 * Unit Test code for StdoutCapture
 *
 * @author yoshiki
 */
public class StdoutCaptureTest {

    @Test
    public void illegalStart() {
        StdoutCapture sc = new StdoutCapture();
        sc.start();

        try {
            sc.start();
        } catch (IllegalStateException e) {
            sc.stop();
            return;
        }
        fail();
    }

    @Test(expected = IllegalStateException.class)
    public void illegalStop() {
        StdoutCapture sc = new StdoutCapture();
        sc.stop();
    }

    @Test(expected = IllegalStateException.class)
    public void illegalStop2() {
        StdoutCapture sc = new StdoutCapture();
        sc.start();
        sc.stop();
        sc.stop();
    }

    @Test
    public void illegalAssertEquals() {
        StdoutCapture sc = new StdoutCapture();
        sc.start();
        try {
            sc.assertEquals("");
        } catch (IllegalStateException e) {
            sc.stop();
            return;
        }
        fail();
    }

    @Test
    public void oneLine() {
        StdoutCapture sc = new StdoutCapture();
        sc.start();
        System.out.println("Hello");
        sc.stop();
        sc.assertEquals("Hello");
    }

    @Test
    public void twoLines() {
        StdoutCapture sc = new StdoutCapture();
        sc.start();
        System.out.println("Hello");
        System.out.println("World");
        sc.stop();
        sc.assertEquals("Hello", "World");
    }
}
