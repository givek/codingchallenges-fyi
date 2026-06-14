package main

import (
	"fmt"
	"io"
	"os"
)

func SizeInBytes(r io.Reader) (int64, error) {
	buf := make([]byte, 128)
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

func main() {
	f, err := os.Open("./test.txt")
	if err != nil {
		panic(err) // TODO: Handle the error gracefully
	}
	defer f.Close()

	sizeInBytes, err := SizeInBytes(f)
	if err != nil {
		panic(err) // TODO: Handle the error gracefully
	}

	// fStat, err := f.Stat()
	// if err != nil {
	// 	panic(err)
	// }

	fmt.Printf("%v %v\n", sizeInBytes, f.Name())
}
