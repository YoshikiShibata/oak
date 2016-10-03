/*
 * Copyright (C) 2015 Yoshiki Shibata. All rights reserved.
 */
package ch08.ex10;

import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.List;
import java.util.Objects;
import java.util.logging.Level;
import java.util.logging.Logger;
import java.util.stream.Collectors;
import java.util.stream.Stream;

/**
 * Unzip the src.zip file from the JDK. Using Files.walk, find all Java files
 * that contain the keywords transient and volatile.
 *
 * JDK に含まれるsrc.zip ファイルを展開しなさい。Files.walk を使用して、予約語 であるtransient とvolatile
 * を含むJava のファイルをすべて見つけなさい。
 */
public final class JavaScanner {

    private final Path jdkPath;
    private final String[] keyWords;

    /**
     * Constructs an instance.
     *
     * @param jdkPath path of a directory where the source code of JDK is
     * located.
     * @param keyWords key words for searching
     * @throws NullPointerException if either jdkPath is null or keyWords is
     * null.
     * @throws IllegalArgumentException if the specified path doesn't exist or
     * is not a directory.
     * @throws IllegalArgumentException if the lenght of keyWords is zero.
     */
    public JavaScanner(Path jdkPath, String... keyWords) {
        Objects.requireNonNull(jdkPath, "jdkPath is null");
        Objects.requireNonNull(keyWords, "keyWords is null");

        if (keyWords.length == 0) {
            throw new IllegalArgumentException("keyWords is an empty array");
        }

        File file = jdkPath.toFile();

        if (!file.exists()) {
            throw new IllegalArgumentException(jdkPath + " Not Found");
        }

        if (!file.isDirectory()) {
            throw new IllegalArgumentException(jdkPath + " Not Directory");
        }

        this.jdkPath = jdkPath;
        this.keyWords = keyWords;
    }

    /**
     * Scans all source files and return a list of names of source files which
     * contains key words.
     *
     * @return a list of names of source files
     * @throws IOException if the jdkDirecto cannot be opened
     */
    public final List<String> scanFiles() throws IOException {
        try (Stream<Path> entries = Files.walk(jdkPath)) {
            return entries.filter((path) -> !path.toFile().isDirectory()).
                    filter((path) -> containsKeyWords(path)).
                    map((path) -> path.toFile().getAbsolutePath()).
                    collect(Collectors.toList());
        }
    }

    private boolean containsKeyWords(Path path) {
        try {
            try (Stream<String> lines = Files.lines(path)) {
                long count = lines.filter((line) -> {
                    for (String keyWord : keyWords) {
                        if (line.contains(keyWord)) {
                            return true;
                        }
                    }
                    return false;
                }).limit(1).count();
                return count == 1;
            }
        } catch (IOException ex) {
            Logger.getLogger(JavaScanner.class.getName()).log(Level.SEVERE, null, ex);
            return false;
        }
    }

    public static void main(String[] args) throws IOException {
		if (args.length == 0) {
			System.out.println("File Path for the source directory of JDK must be specified");
			System.exit(1);
		}
        JavaScanner js = new JavaScanner(Paths.get(args[0]), "transient", "volatile");

        for (String file : js.scanFiles()) {
            System.out.println(file);
        }
    }

}
