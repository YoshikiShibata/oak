/*
 * Copyright (C) 2016 Yoshiki Shibata. All rights reserved.
 */
package jp.ne.sonet.ca2.yshibata.test;

import java.io.ByteArrayOutputStream;
import java.io.PrintStream;
import java.io.UnsupportedEncodingException;
import java.util.logging.Level;
import java.util.logging.Logger;
import junit.framework.AssertionFailedError;

import org.junit.Assert;
import org.junit.internal.ArrayComparisonFailure;

/**
 *
 * @author yoshiki
 */
public final class StdoutCapture {

    private boolean started = false;
    private PrintStream writer;
    private ByteArrayOutputStream baos;

    /**
     * Starts capturing of System.out
     *
     * @throws IllegalStateException Capturing has already started.
     */
    public synchronized void start() {
        if (started) {
            throw new IllegalStateException("Has already started");
        }

        baos = new ByteArrayOutputStream();
        PrintStream ps;
        try {
            ps = new PrintStream(baos, true, "UTF-8");
            writer = System.out;
            System.setOut(ps);
            started = true;
        } catch (UnsupportedEncodingException ex) {
            Logger.getLogger(StdoutCapture.class.getName()).log(Level.SEVERE, null, ex);
        }

    }

    /**
     * Stops capturing of System.out
     *
     * @throws IllegalStateException Capturing has not started yet.
     */
    public synchronized void stop() {
        if (!started) {
            throw new IllegalStateException("Has not started yet");
        }
        started = false;
        System.out.close();
        System.setOut(writer);
    }

    /**
     * Determines if the captured output equals to the specified argument.
     * All CR and LF characters are ignored.
     *
     * @param expected An array of expected output
     */
    public synchronized void assertEquals(String... expected) {
        if (started) {
            throw new IllegalStateException("Has not stopped yet");
        }
        
        String[] trimmedExpected = removeCRLF(expected);
        String[] trimmedResult = removeCRLF(toStringArray(baos.toByteArray()));
        
        try {
            Assert.assertArrayEquals(trimmedExpected,trimmedResult);
        } catch (ArrayComparisonFailure e) {
            System.err.printf("%nResult is%n");
            for (String s: trimmedResult) 
                System.err.printf("%s%n", s);
            
            System.err.printf("%nBut want is%n");
            for (String s: trimmedExpected)
                System.err.printf("%s%n", s);
            
            System.err.println();
            throw e;
        }
    }
    
    private String[] toStringArray(byte[] bytes) {
        try {
            String out = new String(bytes, "UTF-8");
            return out.split(System.lineSeparator());
        } catch (UnsupportedEncodingException ex) {
            Logger.getLogger(StdoutCapture.class.getName()).log(Level.SEVERE, null, ex);
            return new String[0];
        }
    }

    private String[] removeCRLF(String[] lines) {
        String[] result = new String[lines.length];

        for (int i = 0; i < lines.length; i++) {
            String line = lines[i];
            char lastChar = line.charAt(line.length() - 1);
            while (lastChar == '\n' || lastChar == '\r') {
                line = line.substring(0, line.length() - 1);
                lastChar = line.charAt(line.length() - 1);
            }
            result[i] = line;
        }

        return result;
    }
}
