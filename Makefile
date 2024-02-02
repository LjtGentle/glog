# 目标是简化部署流程、运行流程，点单
# 需要准备的东西 docker kafka es  zookeeper

init :
	make kafka
	make es

kafka :
	docker pull bitnami/zookeeper
	docker pull bitnami/kafka
	mkdir -p ~/data/zookeeper/data
	mkdir -p ~/data/kafka/logs
	docker run -d --name zookeeper -p 2181:2181 --restart=always -v ~/data/zookeeper/data:/opt/zookeeper-3.4.13/data -v /etc/localtime:/etc/localtime -e ALLOW_ANONYMOUS_LOGIN=true -t bitnami/zookeeper
	docker run -d --name kafka -p 9092:9092 --restart=always -v ~/data/kafka/logs:/opt/kafka/logs -v /etc/localtime:/etc/localtime \
 	--link zookeeper \
 	-e KAFKA_BROKER_ID=0 -e KAFKA_CFG_ZOOKEEPER_CONNECT=zookeeper:2181 -e KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092 -e KAFKA_CFG_LISTENERS=PLAINTEXT://0.0.0.0:9092 -e ALLOW_PLAINTEXT_LISTENER=yes -t bitnami/kafka

# 密钥登陆的问题待解决
es:
	docker pull docker.elastic.co/elasticsearch/elasticsearch:8.12.0
	docker run --name es01 --net elastic -p 9200:9200 \
	-d -m 1GB docker.elastic.co/elasticsearch/elasticsearch:8.12.0
	docker exec -it es01 /usr/share/elasticsearch/bin/elasticsearch-reset-password -u elastic
	docker exec -it es01 /usr/share/elasticsearch/bin/elasticsearch-create-enrollment-token -s kibana
	export ELASTIC_PASSWORD="your_password"
	docker cp es01:/usr/share/elasticsearch/config/certs/http_ca.crt .
	curl --cacert http_ca.crt -u elastic:$ELASTIC_PASSWORD https://localhost:9200
#	docker exec -it es01 /usr/share/elasticsearch/bin/elasticsearch-create-enrollment-token -s node
#	docker run -e ENROLLMENT_TOKEN="<token>" --name es02 --net elastic -it -m 1GB docker.elastic.co/elasticsearch/elasticsearch:8.12.0
#	curl --cacert http_ca.crt -u elastic:$ELASTIC_PASSWORD https://localhost:9200/_cat/nodes

#.IGNORE:
reset:
# 需要创建网桥，但是重复创建会报错
	docker network create elastic