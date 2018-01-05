package main

import (
        "fmt"
        "github.com/confluentinc/confluent-kafka-go/kafka"
        "os"
        "encoding/json"
        //"syscall"
        //"os/signal"
)

var producer *kafka.Producer
//var consumer *kafka.Consumer
//var sigchan chan os.Signal
var topicList topicType
//var consumerRun map[string]bool
//var consumerRun map[string]chan bool

var (
    brokers = "10.148.0.4:9092"
    consumerRun = make(map[string]bool)
)

type topicType[]string
type ConsumerParam struct {
    Topic string `json:topic`
    Group string `json:group`
}

func (l *topicType) removeElement(item string) {
    l1 := *l
    for i, other := range *l {
        if other == item {
            l1 = append(l1[:i], l1[i+1:]...)
        }
    }
    *l = l1
}

func(l *topicType) addElement(item string) {
    l1 := *l
    l1.removeElement(item)
    l1 = append(l1, item)
    *l = l1
}

func InitKafkaProducer() (err error) {
        //brokers = "10.148.0.4:9092"
        producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": brokers})
        if err!=nil{
                os.Exit(1)
        }
        return
}

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

func deleteConsumerTopic(param ConsumerParam) {
	consumerRun[param.Topic] = false
}

func consumeKafkaAll(param ConsumerParam) (err error) {
        var data interface{}
        topic := param.Topic
        group := param.Group

	//consumerRun[topic] = make(chan bool)
        run := make(chan bool)
        //sigchan := make(chan os.Signal, 1)
        //signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

        consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
                "bootstrap.servers":               brokers,
                "group.id":                        group,
                "session.timeout.ms":              6000,
                "go.events.channel.enable":        true,
                "go.application.rebalance.enable": true,
                "default.topic.config":            kafka.ConfigMap{"auto.offset.reset": "earliest"}})   
        if err!=nil{
                fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
                os.Exit(1)
        }

        err = consumer.SubscribeTopics([]string{topic}, nil)
        if err!=nil{
                os.Exit(2)
        }

	consumerRun[topic] = true
        go func() {
                defer close(run)
                for {
		    ev := <-consumer.Events()
                    //select {
                    //case run := <- consumerRun[topic]:
                    //      fmt.Printf("Terminating topic %s\n", topic)
                    //      consumerRun[topic] = false

                    //case ev := <-consumer.Events():
                    //ev := <-consumer.Events()
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
			    consumerRun[topic] = false
                    }
                }
        }()

	go func() {
		for {
		if consumerRun[topic]==false {
			_ = <- run
			fmt.Fprintf(os.Stderr, "%% Remove consumer %s\n", topic)
		}
		}
	}()

        //fmt.Printf("Closing consumer\n")
        //consumer.Close()
        return
}
