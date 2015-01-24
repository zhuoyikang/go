package main

import (
	"fmt"
)

type Engine struct {
}

//T
func (Engine) work1() {
	fmt.Println("[1]: I am an engine")
}

//*T
func (*Engine) work2() {
	fmt.Println("[2]: I am an engine")
}

//T方式组合
type FoxCar struct {
	Engine
}

//*T方式组合
type SmartCar struct {
	*Engine
}

func main() {
	//fCar := new(FoxCar)
	//sCar := new(SmartCar)
	fCar := FoxCar{}
	sCar := SmartCar{}
	//以T方式组合其他类型的时候可以访问T类型中定义为T和*T的方法
	fCar.work1()
	fCar.work2()

	fmt.Println("--------------------")

	//以*T方式组合其他类型的时候仅可以访问T类型中定义为*T的方法
	sCar.work2()
}
