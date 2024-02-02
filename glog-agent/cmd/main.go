package main

import (
	glog_agent "glog-agent"
)

func main() {
	config := glog_agent.InitConfig()
	tailfMap := glog_agent.InitTail(config)
	glog_agent.InitKafka(config)
	go glog_agent.ReadAll(tailfMap)
	go glog_agent.SendMsgAll(tailfMap)
	select {}
}
