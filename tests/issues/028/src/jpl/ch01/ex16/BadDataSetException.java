package jpl.ch01.ex16;

import java.io.IOException;

public class BadDataSetException extends Exception {
	private IOException ioException;
	private String file;

	BadDataSetException(String file, IOException ioException) {
		this.file = file;
		this.ioException = ioException;
	}

	@Override
	public String toString() {
		return "file=" + file + ", IOException=" + ioException;
	}
}
