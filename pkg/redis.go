package pkg

import (
	"errors"
	"fmt"
	"github.com/dongzhouzhoudz/gotools/config"
	"github.com/go-redis/redis/v7"
	"strconv"
	"sync"
)

// 并发线程安全map
var (
	rs = new(sync.Map)
)

func GetRedisByConnName(connName string) (*redis.Client, error) {
	//获取配置
	configMap := conf.AppConfig
	redisConfigData := configMap.RedisConfig
	//判断配置中是否包含连接名称
	if _, ok := redisConfigData[connName]; !ok || len(connName) <= 0 {
		return nil, errors.New(fmt.Sprintf("Not Found Redis Connection Name: %s,Please Check Your Config Is Right!", connName))
	}

	//获取连接配置
	oneRedisConfigData := redisConfigData[connName]
	redisConnKey := fmt.Sprintf("%s_%d_%s_%d", oneRedisConfigData.Address, oneRedisConfigData.Port, oneRedisConfigData.Password, oneRedisConfigData.DbNumber)

	//判断当前连接对象是否存在，如果存在不需要重复创建
	v, ok := rs.Load(redisConnKey)
	if ok {
		return v.(*redis.Client), nil
	}
	//创建连接
	client := redis.NewClient(&redis.Options{
		Addr:     oneRedisConfigData.Address + ":" + strconv.Itoa(oneRedisConfigData.Port),
		Password: oneRedisConfigData.Password,
		DB:       oneRedisConfigData.DbNumber,
	})
	//连接状态测试
	pong, err := client.Ping().Result()
	fmt.Println(pong)
	if err != nil {
		return nil, err
	} else {
		rs.Store(redisConnKey, client)
		return client, nil
	}

}
