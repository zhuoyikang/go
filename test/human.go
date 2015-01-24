package main

import(
	"fmt"
	"strconv"
)

type Person struct {
	name string
	age int
}


type IntSlice []int
type Float32Slice []float32
type PersonSlice []Person

type MaxInterface interface {
	Len() int
	Get(i int) interface{}
	Bigger(i, j int) bool
}

func (x IntSlice) Len() int { return len(x)}
func (x Float32Slice) Len() int { return len(x)}
func (x PersonSlice) Len() int {return len(x)}


func (x IntSlice) Get(i int) interface{} { return x[i]}
func (x Float32Slice) Get(i int) interface{} { return x[i]}
func (x PersonSlice) Get(i int) interface{} {return x[i]}

func (x IntSlice) Bigger(i ,j int ) bool {
	if (x[i] > x[j]) {
		return true;
	}
	return false;
}

func (x Float32Slice) Bigger(i , j int) bool {
	if (x[i] > x[j]) {
		return true;
	}
	return false;
}

func (x PersonSlice) Bigger(i , j int) bool {
	if (x[i].age > x[j].age) {
		return true;
	}
	return false;
}

func (p Person) String() string {
	return "(name:" +p.name + "- age:" + strconv.Itoa(p.age) + " yeas)"
}


func Max(data MaxInterface) (ok bool, max interface{}) {
	if data.Len() == 0 {
		return false, nil
	}
	if data.Len() == 1 {
		return true, data.Get(1)
	}

	max = data.Get(0)
	m := 0
	for i:=1 ;i < data.Len() ;i++ {
		if data.Bigger(i, m) {
			max = data.Get(i)
			m = i
		}
	}

	return true, max
}

//
func test()  {
	fmt.Printf("%s\n", "Testfwefwe")
}


//
func init()  {
	islice := IntSlice{1, 2, 44, 6, 44, 222}
	fslice := Float32Slice{1.99, 3.14, 24.8}
	group := PersonSlice{
		Person{name:"Bart", age:24},
		Person{name:"Bob", age:23},
		Person{name:"Gertrude", age:104},
		Person{name:"Paul", age:44},
		Person{name:"Sam", age:34},
		Person{name:"Jack", age:54},
		Person{name:"Martha", age:74},
		Person{name:"Leo", age:4},
	}

	_, m := Max(islice)
	fmt.Println("The biggest integer in islice is :", m)
	_, m = Max(fslice)
	fmt.Println("The biggest float in fslice is :", m)
	_, m = Max(group)
	fmt.Println("The oldest person in the group is:", m)
}
