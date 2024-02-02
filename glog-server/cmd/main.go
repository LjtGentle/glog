package main

import (
	"context"
	glog_server "glog-server"
)

func main() {
	config := glog_server.InitConfig()
	ks := glog_server.InitKafka(&config.KafkaConfig)
	glog_server.Do(ks)
	glog_server.InitElastic(&config.ElasticConfig, ks)
	handle(context.TODO())
}

// handle 阻塞主线程
func handle(ctx context.Context) {
	select {
	case <-ctx.Done():
		break
	}
}
