# Install all required GO packages
go get -u github.com/gin-gonic/gin
go get -u gopkg.in/resty.v1

# Compile & run
#go run main.go
go build
./kafka-pubsub

# Preparation
##Edit kafka-pubsub.go
##(1) Variable brokers to point to desired Kafka broker clusters
##(2) Function useConsumer():  to customize your desired commands for consumed data from Kafka

# Get list of topics of the consumer
curl -X GET http://localhost:8020/topic

# Add new topic for kafka consumer to be running forever as a goroutine
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{"topic":"test2","group":"test-group"}' http://localhost:8020/subscribe/add

# Publish data to Kafka
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d @./json/example.json http://localhost:8020/publish/<topic-name>

# Remove a topic, stop its consumer goroutine
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{"topic":"test2","group":"test-group"}' http://localhost:8020/subscribe/delete

