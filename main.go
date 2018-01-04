package main

import (
	"github.com/gin-gonic/gin"
	//"fmt"
)

var router *gin.Engine

func main() {

  // Set the router as the default one provided by Gin
  router = gin.Default()

  // Initialize kafka
  InitKafkaProducer()
  //InitKafkaConsumer()

  // Run kafka consumer asynchronously
  /*
  go func() {
    _, _ := consumeKafka("test2")
    fmt.Printf("Message:\n%s\n", string(data_byte))
  }()
  */
  go consumeKafkaAll("test2")

  // Initialize the routes
  initializeRoutes()

  // Start serving the application
  router.Run("0.0.0.0:8010")

}
