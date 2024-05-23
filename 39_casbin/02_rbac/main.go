package main

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v2"
)

func check(e *casbin.Enforcer, sub, obj, act string) {
	ok, reason, _ := e.EnforceEx(sub, obj, act)
	if ok {
		fmt.Printf("%s CAN %s %s ,reason:%s \n", sub, act, obj, reason)
	} else {
		fmt.Printf("%s CANNOT %s %s\n", sub, act, obj)
	}
}

func main() {
	// NewEnforcer creates an enforcer via file or DB.
	e, err := casbin.NewEnforcer("39_casbin/02_rbac/model.conf", "39_casbin/02_rbac/policy.csv")
	if err != nil {
		log.Fatalf("NewEnforecer failed:%v\n", err)
	}

	check(e, "danny", "data", "read")
	check(e, "danny", "data", "write")
	check(e, "joy", "data", "read")
	check(e, "joy", "data", "write")
}
