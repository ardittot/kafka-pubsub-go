package main

import (
        "fmt"
        "github.com/confluentinc/confluent-kafka-go/kafka"
        "os"
        "encoding/json"
	"syscall"
	"os/signal"
)

var producer *kafka.Producer
//var consumer *kafka.Consumer
var broker string
//var topics []string
//var sigchan chan os.Signal

func InitKafkaProducer() (err error) {
	broker = "10.148.0.4:9092"
	producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err!=nil{
		os.Exit(1)
	}
	return
}

/*
func InitKafkaConsumer() (err error) {
	broker = "10.148.0.4:9092"
	group := "test2-group"

	sigchan = make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":               broker,
		"group.id":                        group,
		"session.timeout.ms":              6000,
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
		"default.topic.config":            kafka.ConfigMap{"auto.offset.reset": "earliest"}})	
	if err!=nil{
		os.Exit(1)
	}
	return
}
*/

func produceKafka(topic string, data interface{}) {
        value, err := json.Marshal(data)
	if err == nil {
		deliveryChan := make(chan kafka.Event)
		err = producer.Produce(&kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny}, Value: []byte(value)}, deliveryChan)

		e := <-deliveryChan
		m := e.(*kafka.Message)
		if m.TopicPartition.Error != nil {
			fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
		}
		close(deliveryChan)
	}
}

/*
func consumeKafka(topic string) {
	var data interface{}

	topics := []string{topic}
	err := consumer.SubscribeTopics(topics, nil)
	if err!=nil{
		os.Exit(2)
	}

	run := true
	for run == true {
	select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false

		case ev := <-consumer.Events():
			switch e := ev.(type) {
			case kafka.AssignedPartitions:
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				consumer.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				consumer.Unassign()
			case *kafka.Message:
				//consumer.Commit()
				data_byte := e.Value
				json.Unmarshal(data_byte, &data)
				fmt.Printf("Message:\n%+v\n", data)
				//fmt.Printf("Message:\n%s\n", string(data_byte))
				//fmt.Printf("%% Message on %s:\n%s\n", e.TopicPartition, string(e.Value))
			case kafka.PartitionEOF:
				fmt.Printf("%% Reached %v\n", e)
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
				run = false
			}
		}
	}
}
*/

func consumeKafkaAll(topic string) {
	var data interface{}

        group := "test2-group"

        sigchan := make(chan os.Signal, 1)
        signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

        consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
                "bootstrap.servers":               broker,
                "group.id":                        group,
                "session.timeout.ms":              6000,
                "go.events.channel.enable":        true,
                "go.application.rebalance.enable": true,
                "default.topic.config":            kafka.ConfigMap{"auto.offset.reset": "earliest"}})   
        if err!=nil{
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
                os.Exit(1)
        }

	topics := []string{topic}
	err = consumer.SubscribeTopics(topics, nil)
	if err!=nil{
		os.Exit(2)
	}

	run := true
	for run == true {
	select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false

		case ev := <-consumer.Events():
			switch e := ev.(type) {
			case kafka.AssignedPartitions:
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				consumer.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				consumer.Unassign()
			case *kafka.Message:
				//consumer.Commit()
				data_byte := e.Value
				json.Unmarshal(data_byte, &data)
				fmt.Printf("Message:\n%+v\n", data)
				//fmt.Printf("Message:\n%s\n", string(data_byte))
				//fmt.Printf("%% Message on %s:\n%s\n", e.TopicPartition, string(e.Value))
			case kafka.PartitionEOF:
				fmt.Printf("%% Reached %v\n", e)
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
				run = false
			}
		}
	}
	fmt.Printf("Closing consumer\n")
	consumer.Close()
}
