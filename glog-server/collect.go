package glog_server

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/sirupsen/logrus"
)

func InitElastic(config *ElasticConfig, ks []KafkaImpl) {
	// ES 配置
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
		//Username: "elasticsearch",
		//Password: "your_password",
	}
	client, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		logrus.Fatalf("new elastic client fialed, addr=%s,err=%+v", config.Addr, err)
		return
	}
	for _, k := range ks {
		go createLog(k, client, config.Index)
	}
	logrus.Info("init elastic success!")
}

func createLog(k KafkaImpl, client *elasticsearch.TypedClient, index string) {
	logrus.Infof("in create log kafkaImpl=%+v", k)
	for msg := range k.GetMsg() {
		m := map[string]string{
			"message": k.topic + ":" + *msg,
		}
		indexRes, err := client.Index(index).Request(m).Do(context.TODO())
		if err != nil {
			logrus.Errorf("create es index failed,index=%s,msg=%+v,err=%+v", index, msg, err)
			continue
		}
		logrus.Infof("success create es index,indexRes=%+v", indexRes)
	}
}
