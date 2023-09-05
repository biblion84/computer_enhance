package main

import (
	"fmt"
	"strconv"
)

func main() {

	for i := 150_000_000; i < 2000000000; i++ {
		fmt.Println(prettyPrint(i))
	}

}

func prettyPrint(x int) string {
	printed := []rune(strconv.Itoa(x))

	prettyPrinted := ""

	for i := 0; i < len(printed); i++ {
		if i%3 == 0 && i != 0 {
			prettyPrinted = "_" + prettyPrinted
		}

		prettyPrinted = string(printed[len(printed)-1-i]) + prettyPrinted
	}

	return prettyPrinted
}
