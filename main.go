package main

import (
	"fmt"
	"syscall"
)

func main() {
	// convertToOpus("test/audio/16k-54s.wav")
	testlib, err := syscall.LoadLibrary("lib/linux/x64/libopus.so")
	if err != nil {
		panic(err)

	}
	fmt.Println(testlib)
}
