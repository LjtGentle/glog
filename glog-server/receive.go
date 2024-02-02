package glog_server

import (
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type KafkaImpl struct {
	topic      string
	msgChan    chan *string
	partitions []int32
	consumer   sarama.Consumer
}

func InitKafka(cfg *KafkaConfig) []KafkaImpl {
	consumer, err := sarama.NewConsumer(cfg.Address, nil)
	if err != nil {
		logrus.Fatalf("new kafka consumer failed,addrs=%+v,err=%+v", cfg.Address, err)
	}
	ks := make([]KafkaImpl, 0, len(cfg.Targets))
	for _, t := range cfg.Targets {
		partitions, err := consumer.Partitions(t.Topic)
		if err != nil {
			logrus.Fatalf("get kafka partitions failed, topic=%s,err=%+v", t.Topic, err)
		}
		k := KafkaImpl{
			topic:      t.Topic,
			msgChan:    make(chan *string, t.ChanSize),
			partitions: partitions,
			consumer:   consumer,
		}
		ks = append(ks, k)
	}
	return ks
}
func Do(ks []KafkaImpl) {
	logrus.Infof("in do ks=%+v", ks)
	for i := 0; i < len(ks); i++ {
		ks[i].listen()
	}
}

func (k *KafkaImpl) listen() {
	for _, partition := range k.partitions {
		pc, err := k.consumer.ConsumePartition(k.topic, partition, sarama.OffsetNewest)
		if err != nil {
			logrus.Errorf("get kafka pc failed,topic=%s,partition=%d,err=%+v", k.topic, partition, err)
			continue
		}
		logrus.Infof("get kafka consumer pc=%+v", pc)
		go func(pConsumer sarama.PartitionConsumer) {
			defer pConsumer.AsyncClose()
			for msg := range pConsumer.Messages() {
				msgStr := string(msg.Value)
				k.sendMsg(&msgStr)
				logrus.Infof("get kafka msg value=%+v,topic=%s", msg.Value, msg.Topic)
			}
		}(pc)
	}
}

func (k *KafkaImpl) sendMsg(msg *string) {
	k.msgChan <- msg
}

func (k *KafkaImpl) GetMsg() <-chan *string {
	return k.msgChan
}
