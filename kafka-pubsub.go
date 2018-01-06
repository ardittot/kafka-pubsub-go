package main

import (
        "fmt"
        "github.com/confluentinc/confluent-kafka-go/kafka"
        "os"
        "encoding/json"
        //"syscall"
        //"os/signal"
	"gopkg.in/resty.v1"
)

// *** Modify your function here ******************** //
func useConsumer(msg *kafka.Message, topic string) {
	var data interface{}
	urlList := consumerUrl[topic]
	msgVal := msg.Value
	json.Unmarshal(msgVal, &data)
	fmt.Printf("**********\nMessage:\n%+v\n", data)
	for _,url := range urlList {
	    fmt.Printf("Send to: %s\n", url)
	    go func() {
		//if len(url)>0 {
	            resty.R().
                    SetBody(data).
                    Post(url)
	        //}
	    }()
	}
	fmt.Printf("**********\n\n")
}
var (
    brokers = "10.148.0.4:9092"
)
// ************************************************** //

var producer *kafka.Producer
//var consumer *kafka.Consumer
var topicList topicType

var (
    consumerRun = make(map[string] (chan bool))
    consumerUrl = make(map[string] urlType)
)

type topicType []string
type urlType []string
type ConsumerParam struct {
    Topic string `json:topic`
    Group string `json:group`
}
type ConsumerUrl struct {
    Topic string `json:topic`
    Url []string `json:url`
}
type ConsumerUrlParam struct {
    Data []ConsumerUrl `json:data`
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

func (l *urlType) removeElement(item string) {
    l1 := *l
    for i, other := range *l {
        if other == item {
            l1 = append(l1[:i], l1[i+1:]...)
        }
    }
    *l = l1
}

func(l *urlType) addElement(item string) {
    l1 := *l
    l1.removeElement(item)
    l1 = append(l1, item)
    *l = l1
}

func addConsumerUrl(param ConsumerUrlParam) {
    arr := param.Data
    for _,par := range arr {
	topic := par.Topic
	urlList := par.Url
        for _,url := range urlList {
	    l := consumerUrl[topic]
	    l.removeElement(url)
	    l.addElement(url)
	    consumerUrl[topic] = l
	}
    }
}

func deleteConsumerUrl(param ConsumerUrlParam) {
    arr := param.Data
    for _,par := range arr {
	topic := par.Topic
	urlList := par.Url
        for _,url := range urlList {
	    l := consumerUrl[topic]
	    l.removeElement(url)
	    consumerUrl[topic] = l
	}
    }
}

func InitKafkaProducer() (err error) {
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
        consumerRun[param.Topic] <- true
	close(consumerRun[param.Topic])
	topicList.removeElement(param.Topic)
	fmt.Printf("** Deleting topic; Current topics: %v",topicList)
}

func consumeKafkaAll(param ConsumerParam) (err error) {
        topic := param.Topic
        group := param.Group
	exit := false
	consumerRun[topic] = make(chan bool)

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

	topicList.addElement(topic)
	fmt.Printf("** Adding topic; Current topics: %v\n\n",topicList)

        go func() {
                for {
                    select {
                    case <-consumerRun[topic]:
                        //fmt.Printf("Terminating topic %s\n", topic)
			exit = true
			//break

                    case ev := <-consumer.Events():
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
			    useConsumer(e, topic)
                            //data_byte := e.Value
                            //json.Unmarshal(data_byte, &data)
                            //fmt.Printf("Message:\n%+v\n", data)
                            //fmt.Printf("Message:\n%s\n", string(data_byte))
                            //fmt.Printf("%% Message on %s:\n%s\n", e.TopicPartition, string(e.Value))
                        case kafka.PartitionEOF:
                            fmt.Printf("%% Reached %v\n\n", e)
                        case kafka.Error:
                            fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
                        }
		    }
		    if exit==true {
		        break
		    }
                }
        	//fmt.Printf("Closing consumer\n")
        	consumer.Close()
        }()

        return
}
