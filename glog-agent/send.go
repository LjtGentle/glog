package glog_agent

import (
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

var kafkaClient sarama.SyncProducer

// InitKafka 初始化kafka
func InitKafka(config *AgentConfig) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.RequiredAcks(config.RequiredAck)
	kafkaConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	kafkaConfig.Producer.Return.Successes = true
	var err error
	kafkaClient, err = sarama.NewSyncProducer(config.Address, kafkaConfig)
	if err != nil {
		logrus.Fatalf("new kafka client err=%+v", err)
	}
	logrus.Infof("success init kafka!")
}

// SendMsgAll 异步发送消息到kafka
func SendMsgAll(tailfMap map[string]*Tailf) {
	for _, t := range tailfMap {
		go sendMsg(t)
	}
}

// sendMsg 发送单个topic的消息到kafka
func sendMsg(t *Tailf) {
	for {
		select {
		case msg := <-t.GetMsg():
			_, _, err := kafkaClient.SendMessage(msg)
			if err != nil {
				logrus.Warnf("kafka send msg failed,topic=%s,msg=%+v,err=%+v\n", t.topic, msg, err)
				continue
			}
			logrus.Println("topic=%s,send msg=", t.topic, msg.Value)
		}
	}
}
