package main

import (
	"fmt"
	"time"
)

// now work out how to pass this around as objects
func main() {
	messages := NewLogger()
	messages.Run()

	go func() {
		i := 1
		for {
			messages.Log(fmt.Sprintf("piing%v", i))
			time.Sleep(1 * time.Second)
			i++
		}
	}()
	for {
		fmt.Println("j-chillin")
		time.Sleep(4 * time.Second)
	}
}
