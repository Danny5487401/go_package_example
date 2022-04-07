package main

import (
	"strconv"
	"time"

	"go_package_example/05_rds/02_xorm/models"
	"go_package_example/05_rds/02_xorm/util"

	"go.uber.org/zap"
	"xorm.io/xorm"
)

func main() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	today := time.Now().Format("20060102")
	activeDate, _ := strconv.Atoi(today)
	eg := util.GetEngineGroup()
	//eg.Sync2(models.UserActive{}, models.UserActiveRecord{})
	ok, err := eg.Transaction(func(session *xorm.Session) (interface{}, error) {
		var data = models.UserActiveRecord{
			Uid:         19,
			ActiveDate:  activeDate,
			ExtendField: "1",
		}
		_, err := session.Insert(data)
		if err != nil {
			zap.S().Info("插入数据存在错误", err.Error())
			// 错误会自动回滚 rollback()
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
				zap.S().Info("出现错误", err.Error())
				return nil, err
			}
			if ok {
				zap.S().Infof("获取到数据%+v", actData)
			}
			zap.S().Infof("没有获取到%+v", actData)
			actData.LatestDate = int64(activeDate)
			actData.TotalDays += 1
			affected, err := session.Where("uid=?", uid).Cols("total_days,updated").Update(&actData)
			zap.S().Info(affected, err)
		} else {
			actData.Uid = uid
			actData.LatestDate = int64(activeDate)
			actData.TotalDays = 1
			affected, err := session.Insert(&actData)
			if err != nil {
				zap.S().Info("错误信息是", err)
				return nil, err
			}
			zap.S().Info("用户总活跃表返回Id", actData.Id, "影响的行数", affected)
		}
		return nil, nil
	})
	if err != nil {
		zap.S().Info("出现错误", err.Error())
	}
	zap.S().Info("事务结果", ok)

}

/*
源码： session配置
func newSession(engine *Engine) *Session {
	var ctx context.Context
	if engine.logSessionID {
		ctx = context.WithValue(engine.defaultContext, log.SessionIDKey, newSessionID())
	} else {
		ctx = engine.defaultContext
	}

	session := &Session{
		ctx:    ctx,
		engine: engine,
		tx:     nil,
		statement: statements.NewStatement(
			engine.dialect,
			engine.tagParser,
			engine.DatabaseTZ,
		),
		isClosed:               false,
		isAutoCommit:           true,
		isCommitedOrRollbacked: false,
		isAutoClose:            false,
		autoResetStatement:     true,
		prepareStmt:            false,

		afterInsertBeans: make(map[interface{}]*[]func(interface{}), 0),
		afterUpdateBeans: make(map[interface{}]*[]func(interface{}), 0),
		afterDeleteBeans: make(map[interface{}]*[]func(interface{}), 0),
		beforeClosures:   make([]func(interface{}), 0),
		afterClosures:    make([]func(interface{}), 0),
		afterProcessors:  make([]executedProcessor, 0),
		stmtCache:        make(map[uint32]*core.Stmt),

		lastSQL:     "",
		lastSQLArgs: make([]interface{}, 0),

		sessionType: engineSession,
	}
	if engine.logSessionID {
		session.ctx = context.WithValue(session.ctx, log.SessionKey, session)
	}
	return session
}
*/
