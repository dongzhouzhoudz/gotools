package pkg

import (
	"errors"
	"fmt"
	"github.com/dongzhouzhoudz/gotools/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"sync"
	"time"
)

//并发线程安全map
var (
	dbs = new(sync.Map)
)

func GetMysqlConnectByName(connName string) (*gorm.DB, error) {
	//获取配置
	mysqlConfigMap := conf.AppConfig.MysqlConfig
	//判断配置中是否包含连接名称
	if _, ok := mysqlConfigMap[connName]; !ok || len(connName) <= 0 {
		return nil, errors.New(fmt.Sprintf("Not Found Mysql Connection Name: %s,Please Check Your Config Is Right!", connName))
	}
	//获取连接配置
	oneMysqlConfig := mysqlConfigMap[connName]
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4,utf8&parseTime=True&loc=Local", oneMysqlConfig.UserName, oneMysqlConfig.UserPassword, oneMysqlConfig.Address, oneMysqlConfig.Port, oneMysqlConfig.DbName)
	//判断当前连接对象是否存在，如果存在不需要重复创建
	v, ok := dbs.Load(dsn)
	if ok {
		return v.(*gorm.DB), nil
	}

	slowLogQueryMillSeconds := 1000
	if oneMysqlConfig.SlowQueryMillSeconds > 0 {
		slowLogQueryMillSeconds = oneMysqlConfig.SlowQueryMillSeconds

	}

	logger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Duration(slowLogQueryMillSeconds) * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	//创建连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger})

	if err != nil {
		return nil, err
	}
	//创建实体链接对象
	dbObject, dbErr := db.DB()

	if dbErr != nil {
		return nil, dbErr
	}

	if err := dbObject.Ping(); err != nil {
		return nil, err
	}

	//设置数据库最大连接数量
	dbObject.SetMaxOpenConns(oneMysqlConfig.MaxConnNumber)
	//设置数据库空闲最大连接数量
	dbObject.SetMaxIdleConns(oneMysqlConfig.MaxIdleConnNumber)
	//存储连接对象
	dbs.Store(dsn, db)
	return db, err

}
