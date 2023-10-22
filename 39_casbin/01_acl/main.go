package main

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v2"
)

func check(e *casbin.Enforcer, sub, obj, act string) {
	ok, _ := e.Enforce(sub, obj, act)
	if ok {
		fmt.Printf("%s CAN %s %s\n", sub, act, obj)
	} else {
		fmt.Printf("%s CANNOT %s %s\n", sub, act, obj)
	}
}

func main() {
	// NewEnforcer creates an enforcer via file or DB.
	e, err := casbin.NewEnforcer("39_casbin/01_acl/model.conf", "39_casbin/01_acl/policy.csv")
	if err != nil {
		log.Fatalf("NewEnforecer failed:%v\n", err)
	}

	check(e, "danny", "data1", "read")
	check(e, "joy", "data2", "write")
	check(e, "danny", "data1", "write")
	check(e, "danny", "data2", "read")
}
