package main

import (
	"fmt"
)

func Test() {
	fmt.Printf("%s\n", "fewfew")
}


func main() {
	var g interface{}
	g=Test
	s := g.(func())
	s()
}
