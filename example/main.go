package main

import (
	"fmt"
	"os"

	"xyzc.dev/go/profile/gc"
)

func main() {
	timer := gc.NewSectionTimer("main")
	doAllocates()
	timer.Mark("middle")
	doAllocates()
	timer.End("end")
	fmt.Printf("timer: %v\n", timer)
	fmt.Println(timer.CSVHeader())
	fmt.Println(timer.CSVString())
}

func doAllocates() {
	for i := 0; i < 3000000; i++ {
		fmt.Fprintf(os.Stderr, "this is a string: %v", i)
	}
}
