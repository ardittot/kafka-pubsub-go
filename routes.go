package main

func initializeRoutes() {

  // Handle the index route
  router.GET("/topic", GetTopics)
  router.POST("/subscribe/add", AddConsumerTopic)
  router.POST("/subscribe/delete", DeleteConsumerTopic)
  router.POST("/subscribe/url/add", AddConsumerUrl)
  router.POST("/subscribe/url/delete", DeleteConsumerUrl)
  router.POST("/publish/:topic", RunProduce)
}
