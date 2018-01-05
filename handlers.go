package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
    //"strconv"
    //"fmt"
)

func AddConsumerTopic(c *gin.Context) {
    var param ConsumerParam
    if err := c.ShouldBindJSON(&param); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    }
    err := consumeKafkaAll(param) // Add new Kafka topic 
    if err==nil{
        c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
    } else {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    }
}

func DeleteConsumerTopic(c *gin.Context) {
    var param ConsumerParam
    if err := c.ShouldBindJSON(&param); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    }
    deleteConsumerTopic(param) // Add new Kafka topic 
    c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

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

