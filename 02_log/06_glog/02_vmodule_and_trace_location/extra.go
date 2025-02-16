package main

import "github.com/golang/glog"

func bar() {
	glog.V(4).Info("LEVEL 4: level 4 message in bar.go")
}
