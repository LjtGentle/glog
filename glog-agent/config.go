package glog_agent

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"os"
)

type AgentConfig struct {
	Targets     []AgentTarget `json:"targets"`
	Address     []string      `json:"address"`
	RequiredAck int           `json:"required_ack"`
}

type AgentTarget struct {
	FileName string `json:"file_name"`
	Topic    string `json:"topic"`
	ChanSize int    `json:"chan_size"`
}

// InitConfig 初始化配置
func InitConfig() *AgentConfig {
	bs, err := os.ReadFile("./config.json")
	if err != nil {
		logrus.Fatalf("config read err=%+v", err)
	}
	config := new(AgentConfig)
	err = json.Unmarshal(bs, config)
	if err != nil {
		logrus.Fatalf("config unmarshal err=%+v", err)
	}
	return config
}
