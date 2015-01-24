package main

import (
	"fmt"
	"log"
)

//
func init()  {
	fmt.Println("test1");
}

func init()  {
	fmt.Println("test2");
}


//
func main()  {
	log.Println("error opening file %v\n", 2323);
}
