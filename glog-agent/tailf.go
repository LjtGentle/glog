package glog_agent

import (
	"github.com/IBM/sarama"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

// Tailf log agent handle struct
type Tailf struct {
	msgChan chan *sarama.ProducerMessage // 消息队列
	topic   string                       // kafka topic
	tail    *tail.Tail                   // 监听的文件
}

// InitTail 初始化 tail -f
func InitTail(agentConfig *AgentConfig) map[string]*Tailf {
	config := tail.Config{
		Location: &tail.SeekInfo{
			Offset: 0,
			Whence: 2,
		},
		ReOpen: true,
		Poll:   true,
		Follow: true,
	}
	tailfMap := make(map[string]*Tailf, len(agentConfig.Targets))
	for _, t := range agentConfig.Targets {
		tailFile, err := tail.TailFile(t.FileName, config)
		if err != nil {
			logrus.Fatalf("tail file err=%+v", err)
		}
		tailfMap[t.FileName] = &Tailf{
			msgChan: make(chan *sarama.ProducerMessage, t.ChanSize),
			tail:    tailFile,
			topic:   t.Topic,
		}
	}
	return tailfMap
}

// ReadAll 读取多个文件 tail -f 的信息
func ReadAll(tailfMap map[string]*Tailf) {
	for fileName, t := range tailfMap {
		go t.read(fileName)
	}
}

// Read 阻塞读取一个文件 tail -f 命令
func (t *Tailf) read(fileName string) {
	for {
		msg, ok := <-t.tail.Lines
		if !ok {
			logrus.Warnf("tailf lines failed,fileName=%s,err=%+v", fileName, msg.Err)
			time.Sleep(time.Second)
			continue
		}
		if len(strings.TrimSpace(msg.Text)) == 0 {
			continue
		}
		kafkaMsg := &sarama.ProducerMessage{
			Topic: t.topic,
			Value: sarama.StringEncoder(msg.Text),
		}
		t.msgChan <- kafkaMsg
		logrus.Infof("read text=%s,kafkaMsg=%+v", msg.Text, kafkaMsg)
	}
}

// GetMsg 获取消息
func (t *Tailf) GetMsg() <-chan *sarama.ProducerMessage {
	return t.msgChan
}
