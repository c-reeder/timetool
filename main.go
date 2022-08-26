package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func printError(msg string) {
	fmt.Fprintln(os.Stderr, msg)
}

// Default format for timestamp with tz in DBeaver
const dbFormat = "2006-01-02 15:04:05.000 -0700"

func main() {
	var formatStr = flag.String("f", "3339", "The format of timestamp to use")
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		printError("Must provide command")
		os.Exit(-1)
	}

	var format string
	switch *formatStr {
	case "3339":
		format = time.RFC3339
	case "db":
		format = dbFormat
	default:
		format = time.RFC3339
	}

	switch args[0] {
	case "now":
		now(format)
	case "diff":
		diff(format, args[1:])
	default:
		now(format)
	}
}

func now(format string) {
	t1 := time.Now()
	fmt.Println(t1.Format(format))
}

func diff(format string, args []string) {
	if len(args) != 2 {
		printError("Must provide two arguments to diff")
		os.Exit(-1)
	}

	t1, err := time.Parse(format, args[0])
	if err != nil {
		fmt.Println(args[1])
		printError("Error parsing first time stamp")
		printError(err.Error())
		os.Exit(-1)
	}

	t2, err := time.Parse(format, args[1])
	if err != nil {
		fmt.Println(args[2])
		printError("Error parsing second time stamp")
		os.Exit(-1)
	}

	diff := t2.Sub(t1)

	fmt.Println(diff.Milliseconds())
}
