// Copyright (C) 2016 Yoshiki Shibata. All rights reserved.

package main

const runner="jp.ne.sonet.ca2.yshibata.JUnitRunner"

const runnerJavaSrc = `
/*
 * Copyright (C) 2016 Yoshiki Shibata. All rights reserved.
 */
package jp.ne.sonet.ca2.yshibata;

import java.io.PrintStream;
import java.text.NumberFormat;
import java.util.ArrayList;
import java.util.List;
import java.util.logging.Level;
import java.util.logging.Logger;
import org.junit.runner.Description;
import org.junit.runner.JUnitCore;
import org.junit.runner.Result;
import org.junit.runner.notification.Failure;
import org.junit.runner.notification.RunListener;

/**
 * JUnitRunner supports -v option
 *
 * @author yoshiki
 */
public class JUnitRunner {

    private static boolean verbose = false;

    public static void main(String[] args) {
        if (args.length == 0) {
            showUsage();
        }

        for (int i = 0; i < args.length; i++) {
            if (args[i].equals("-v")) {
                verbose = true;
                args[i] = null;
            }
        }

        List<Class<?>> classes = new ArrayList<>();
        for (String testClassName : args) {
            if (testClassName == null) {
                continue;
            }

            try {
                classes.add(Class.forName(testClassName));
            } catch (ClassNotFoundException ex) {
                Logger.getLogger(JUnitRunner.class.getName()).log(Level.SEVERE, null, ex);
                System.exit(1);
            }
        }

        JUnitCore core = new JUnitCore();
        core.addListener(new TestListener(System.out));
        Result result = core.run(classes.toArray(new Class<?>[0]));
        if (result.getFailureCount() != 0) {
            System.exit(1);
        }
    }

    private static void showUsage() {
        System.err.println("Usage: JUnitRunner [Test Class Names]");
        System.exit(1);
    }

    /**
     * This TestListener is basically a copy from org.juit.internal.TextListener
     */
    private static class TestListener extends RunListener {

        private final PrintStream writer;

        TestListener(PrintStream writer) {
            this.writer = writer;
        }

        @Override
        public void testRunFinished(Result result) {
            printHeader(result.getRunTime());
            printFailures(result);
            printFooter(result);
        }

        @Override
        public void testStarted(Description description) {
            if (verbose) {
                writer.printf("%s # %s: %n", description.getClassName(), description.getMethodName());
            } else {
                writer.append('.');
            }
        }

        @Override
        public void testFailure(Failure failure) {
            writer.append('E');
        }

        @Override
        public void testIgnored(Description description) {
            writer.append('I');
        }

        protected void printHeader(long runTime) {
            writer.println();
            writer.println("Time: " + elapsedTimeAsString(runTime));
        }

        protected void printFailures(Result result) {
            List<Failure> failures = result.getFailures();
            if (failures.isEmpty()) {
                return;
            }
            if (failures.size() == 1) {
                writer.println("There was " + failures.size() + " failure:");
            } else {
                writer.println("There were " + failures.size() + " failures:");
            }
            int i = 1;
            for (Failure each : failures) {
                printFailure(each, "" + i++);
            }
        }

        private void printFailure(Failure each, String prefix) {
            writer.println(prefix + ") " + each.getTestHeader());
            writer.print(each.getTrace());
        }

        private void printFooter(Result result) {
            if (result.wasSuccessful()) {
                writer.println();
                writer.print("OK");
                writer.println(" (" + result.getRunCount() + " test" + (result.getRunCount() == 1 ? "" : "s") + ")");

            } else {
                writer.println();
                writer.println("FAILURES!!!");
                writer.println("Tests run: " + result.getRunCount() + ",  Failures: " + result.getFailureCount());
            }
            writer.println();
        }

        private String elapsedTimeAsString(long runTime) {
            return NumberFormat.getInstance().format((double) runTime / 1000);
        }
    }
}
`
