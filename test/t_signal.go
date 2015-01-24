package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)


func main() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGTERM)

	for {
		msg := <-ch
		switch msg {
		case syscall.SIGHUP:
			fmt.Printf("hub\n");
		case syscall.SIGTERM:
			fmt.Printf("term\n");
		}
	}
}
