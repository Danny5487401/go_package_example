package _2_clockwork

import (
	"fmt"
	"github.com/jonboulle/clockwork"
	"time"
)

func myFunc(clock clockwork.Clock) {
	clock.Sleep(3 * time.Second)
	doSomething()
}

func doSomething() {
	fmt.Println("do something")
}
