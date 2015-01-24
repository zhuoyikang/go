package main

import (
	"io"
	"log"
	"fmt"
	"net/http"
)

//
func HelloServer(w http.ResponseWriter, req *http.Request)  {
	io.WriteString(w, "hello,world");
}

func main() {
	fmt.Printf("%s\n", "fweew");
	http.HandleFunc("/hello", HelloServer)
	err := http.ListenAndServe("0.0.0.0:6060", nil)
	if(err != nil) {
		log.Fatal("listen failed");
	}
}
