package main

import (
	"fmt"
	"strconv"
	"time"

	"go_grpc_example/05_orm/02_xorm/models"
	"go_grpc_example/05_orm/02_xorm/util"

	"go.uber.org/zap"
	"xorm.io/xorm"
)

func main() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	today := time.Now().Format("20060102")
	activeDate, _ := strconv.Atoi(today)
	eg := util.GetEngineGroup()
	eg.Sync2(models.UserActive{}, models.UserActiveRecord{})
	ok, err := eg.Transaction(func(session *xorm.Session) (interface{}, error) {
		var data = models.UserActiveRecord{
			Uid:        19,
			ActiveDate: activeDate,
		}
		_, err := session.Insert(data)
		if err != nil {
			zap.S().Info("判断是否存在错误", err.Error())
			return nil, err
		}
		var actData = models.UserActive{}
		//eg.Insert(actData)
		var uid int64 = 19
		//actData.LatestDate = time.Now().Format("20060102")

		has, err := session.Where("uid=?", uid).Exist(&actData)
		if err != nil {
			zap.S().Info("判断是否存在错误", err.Error())
			return nil, err
		}
		if has {
			ok, err := session.Where("uid=?", uid).Get(&actData)
			if err != nil {
				fmt.Println("出现错误", err.Error())
				return nil, err
			}
			if ok {
				fmt.Printf("获取到数据%+v", actData)
			}
			fmt.Printf("没有获取到%+v", actData)
			actData.LatestDate = int64(activeDate)
			actData.TotalDays += 1
			affected, err := session.Where("uid=?", uid).Cols("total_days,updated").Update(&actData)
			fmt.Println(affected, err)
		} else {
			actData.Uid = uid
			actData.LatestDate = int64(activeDate)
			actData.TotalDays = 1
			affected, err := session.Insert(&actData)
			if err != nil {
				fmt.Println("错误信息是", err)
				return nil, err
			}
			zap.S().Info("用户总活跃表返回Id", actData.Id, "影响的行数", affected)
		}
		return nil, nil
	})
	if err != nil {
		fmt.Println("出现错误", err.Error())
	}
	fmt.Println("事务结果", ok)

}
