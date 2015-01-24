package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/rpc"
)

const (
	URL = "127.0.0.1:12981"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int


//
func (t*Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

//
func (t*Arith) Divide(args *Args, quo *Quotient) error {
	if(args.B == 0) {
		return errors.New("divide by zero!")
	}

	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B

	return nil
}

func main() {
	arith := new(Arith)
	rpc.Register(arith)
	rpc.HandleHTTP()

	err := http.ListenAndServe(URL, nil)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
}
