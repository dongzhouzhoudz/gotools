package pkg

import (
	"errors"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/dongzhouzhoudz/gotools/config"
	"sync"
)

// 并发线程安全map
var (
	pulsars = new(sync.Map)
)

func GetPulsarByConnName(connName string) (pulsar.Client, error) {

	pulsarConfig := conf.AppConfig.PulsarConfig
	//判断配置中是否包含连接名称
	if _, ok := pulsarConfig[connName]; !ok || len(connName) <= 0 {
		return nil, errors.New(fmt.Sprintf("Not Found Mysql Connection Name: %s,Please Check Your Config Is Right!", connName))
	}

	onePulsarConfig := pulsarConfig[connName]
	onePulsarKey := fmt.Sprintf("%s_%d_%d", onePulsarConfig.BrokerUrl, onePulsarConfig.ConnectionTimeOut, onePulsarConfig.OperationTimeOut)
	onePulsarObject, ok := pulsars.Load(onePulsarKey)
	if ok {
		return onePulsarObject.(pulsar.Client), nil
	}

	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: onePulsarConfig.BrokerUrl,
	})

	if err != nil {
		return nil, err
	} else {
		pulsars.Store(onePulsarKey, client)
		return client, nil
	}

}

func GetNonConfigPulsarClientByBrokerUrl(brokerUrl string) (pulsar.Client, error) {
	onePulsarKey := "inputBrokerUrlPulsar"
	onePulsarObject, ok := pulsars.Load(onePulsarKey)
	if ok {
		return onePulsarObject.(pulsar.Client), nil
	}
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: brokerUrl,
	})
	if err != nil {
		return nil, err
	} else {
		pulsars.Store(onePulsarKey, client)
		return client, nil
	}
}
