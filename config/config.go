package conf

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

var AppConfig = &ConfigMap{}

type ConfigMap struct {
	MysqlConfig  map[string]Mysql  `yaml:"mysql"`
	RedisConfig  map[string]Redis  `yaml:"redis"`
	KafkaConfig  map[string]Kafka  `yaml:"kafka"`
	PulsarConfig map[string]Pulsar `yaml:"pulsar"`
}

// Mysql 实体类型基础配置
type Mysql struct {
	Address              string `yaml:"address"`
	Port                 int    `yaml:"port"`
	UserName             string `yaml:"uname"`
	UserPassword         string `yaml:"upass"`
	DbName               string `yaml:"dbname"`
	MaxConnNumber        int    `yaml:"maxconn"`
	MaxIdleConnNumber    int    `yaml:"maxidleconn"`
	SlowQueryMillSeconds int    `yaml:"slowquerymillseconds"`
}

// Redis 实体类型基础配置
type Redis struct {
	Address  string `yaml:"address"`
	Port     int    `yaml:"port"`
	DbNumber int    `yaml:"db"`
	Password string `yaml:"password"`
}

// Kafka 实体类型基础配置
type Kafka struct {
	Address         string `yaml:"address"`
	ProducerTimeOut int    `yaml:"producertimeout"`
	ConsumerTimeOut int    `yaml:"consumertimeout"`
}

// Pulsar 实体类型基础配置
type Pulsar struct {
	BrokerUrl         string `yaml:"broker"`
	OperationTimeOut  int    `yaml:"operationtimeout"`
	ConnectionTimeOut int    `yaml:"connectiontimeout"`
}

// ConfigInit 初始化基础配置
func ConfigInit(configYamlPath string) {
	dataBytes, readFileError := os.ReadFile(configYamlPath)
	if readFileError != nil {
		fmt.Println(readFileError)
		return
	}
	yamlReadError := yaml.Unmarshal(dataBytes, AppConfig)
	if yamlReadError != nil {
		fmt.Println(yamlReadError)
		return
	}
}
