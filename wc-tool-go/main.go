package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/givek/codingchallenges-fyi/wc-tool-go/counter"
)

func main() {
	userArgs := os.Args[1:]

	mode := "-default"

	var rs io.ReadSeeker

	fileName := ""

	if len(userArgs) == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		rs = bytes.NewReader(data)
	} else if len(userArgs) == 1 {
		arg := userArgs[0]
		if strings.HasPrefix(arg, "-") {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				panic(err)
			}

			rs = bytes.NewReader(data)
			mode = arg
		} else {
			f, err := os.Open(arg)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			rs = f
			fileName = f.Name()
		}
	} else {
		mode = userArgs[0]
		f, err := os.Open(userArgs[1])
		if err != nil {
			panic(err)
		}
		defer f.Close()
		rs = f
		fileName = f.Name()
	}

	switch mode {
	case "-c":
		sizeInBytes, err := counter.SizeInBytes(rs)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%v %v\n", sizeInBytes, fileName)
	case "-l":
		totalLineCount, err := counter.NumberOfLines(rs)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%v %v\n", totalLineCount, fileName)
	case "-w":
		totalWordCount, err := counter.NumberOfWords(rs)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%v %v\n", totalWordCount, fileName)
	case "-m":
		totalCharCount, err := counter.NumberOfChars(rs)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%v %v\n", totalCharCount, fileName)

	case "-default":
		sizeInBytes, err := counter.SizeInBytes(rs)
		if err != nil {
			panic(err)
		}
		rs.Seek(0, io.SeekStart)
		totalLineCount, err := counter.NumberOfLines(rs)
		if err != nil {
			panic(err)
		}
		rs.Seek(0, io.SeekStart)
		totalWordCount, err := counter.NumberOfWords(rs)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%v %v %v %v\n", totalLineCount, totalWordCount, sizeInBytes, fileName)
	}
}
