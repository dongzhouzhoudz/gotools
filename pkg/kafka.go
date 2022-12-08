package pkg

import (
	"context"
	"errors"
	"fmt"
	"github.com/dongzhouzhoudz/gotools/config"
	"github.com/segmentio/kafka-go"
	"sync"
	"time"
)

var (
	kafkaConn = new(sync.Map)
)

func GetKafkaProducerConnPerTopicPartition(connName string, topicName string, partitionNumber int) (*kafka.Conn, error) {
	//获取配置
	kafkaConfigMap := conf.AppConfig.KafkaConfig
	//判断配置中是否包含连接名称
	if _, ok := kafkaConfigMap[connName]; !ok || len(connName) <= 0 {
		return nil, errors.New(fmt.Sprintf("Not Found Mysql Connection Name: %s,Please Check Your Config Is Right!", connName))
	}
	//获取连接配置
	kafkaConfig := kafkaConfigMap[connName]
	connPrefix := fmt.Sprintf("%s_Producer_%s_%d", connName, topicName, partitionNumber)
	load, ok := kafkaConn.Load(connPrefix)
	if ok {
		return load.(*kafka.Conn), nil
	}

	leader, err := kafka.DialLeader(context.Background(), "tcp", kafkaConfig.Address, topicName, partitionNumber)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Connect To Kafka Error %s,Please Check Connection Is Right!", connPrefix))
	}
	leader.SetWriteDeadline(time.Now().Add(time.Duration(kafkaConfig.ProducerTimeOut) * time.Second))
	kafkaConn.Store(connPrefix, leader)
	return leader, nil

}

func GetKafkaConsumerConnPerTopicPartition(connName string, topicName string, partitionNumber int) (*kafka.Conn, error) {

	//获取配置
	kafkaConfigMap := conf.AppConfig.KafkaConfig
	//判断配置中是否包含连接名称
	if _, ok := kafkaConfigMap[connName]; !ok || len(connName) <= 0 {
		return nil, errors.New(fmt.Sprintf("Not Found Mysql Connection Name: %s,Please Check Your Config Is Right!", connName))
	}
	//获取连接配置
	kafkaConfig := kafkaConfigMap[connName]
	connPrefix := fmt.Sprintf("%s_Producer_%s_%d", connName, topicName, partitionNumber)
	load, ok := kafkaConn.Load(connPrefix)
	if ok {
		return load.(*kafka.Conn), nil
	}

	leader, err := kafka.DialLeader(context.Background(), "tcp", kafkaConfig.Address, topicName, partitionNumber)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Connect To Kafka Error %s,Please Check Connection Is Right!", connPrefix))
	}
	leader.SetWriteDeadline(time.Now().Add(time.Duration(kafkaConfig.ProducerTimeOut) * time.Second))
	kafkaConn.Store(connPrefix, leader)
	return leader, nil

}
