package main

import (
	"fmt"
	"time"
)

func main() {
	s := NewSlidingWindowLimter()

	for {

		if isOverLimit := s.AddCount(); isOverLimit {
			fmt.Println("over limit !!!")
			time.Sleep(100 * time.Millisecond)
		} else {
			fmt.Println("add success, current count:", s.GetTotalCount())
		}
	}

}
