package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// Print error to stderr
func printError(msg string) {
	fmt.Fprintln(os.Stderr, msg)
}

type Format string

// Format to/from functions
type fromFunc func(string) (time.Time, error)
type toFunc func(time.Time) string

const (
	rfc3339 Format = "rfc3339" // RFC3339
	db             = "db"      // Default format for timestamp with tz in DBeaver
	pb             = "pb"      // Protobuf components (useful in Bloom)
	ms             = "ms"      // Unix milliseconds
)

// Default format for timestamp with tz in DBeaver
const dbFormat = "2006-01-02 15:04:05.000 -0700"

var fromFormats = map[Format]fromFunc{
	rfc3339: func(input string) (output time.Time, err error) {
		return time.Parse(time.RFC3339, input)
	},
	db: func(input string) (output time.Time, err error) {
		return time.Parse(dbFormat, input)
	},
	pb: func(input string) (output time.Time, err error) {
		return time.Parse(dbFormat, input)
	},
	ms: func(input string) (output time.Time, err error) {
		milli, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			return time.Time{}, err
		}
		return time.UnixMilli(milli), nil
	},
}

var toFormats = map[Format]toFunc{
	rfc3339: func(input time.Time) (output string) {
		return input.Format(time.RFC3339)
	},
	db: func(input time.Time) (output string) {
		return input.Format(dbFormat)
	},
	ms: func(input time.Time) (output string) {
		return strconv.FormatInt(input.UnixMilli(), 10)
	},
	pb: func(input time.Time) (output string) {
		pbTimestamp := timestamppb.New(input)
		return fmt.Sprintf("UTC Timestamp %v\nLocal Timestamp %v\nSeconds %v\nNanos %v\n", pbTimestamp.AsTime().UTC(), pbTimestamp.AsTime().Local(), pbTimestamp.Seconds, pbTimestamp.Nanos)
	},
}

var commands = map[string]func(fromFunc, toFunc, []string){
	"now":  now,
	"diff": diff,
	"conv": convert,
}

func main() {
	var inputFormatStr = flag.String("i", "rfc3339", "The format of timestamp to use")
	var outputFormatStr = flag.String("o", "rfc3339", "The format of timestamp to use")
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		printError("Must provide command")
		os.Exit(-1)
	}

	inputFormat := Format(*inputFormatStr)
	outputFormat := Format(*outputFormatStr)

	inputFunc := fromFormats[inputFormat]
	outputFunc := toFormats[outputFormat]

	if inputFunc == nil {
		printError("Invalid input format!")
		os.Exit(-1)
	}
	if outputFunc == nil {
		printError("Invalid output format!")
		os.Exit(-1)
	}

	command := args[0]

	commands[command](inputFunc, outputFunc, args[1:])
}

// Print current timestamp
func now(inputFunc fromFunc, outputFunc toFunc, args []string) {
	t1 := time.Now()
	fmt.Println(outputFunc(t1))
}

// Print difference between two provided timestamps in milliseconds
func diff(inputFunc fromFunc, outputFunc toFunc, args []string) {
	if len(args) != 2 {
		printError("Must provide two arguments to diff")
		os.Exit(-1)
	}

	t1, err := inputFunc(args[0])
	if err != nil {
		fmt.Println(args[1])
		printError("Error parsing first time stamp")
		printError(err.Error())
		os.Exit(-1)
	}

	t2, err := inputFunc(args[1])
	if err != nil {
		fmt.Println(args[2])
		printError("Error parsing second time stamp")
		os.Exit(-1)
	}

	diff := t2.Sub(t1)

	fmt.Println(diff.Milliseconds())
}

// Convert from one format to another
func convert(inputFunc fromFunc, outputFunc toFunc, args []string) {
	if len(args) != 1 {
		printError("Must provide two arguments to diff")
		os.Exit(-1)
	}

	t1, err := inputFunc(args[0])
	if err != nil {
		fmt.Println(args[0])
		printError("Error parsing time stamp")
		printError(err.Error())
		os.Exit(-1)
	}
	fmt.Println(outputFunc(t1))
}
