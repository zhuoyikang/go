package main

import (
	"github.com/peterbourgon/g2s"
)
//
func main( )  {
	s, err := g2s.Dial("udp", "54.200.145.61:8125")
	if err != nil {
		return;
	}
	s.Counter(1.0, "test.g2s",1)
}
