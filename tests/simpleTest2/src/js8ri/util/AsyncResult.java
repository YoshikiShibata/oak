/*
 * Copyright (C) 2014 Yoshiki Shibata. All rights reserved.
 */
package js8ri.util;

import java.util.logging.Level;
import java.util.logging.Logger;

/**
 * Utility class to wait for the result of an Async operation.
 */
public class AsyncResult {

    private Boolean result = null;

    public synchronized boolean waitForResult() {
        while (result == null) {
            try {
                wait();
            } catch (InterruptedException ex) {
                Logger.getLogger(AsyncResult.class.getName()).log(Level.SEVERE, null, ex);
            }
        }
        return result;
    }

    public synchronized void setResult(boolean result) {
        this.result = result;
        notifyAll();
    }
}
