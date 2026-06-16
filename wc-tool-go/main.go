package main

import (
	"fmt"
	"os"

	"github.com/givek/codingchallenges-fyi/wc-tool-go/counter"
)

func main() {
	userArgs := os.Args[1:]

	if len(userArgs) < 2 {
		// TODO: FIXME
		panic("less than 2 args")
	}

	mode := userArgs[0]
	file := userArgs[1]

	f, err := os.Open(file)
	if err != nil {
		panic(err) // TODO: Handle the error gracefully
	}
	defer f.Close()

	switch mode {
	case "-c":
		sizeInBytes, err := counter.SizeInBytes(f)
		if err != nil {
			panic(err) // TODO: Handle the error gracefully
		}
		// fStat, err := f.Stat()
		// if err != nil {
		// 	panic(err)
		// }
		fmt.Printf("%v %v\n", sizeInBytes, f.Name())
	case "-l":
		totalLineCount, err := counter.NumberOfLines(f)
		if err != nil {
			panic(err) // TODO: Handle the error gracefully
		}
		fmt.Printf("%v %v\n", totalLineCount, f.Name())
	case "-w":
		totalWordCount, err := counter.NumberOfWords(f)
		if err != nil {
			panic(err) // TODO: Handle the error gracefully
		}
		fmt.Printf("%v %v\n", totalWordCount, f.Name())
	case "-m":
		totalCharCount, err := counter.NumberOfChars(f)
		if err != nil {
			panic(err) // TODO: Handle the error gracefully
		}

		fmt.Printf("%v %v\n", totalCharCount, f.Name())
	}
}
