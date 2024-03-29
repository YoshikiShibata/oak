# oak
Oak is a command line tool for Java. Oak is inspired by "go" tool for the Go programming language.

Oak is not a general purpose tool to compile and run any Java projects, but a utility tool for me to compile solutions of Java exercises for in-house Java training courses.

Currently following commands are supported:

* oak help 
* oak run [Java source file] [arguments]
* oak run [jar file]
* oak test [-v]
* oak version

##  **How To Install:**

To install oak, you need [**go** command](https://golang.org/), which must be Go 1.18+.

1. You have to install oak with the following command:
```
go install github.com/YoshikiShibata/oak@latest 
```

2. **oak** command will be build as `$GOPATH/bin/oak`.

3. After installing, then set `OAK_HOME` environment variable to the installed directory which is typically `$GOPATH/src/github.com/YoshikiShibata/oak`.

4. For JavaFX application, you need to install JavaFX SDK from [OpenJFX](https://openjfx.io/). After installing JavaFX SDK, set `PATH_TO_FX` environment variable to its `lib` directory. Now oak automatically use `javafx.controls`, `javafx.fxml`, `javafx.web`, and `javafx.swing` modules if the environment variable is set.

## **Commands**

To use oak, you have to change the currently directory to the directory where Java source files are located and then run oak command.

* **help** command shows the help messages.

* **run** command compiles the specified Java file and run its main method. If java file is not specified such as `oak run`, then all `.java` files are searched locally and the first one which has a line starting with either `public static void main` or `static public void main` will be considered as the java file. If no such java file is not found, then jar files are searched locally and the first one will be executed with `-jar`.

* **test** command compiles all JUnit-based test Java files and run all test methods.
`-v` option shows both names of test class and test method. `-run=` option accept a regular expression for filtering test methods. For example, `-run=.` will run all tests.
If tests will not be completed within one minute, then they will be aborted.

* **version** command shows the version of the oak command.

## **Notes**

* If source files and test files are located separately, then they must be located under `src` and `test` directories respectively.  
* The encoding of Java source files must be UTF-8.


Copyright (C) 2016, 2017, 2020 - 2022 Yoshiki Shibata. All rights reserved.
