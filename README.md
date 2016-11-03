# oak
Oak is a command line tool for Java. Oak is inspired by "go" tool for the Go programming language.

Oak is not a general purpose tool to compile and run any Java projects, but an assitant tool for me to compile solutions of Java exercises for in-house Java training courses.

Currently following commands are supported:

* oak help 
* oak run [Java file] [arguments]
* oak test [-v]
* oak version

##  **How To Install:**

To install oak, you need [**go** command](https://golang.org/). 

1. You have to install oak with the following command:
`go get github.com/YoshikiShibata/oak`. 

2. If you want to update the oak to the latest version:
`go get -u github.com/YoshikiShibata/oak`.

3. **oak** command will be build as `$GOPATH/bin/oak`.

4. After installing, then set `OAK_HOME` environment variable to the installed directory which is typically `$GOPATH/src/github.com/YoshikiShibata/oak`.

## **Commands**

To use oak, you have to change the currently directory to the directory where Java source files are located and then run oak command.

* **help** command shows the help messages.

* **run** command compiles the specified Java file and run its main method. If java file is not specified such as `oak run`, then all `.java` files are searched locally and the first one which has a line starting with either `public static void main(String` or `static public void main(String` will be considered as the java file.

* **test** command compiles all test Java files of which file names end with "Test.java" and run all test methods.
-v option shows both names of test class and test method. 
If tests will not be completed within one minute, then they will be aborted.

* **version** command shows the version of the oak command.

## **Notes**

* If source files and test files are located separately, then they must be located under `src` and `test` directories respectively.  
* The encoding of Java source files must be UTF-8.
