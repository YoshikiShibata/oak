package jpl.ch01.ex16;

import java.io.FileInputStream;
import java.io.IOException;

public class MyUtilities {
	public double[] getDataSet(String setName) throws BadDataSetException {
		String file = setName + ".dset";
		try(FileInputStream in = new FileInputStream(file)) {
			return readDataSet(in);
		} catch (IOException e) {
			throw new BadDataSetException(file, e);
		}
	}

	public double[] readDataSet(FileInputStream in) throws IOException {
		throw new IOException();
	}
}
