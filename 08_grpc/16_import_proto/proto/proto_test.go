package proto

import (
	"fmt"
	"testing"
)

func TestPrint(t *testing.T) {
	c := Computer{
		Name: "mac",
		Cpu: &CPU{
			Name:      "intel",
			Frequency: 4096,
		},
		Memory: &Memory{
			Name: "芝奇",
			Cap:  8192,
		},
	}
	fmt.Println(c.String())
}
