package main

import (
	"context"
	"fmt"
	"github.com/twmb/franz-go/pkg/kgo"
	"time"
)

// connect kafka using kgo

// create a record using kgo.Record{key:value}

// produceRecord:=client.ProduceSync(ctx,record)

// produceRecord.First()

func main() {
	// Specify the Kafka seed brokers (addresses of Kafka servers)
	seeds := []string{"kafka-service:9092"}

	var client *kgo.Client // Declare a Kafka client variable
	var err error          // Declare an error variable to hold potential errors

	for i := 0; i < 5; i++ { // Retry up to 5 times
		// Initialize the Kafka client with configuration options
		client, err = kgo.NewClient(
			kgo.SeedBrokers(seeds...),    // Specify the Kafka brokers to connect to
			kgo.AllowAutoTopicCreation(), // Allow automatic topic creation if the topic doesn't exist
		)
		if err != nil { // Check if client initialization failed
			time.Sleep(2 * time.Second) // Wait for 2 seconds before retrying
			continue                    // Skip the rest of the loop and retry
		}

		// Check connectivity to Kafka broker by sending a Ping request
		err = client.Ping(context.Background())
		if err != nil { // If the ping fails
			time.Sleep(2 * time.Second) // Wait for 2 seconds before retrying
			continue                    // Retry the connection
		}

		// If the connection is successful, break out of the loop
		break
	}

	// If after all attempts there's still an error, terminate the program
	if err != nil {
		panic(err) // Terminate the program with the last error message
	}

	// If the client is nil (unexpected case), terminate the program
	if client == nil {
		panic("client is nil")
	}

	// Queue client closure when the program exits
	defer client.Close()

	// Create a Kafka record to send a message
	record := &kgo.Record{
		Topic: "hello",          // The Kafka topic to which the message will be sent
		Value: []byte("barbiz"), // The message payload to be sent
	}

	// Perform a synchronous Kafka produce operation to send the message
	pr := client.ProduceSync(context.Background(), record)

	// Retrieve the response for the first record
	rec, err := pr.First()
	if err != nil { // If there's an error while producing the message
		panic(err) // Terminate the program with the error
	}

	// Print the value of the successfully produced record
	fmt.Println(string(rec.Value))
}
