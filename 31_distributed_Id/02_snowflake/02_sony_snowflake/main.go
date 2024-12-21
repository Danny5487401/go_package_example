package main

import (
	"fmt"
	"time"

	"github.com/sony/sonyflake"
)

func getMachineID() (uint16, error) {
	return uint16(1), nil
}

func checkMachineID(machineID uint16) bool {
	return true
}

func main() {
	t := time.Unix(0, 0)
	settings := sonyflake.Settings{
		StartTime:      t,
		MachineID:      getMachineID,
		CheckMachineID: checkMachineID,
	}
	sf := sonyflake.NewSonyflake(settings)
	id, _ := sf.NextID()
	fmt.Println(id)
}
