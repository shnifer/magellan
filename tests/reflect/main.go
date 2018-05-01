package main

import (
	"log"
	"reflect"
)

type myStruct struct{
	X,Y float64
}

func (a myStruct) Mul (b myStruct) (res myStruct) {
	va:=reflect.ValueOf(a)
	vb:=reflect.ValueOf(b)
	vr:=reflect.ValueOf(&res).Elem()
	t:=reflect.TypeOf(a)
	fc:=t.NumField()
	for i:=0;i<fc;i++{
		x:=va.Field(i).Float()*vb.Field(i).Float()
		vr.Field(i).SetFloat(x)
	}
	return res
}

func main(){
	a:=myStruct{0.5,2.0}
	b:=myStruct{1.5,3.0}
	c:=a.Mul(b)
	log.Println(a)
	log.Println(b)
	log.Println(c)
}
