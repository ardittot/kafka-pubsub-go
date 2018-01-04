package main

func initializeRoutes() {

  // Handle the index route
  //router.GET("/update", UpdateTopics)
  router.POST("/publish/:topic", RunProduce)
}
