package main

func initializeRoutes() {

  // Handle the index route
  router.GET("/topic", GetTopics)
  router.POST("/subscribe/add", AddConsumerTopic)
  router.POST("/subscribe/delete", DeleteConsumerTopic)
  router.POST("/publish/:topic", RunProduce)
}
