package main

import "fmt"

type Ints2 []int;

//
func (i Ints2) Iterator() func() (int, bool) {
	index := 0
	return func() (val int, ok bool) {
		if index >= len(i) {
			return
		}
		val, ok = i[index], true
		index++
		return
	}
}


func main() {
	ints := Ints2{1,2,3}
	it := ints.Iterator()
	for  {
		val, ok := it()
		if !ok {
			break
		}
		fmt.Printf("%d\n", val)
	}
}
