# oak
oak is a command line tool for Java. oak is inspired by "go" tool for the Go programming language.

Currently following commands are supported:

* oak help 
* oak run [Java file] [arguments]
* oak test [-v]
* oak version

help command shows the help messages.

run command compiles the specified Java file and run its main method. 

test command compiles all test Java files of which file names end with "Test.java" and run all test methods.
-v option shows both names of test class and test method. To use the test command, you must set JUNIT_HOME environment variable which points to a directory where hamcrest-core-*xxx*.jar and junit.jar files are found.

version command shows the version of the oak command.
