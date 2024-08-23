package main

import (
	"fmt"
)

func errCheckFunc() {
	fmt.Printf("%d", "hello") // want `fmt.Printf format %d has arg "hello" of wrong type string`
}
