# Confluent-Kafka-Client using Go

## Preparation

### Install librdkafka
```
cd ~
git clone https://github.com/edenhill/librdkafka.git
sudo mv librdkafka /opt/
```

```
cd /opt/librdkafka
./configure --prefix /usr
make
sudo make install
echo 'export PKG_CONFIG_PATH="/opt/librdkafka/src"' >> /etc/profile.d/custom.sh
source /etc/profile.d/custom.sh
cd ~
```

### Install all required packages
```
go get -u github.com/confluentinc/confluent-kafka-go/kafka
go get -u github.com/gin-gonic/gin
go get -u gopkg.in/resty.v1
```

## Preparation
(1) Edit kafka-pubsub.go
- Variable **brokers** to point to the desired Kafka broker clusters address
- Function **useConsumer()** to customize your desired commands for consumed data from Kafka

(2) Edit main.go
- Change the url route

## Compile & run
```
#go run main.go
go build
./kafka-pubsub
```

## API Specifications

### Get list of topics of the consumer
```
curl -X GET http://localhost:8020/topic
```

### Add new topic for kafka consumer to be running forever as a goroutine
```
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{"topic":"test2","group":"test-group"}' http://localhost:8020/subscribe/add
```

### Add microservices URL as kafka subscriber via REST
```
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{"data" : [{"topic":"test1","url":["http://0.0.0.0:8000","http://127.0.0.1:8000"]}, {"topic":"test2","url":["http://0.0.0.0:8080","http://127.0.0.1:8080"]}]}' http://localhost:8020/subscribe/url/add
```

### Delete microservices URL from kafka subscriber via REST
```
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{"data" : [{"topic":"test1","url":["http://0.0.0.0:8000","http://127.0.0.1:8000"]}, {"topic":"test2","url":["http://0.0.0.0:8080","http://127.0.0.1:8080"]}]}' http://localhost:8020/subscribe/url/delete
```

### Publish data to Kafka
```
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d @./json/example.json http://localhost:8020/publish/<topic-name>
```

### Remove a topic, stop its consumer goroutine
```
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{"topic":"test2","group":"test-group"}' http://localhost:8020/subscribe/delete
```
