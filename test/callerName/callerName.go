package main

import (
	"fmt"
	"os"
	"runtime"
)

func test() {
	callerName()
}

func main() {
	test()

}

func callerName() {
	counter, _, _, success := runtime.Caller(1)

	if !success {
		println("functionName: runtime.Caller: failed")
		os.Exit(1)
	}

	callerName := runtime.FuncForPC(counter).Name()

	fmt.Println(callerName)
}
