package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
    //"strconv"
    //"fmt"
)

/*
func UpdateTopics(c *gin.Context) {
    topicAll,err := updateTopics("topics.txt")
    if err==nil {
	topics = topicAll
	//consumer.SubscribeTopics(topics, nil)
	for _, element := range someSlice {
	    go consumeKafkaAll(element)
	}
        c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
    } else {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    }
}
*/

func RunProduce(c *gin.Context) {
    var data    interface{}
    topic := c.Param("topic")
    if err := c.ShouldBindJSON(&data); err == nil {
	produceKafka(topic, data) // Produce data to Kafka topic
        c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
    } else {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    }
}

