package glog_server

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"os"
)

type ServerConfig struct {
	KafkaConfig   KafkaConfig   `json:"kafka_config"`
	ElasticConfig ElasticConfig `json:"elastic_config"`
}
type KafkaConfig struct {
	Address []string `json:"address"`
	Targets []Target `json:"targets"`
}

type Target struct {
	ChanSize int    `json:"chan_size"`
	Topic    string `json:"topic"`
}

type ElasticConfig struct {
	Addr  string `json:"addr"`
	Index string `json:"index"`
}

func InitConfig() *ServerConfig {
	fileName := "./config.json"
	b, err := os.ReadFile(fileName)
	if err != nil {
		logrus.Fatalf("read config failed,fileName=%s,err=%+v", fileName, err)
	}
	cfg := new(ServerConfig)
	err = json.Unmarshal(b, cfg)
	if err != nil {
		logrus.Fatalf("json unmarshal failed,err=%+v", err)
	}
	logrus.Infof("init config success,cfg=%+v", cfg)
	return cfg
}
