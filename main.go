package main

import (
	"fmt"
	"os"
	"time"
)

func printError(msg string) {
	fmt.Fprintln(os.Stderr, msg)
}

func main() {
	if len(os.Args) != 3 {
		printError("Must receive 2 arguments!")
		os.Exit(-1)
	}

	t1, err := time.Parse(time.RFC3339, os.Args[1])
	if err != nil {
		fmt.Println(os.Args[1])
		printError("Error parsing first time stamp")
		printError(err.Error())
		os.Exit(-1)
	}

	t2, err := time.Parse(time.RFC3339, os.Args[2])
	if err != nil {
		fmt.Println(os.Args[2])
		printError("Error parsing second time stamp")
		os.Exit(-1)
	}

	diff := t2.Sub(t1)

	fmt.Println(diff.Milliseconds())
}
