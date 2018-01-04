go get -u github.com/gin-gonic/gin
go run main.go
go build
./kafka-pubsub

# Update list of topics for async subscriber
#curl -X GET http://localhost:8010/update
# Publish data to Kafka
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d @./json/example.json http://localhost:8010/publish/<topic-name>
