package main

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"log"
)

type Object struct {
	Name  string
	Owner string
}

type Subject struct {
	Name string
	Hour int
}

func check(e *casbin.Enforcer, sub Subject, domain string, obj Object, act string) {
	ok, _ := e.Enforce(sub, domain, obj, act)
	if ok {
		fmt.Printf("%s CAN %s %s in %s at %d:00\n", sub.Name, act, obj.Name, domain, sub.Hour)
	} else {
		fmt.Printf("%s CANNOT %s %s in %s at %d:00\n", sub.Name, act, obj.Name, domain, sub.Hour)
	}
}

func main() {
	e, err := casbin.NewEnforcer("39_casbin/04_abac/model.conf")
	if err != nil {
		log.Fatalf("NewEnforecer failed:%v\n", err)
	}
	_, err = e.AddPolicies([][]string{
		{"(r.sub.Hour>= 9 && r.sub.Hour < 18)|| r.sub.Name == r.obj.Owner", "library", "data", "read"}, // 正常工作时间9:00-18:00所有人都可以在library 读 data
	})
	if err != nil {
		log.Fatalf("add policy failed:%v\n", err)
	}
	o := Object{"data", "danny"} // 创建被访问资源对象

	domain := "library"
	wrongDomain := "coffee shop"

	s1 := Subject{"danny", 2}
	check(e, s1, domain, o, "read") // danny CAN read data in library at 2:00

	s2 := Subject{"joy", 10}
	check(e, s2, domain, o, "read") // joy CAN read data in library at 10:00

	s3 := Subject{"danny", 20}
	check(e, s3, domain, o, "read") // danny CAN read data in library at 20:00

	s4 := Subject{"danny", 20}
	check(e, s4, wrongDomain, o, "read") // danny CANNOT read data in coffee shop at 20:00

	s5 := Subject{"joy", 20}
	check(e, s5, domain, o, "read") // joy CANNOT read data in library at 20:00
}
