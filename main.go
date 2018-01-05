package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
)

var router *gin.Engine

func main() {

  // Set the router as the default one provided by Gin
  router = gin.Default()

  // Initialize kafka
  InitKafkaProducer()
  //InitKafkaConsumer()

  // Initialize the routes
  initializeRoutes()

  // Start serving the application
  router.Run("0.0.0.0:8020")
  fmt.Printf("tes")
}
