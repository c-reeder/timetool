package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func printError(msg string) {
	fmt.Fprintln(os.Stderr, msg)
}

// Default format for timestamp with tz in DBeaver
const dbFormat = "2006-01-02 15:04:05.000 -0700"

func main() {
	var inputFormatStr = flag.String("f", "rfc3339", "The format of timestamp to use")
	var outputFormatStr = flag.String("o", "3339", "The format of timestamp to use")
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		printError("Must provide command")
		os.Exit(-1)
	}

	inputFormat := expandFormatName(*inputFormatStr)
	outputFormat := expandFormatName(*outputFormatStr)

	switch args[0] {
	case "now":
		now(inputFormat)
	case "diff":
		diff(inputFormat, args[1:])
	case "conv":
		convert(inputFormat, outputFormat, args[1:])
	case "pb":
		pbComponents(inputFormat, args[1:])
	default:
		now(inputFormat)
	}
}

func expandFormatName(abr string) string {
	switch abr {
	case "db":
		return dbFormat
	default:
		return time.RFC3339
	}
}

// Print current timestamp
func now(format string) {
	t1 := time.Now()
	fmt.Println(t1.Format(format))
}

// Print difference between two provided timestamps in milliseconds
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

// Convert from one format to another
func convert(inputFormat, outputFormat string, args []string) {
	if len(args) != 1 {
		printError("Must provide two arguments to diff")
		os.Exit(-1)
	}

	t1, err := time.Parse(inputFormat, args[0])
	if err != nil {
		fmt.Println(args[0])
		printError("Error parsing time stamp")
		printError(err.Error())
		os.Exit(-1)
	}
	fmt.Println(t1.Format(outputFormat))
}

// Break timestamp into pb components
func pbComponents(inputFormat string, args []string) {
	t1, err := time.Parse(inputFormat, args[0])
	if err != nil {
		fmt.Println(args[0])
		printError("Error parsing time stamp")
		printError(err.Error())
		os.Exit(-1)
	}

	pbTimestamp := timestamppb.New(t1)
	fmt.Printf("UTC Timestamp %v\n", pbTimestamp.AsTime().UTC())
	fmt.Printf("Local Timestamp %v\n", pbTimestamp.AsTime().Local())
	fmt.Printf("Seconds %v\n", pbTimestamp.Seconds)
	fmt.Printf("Nanos %v\n", pbTimestamp.Nanos)
}
