package main

import "fmt"

var globalCounter = 0

func main() {
	for i := 0; i < 100; i++ {

		func() func() {
			c := globalCounter
			globalCounter++
			return func() {
				fmt.Println(c)
			}
		}()()

	}

}
