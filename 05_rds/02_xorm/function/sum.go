package main

import (
	"go_package_example/05_rds/02_xorm/models"
	"go_package_example/05_rds/02_xorm/util"

	"fmt"
)

func main() {
	var err error
	eg := util.GetEngineGroup()

	ms := new(models.MasterSlaveTable)
	total, err := eg.Where("id >?", 2).Sum(ms, "id")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("总和是total:", total)

}
