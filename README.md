# oak
Oak is a command line tool for Java. Oak is inspired by "go" tool for the Go programming language.

Oak is not a general purpose tool to compile and run any Java projects, but an assitant tool for me to compile solutions of Java exercises for in-house Java training courses.

Currently following commands are supported:

* oak help 
* oak run [Java file] [arguments]
* oak test [-v]
* oak version

**How To Install**
`
go get github.com/YoshikiShibata/oak
`

To use oak, you have to change the currently directory to the directory where Java source files are located and then run oak command.

**help** command shows the help messages.

**run** command compiles the specified Java file and run its main method. 

**test** command compiles all test Java files of which file names end with "Test.java" and run all test methods.
-v option shows both names of test class and test method. To use the test command, you must set JUNIT_HOME environment variable which points to a directory where hamcrest-core-*x.x*.jar and junit-*x.xx*.jar files are found.

**version** command shows the version of the oak command.
