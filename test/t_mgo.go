package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"log"
)

const (
	URL string = "127.0.0.1:27017"
)

type Person struct {
	Name  string
	Phone string
}

type Men struct {
	Persons []Person
}

func main() {
	session, err := mgo.Dial(URL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	db := session.DB("mydb")     //数据库名称
	collection := db.C("person") //如果该集合已经存在的话，则直接返回

	countNum, err := collection.Count()
	if err != nil {
		panic(err)
	}
	fmt.Println("Things objects count: ", countNum)

	temp := &Person{
		Phone: "18811577546",
		Name:  "zhangzheHero",
	}
	//一次可以插入多个对象 插入两个Person对象
	err = collection.Insert(&Person{"Ale", "+55 53 8116 9639"}, temp)
	if err != nil {
		panic(err)
	}

	result := Person{}
	err = collection.Find(bson.M{"phone": "456"}).One(&result)
	fmt.Println("Phone:", result.Name, result.Phone) //

	var personAll Men //存放结果
	iter := collection.Find(nil).Iter()
	for iter.Next(&result) {
		fmt.Printf("Result: %v\n", result.Name)
		personAll.Persons = append(personAll.Persons, result)
	}
}
