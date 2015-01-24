package main

import (
	"runtime"
	"fmt"
	"log"
	"github.com/user/stringutil"
)

const (
	TEST_SIZE =  1000
)

//
func test_logger()  {
	pipe := make(chan int, TEST_SIZE)
	for i:=0; i< TEST_SIZE ;i++ {
		go func() {
			log.Println("log_bufsize send",i);
			pipe <- 1
		}()
	}

	for i := 1; i< TEST_SIZE; i++  {
		<- pipe
		log.Println("log_bufsize recv",i);
	}
	fmt.Println("test over");
}

func test_reverse () {
	log.Println("log_bufsize:",23)
	stringutil.Good("!oG ,olleH")
	fmt.Printf(stringutil.Reverse("!oG ,olleH"))
}

//
func test_struct()  {
	//_s := stringutil.Node{"nice"}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	test_logger()
}
