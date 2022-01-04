// Copyright © 2016 Yoshiki Shibata. All rights reserved.

package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Basically most of code here are copied from the source code for the go command and modified for oak.

var helpTemplate = `{{if .Runnable}}usage: oak {{.UsageLine}}

{{end}}{{.Long | trim}}
`

var usageTemplate = `Oak is a tool for compiling and running Java source code.

Usage:

    oak command [arguments]

The commands are:
{{range .}}{{if .Runnable}}
{{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}

Use "go help [command]" for more information about a command.
`

func help(args []string) {
	if len(args) == 0 {
		printUsage(os.Stdout)
		return
	}

	arg := args[0]

	for _, cmd := range commands {
		if cmd.Name() == arg {
			tmpl(os.Stdout, helpTemplate, cmd)
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknow help topic %#q. Run 'oak help'.\n", arg)
	os.Exit(2)
}

func printUsage(w io.Writer) {
	bw := bufio.NewWriter(w)
	tmpl(bw, usageTemplate, commands)
	bw.Flush()
}

// An errWriter wraps a writer, recording whether a write error occurred.
type errWriter struct {
	w   io.Writer
	err error
}

func (w *errWriter) Write(b []byte) (int, error) {
	n, err := w.w.Write(b)
	if err != nil {
		w.err = err
	}
	return n, err
}

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data any) {
	t := template.New("top")
	t.Funcs(template.FuncMap{"trim": strings.TrimSpace, "capitalize": capitalize})
	template.Must(t.Parse(text))
	ew := &errWriter{w: w}
	err := t.Execute(ew, data)
	if ew.err != nil {
		// I/O error writing. Ignore write on closed pipe.
		if strings.Contains(ew.err.Error(), "pipe") {
			os.Exit(1)
		}
		exit(fmt.Errorf("writing output: %v", ew.err), 1)
	}
	if err != nil {
		panic(err)
	}
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToTitle(r)) + s[n:]
}
