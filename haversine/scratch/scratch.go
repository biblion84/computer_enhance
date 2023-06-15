package main

import (
	"fmt"
	"strconv"
)

func main() {

	b := make([]byte, 0, 20)

	b = append(b, '1')
	b = append(b, '2')

	parsed, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		panic(err)
	}

	fmt.Println(parsed)
	fmt.Println(cap(b))

	b = b[:0]

	fmt.Println(cap(b))
	fmt.Println(string(b))
}
