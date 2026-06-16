package counter

import (
	"bytes"
	"io"
	"slices"
	"strings"
)

func SizeInBytes(r io.Reader) (int64, error) {
	buf := make([]byte, 16*1024)
	var sizeInBytes int64 = 0

	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}

		sizeInBytes += int64(n)
	}

	return sizeInBytes, nil
}

func NumberOfLines(r io.Reader) (int64, error) {
	buf := make([]byte, 16*1024)
	var totalLineCount int64 = 0

	for {
		n, err := r.Read(buf)

		// 1. Always process the data read *before* checking the error
		if n > 0 {
			lineCount := bytes.Count(buf[:n], []byte{'\n'})
			totalLineCount += int64(lineCount)
		}

		// 2. Now check if we hit the end of the file or an error
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}
	}

	return totalLineCount, nil
}

func NumberOfWords(r io.Reader) (int64, error) {
	// buf := make([]byte, 16*1024)
	// var totalWordCount int64 = 0
	//
	// lastWord := ""
	//
	// for {
	// 	n, err := r.Read(buf)
	//
	// 	if n > 0 {
	// 		currStr := lastWord + string(buf[:n])
	//
	// 		isLastCharSpace := unicode.IsSpace(rune(currStr[len(currStr)-1]))
	//
	// 		words := strings.Fields(currStr)
	//
	// 		if len(words) > 0 {
	// 			if isLastCharSpace {
	// 				totalWordCount += int64(len(words))
	// 				lastWord = ""
	// 			} else {
	// 				totalWordCount += int64(len(words) - 1)
	// 				lastWord = words[len(words)-1]
	// 			}
	// 		}
	//
	// 	}
	//
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	if err != nil {
	// 		return 0, err
	// 	}
	// }
	//
	// if len(lastWord) != 0 {
	// 	totalWordCount += 1
	// }
	//
	// return totalWordCount, nil

	buf := make([]byte, 16*1024)
	var totalWordCount int64 = 0

	remaining := []byte{}

	for {
		n, err := r.Read(buf)

		if n > 0 {
			currStr := slices.Concat(remaining, buf[:n])

			lastNewLineIdx := bytes.LastIndex(currStr, []byte{'\n'})

			if lastNewLineIdx == -1 {
				remaining = currStr
				continue
			}

			ws := currStr[:lastNewLineIdx+1]

			remaining = currStr[lastNewLineIdx+1:]

			words := strings.Fields(string(ws))
			if len(words) > 0 {
				totalWordCount += int64(len(words))
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}
	}

	if len(remaining) > 0 {
		words := strings.Fields(string(remaining))
		if len(words) > 0 {
			totalWordCount += int64(len(words))
		}
	}

	return totalWordCount, nil
}

func NumberOfChars(r io.Reader) (int64, error) {
	buf := make([]byte, 16*1024)
	var totalCharCount int64 = 0

	remaining := []byte{}

	for {
		n, err := r.Read(buf)

		if n > 0 {
			currStr := slices.Concat(remaining, buf[:n])

			lastNewLineIdx := bytes.LastIndex(currStr, []byte{'\n'})

			if lastNewLineIdx == -1 {
				remaining = currStr
				continue
			}

			ws := currStr[:lastNewLineIdx+1]

			remaining = currStr[lastNewLineIdx+1:]

			for range string(ws) {
				totalCharCount += 1
			}

		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}
	}

	if len(remaining) > 0 {
		for range string(remaining) {
			totalCharCount += 1
		}
	}

	return totalCharCount, nil
}
