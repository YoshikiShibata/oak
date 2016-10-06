/*
 * Copyright (C) 2014 Yoshiki Shibata. All rights reserved.
 */
import java.io.File;
import js8ri.util.Directories;
import static org.junit.Assert.*;
import org.junit.Test;

/**
 * Test code for ListAllSubdirectories. This code must be run on a Unix such
 * as Linux or Mac OS X.
 */
public class ListAllSubdirectories2Test {

    @Test
    public void listAllSubDirectories() {
        // Prepare
        File dir = Directories.toDirectory("/usr/include");
        assertNotNull(dir);

        // Action
        File[] subdirectories = ListAllSubdirectories.listAllSubDirectories(dir);

        // Check
        for (File sub : subdirectories) {
            if (!sub.isDirectory()) {
                fail(sub.getName() + " is not a directory");
            }
        }
    }
}
