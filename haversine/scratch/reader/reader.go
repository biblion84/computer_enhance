package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"
)

func p(err error) {
	if err != nil {
		panic(err)
	}
}

func printTiming(t *time.Time, fileSize int64, testCase string) {
	timeElapsed := time.Since(*t)
	mbPerSecond := int(((float64(fileSize) / 1_000_000.0) / float64(timeElapsed.Milliseconds())) * 1000)
	fmt.Printf("%s: took %d ms, %d MB/s\n", testCase, timeElapsed.Milliseconds(), mbPerSecond)
	*t = time.Now()
}

func main() {
	filepath := path.Join("data", "haversine.json")
	fileinfo, err := os.Stat(filepath)
	p(err)

	file, err := os.Open(filepath)
	p(err)

	start := time.Now()

	content, err := io.ReadAll(file)
	p(err)
	println(len(content))

	printTiming(&start, fileinfo.Size(), "warming up")
	file.Seek(0, 0)

	_, err = io.ReadAll(file)
	p(err)
	println(len(content))
	printTiming(&start, fileinfo.Size(), "basic")
	file.Seek(0, 0)

	content, err = io.ReadAll(file)
	p(err)
	println(len(content))
	printTiming(&start, fileinfo.Size(), "basic2")
	file.Seek(0, 0)

}
