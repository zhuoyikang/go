package main

import "encoding/json"
import "fmt"

type Message struct {
	Name string
	Body string
	Time int64
}

func main() {
	m := Message{"Alice", "Hello", 1294706395881547000}
	b, _ := json.Marshal(m)
	fmt.Printf("%s\n", b)

	var m1 Message
	json.Unmarshal(b, &m1)
	fmt.Printf("%v\n", m1)

	test()
}

//
func test()  {
	b := []byte(`{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}`)
	var f interface{}
	json.Unmarshal(b, &f)
	fmt.Printf("%v\n", f)
}
