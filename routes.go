package main

func initializeRoutes() {

  // Handle the index route
  router.GET("/subscribe/topic", GetTopics)
  router.POST("/subscribe/topic/add", AddConsumerTopic)
  router.POST("/subscribe/topic/delete", DeleteConsumerTopic)
  router.POST("/subscribe/url/add", AddConsumerUrl)
  router.POST("/subscribe/url/delete", DeleteConsumerUrl)
  router.POST("/publish/:topic", RunProduce)
}
