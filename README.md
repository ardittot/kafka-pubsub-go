# Confluent-Kafka-Client using Go

This example is using Ubuntu 16.04 (LTS)

## Preparation

Update Ubuntu
```
sudo apt-get -f update
sudo apt-get -y upgrade
```

Create a new shell bootstrap script _/etc/profile.d/custom.sh_
```
sudo touch /etc/profile.d/custom.sh
```

### Install Go
```
sudo curl -O https://storage.googleapis.com/golang/go1.9.2.linux-amd64.tar.gz
sudo tar -xvf go1.9.2.linux-amd64.tar.gz
sudo mv go /usr/local
sudo ln -s /usr/local/go /opt/go
```

Create a new directory for Go third party packages
```
sudo mkdir $GOROOT/work
sudo chmod 777 /opt/go/work
```

Setup environmet variables for Go
> export GOROOT=/opt/go
>
> export PATH=$PATH:$GOROOT/bin
>
> export GOPATH=/opt/go/work

Source the shell bootstrap script
```
source /etc/profile.d/custom.sh
```

Update the required dependencies for librdkafka (explained next step)
```
sudo apt-get install -y pkg-config lxc-dev
```

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
```

Insert this to _/etc/profile.d/custom.sh_
> export PKG_CONFIG_PATH="/opt/librdkafka/src"

```
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
curl -X GET http://localhost:8020/subscribe/topic
```

### Add new topic for kafka consumer 
It will keep consuming data from kafka topic forever as a goroutine, and pass the data to any microservices (use the function useConsumer() inside kafka-pubsub.go)
```
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{"topic":"test2","group":"test-group"}' http://localhost:8020/subscribe/topic/add
```

### Remove a topic from kafka consumer
Stop its consumer goroutine
```
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{"topic":"test2"}' http://localhost:8020/subscribe/topic/delete
```

### Add microservices URL as kafka subscriber via REST
```
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{"data" : [{"topic":"test1","url":["http://0.0.0.0:8000","http://127.0.0.1:8000"]}, {"topic":"test2","url":["http://0.0.0.0:8080","http://127.0.0.1:8080"]}]}' http://localhost:8020/subscribe/url/add
```

### Delete microservices URL from kafka subscriber
```
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{"data" : [{"topic":"test1","url":["http://0.0.0.0:8000","http://127.0.0.1:8000"]}, {"topic":"test2","url":["http://0.0.0.0:8080","http://127.0.0.1:8080"]}]}' http://localhost:8020/subscribe/url/delete
```

### Publish data to Kafka
```
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d @./json/example.json http://localhost:8020/publish/<topic-name>
```
