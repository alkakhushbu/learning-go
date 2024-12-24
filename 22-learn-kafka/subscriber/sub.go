package main

import (
	"context"
	"fmt"
	"github.com/twmb/franz-go/pkg/kgo"
	"time"
)

func main() {
	// Specify the Kafka seed brokers (addresses of Kafka servers)
	seeds := []string{"kafka-service:9092"}

	// Initialize a Kafka client with both producer and consumer capabilities
	// The client joins a consumer group and consumes messages from the "hello" topic.
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),                // Specify the Kafka brokers for connectivity
		kgo.ConsumerGroup("my-group-identifier"), // Use a consumer group ("my-group-identifier") for consuming messages
		kgo.ConsumeTopics("hello"),               // Specify the topic this client will consume from ("hello")
		kgo.FetchMinBytes(1),                     // Set the minimum number of bytes to fetch before responding (1 byte)
		kgo.FetchMaxWait(10*time.Millisecond),    // Set the maximum wait time for fetch requests (10 ms)

		// Uncomment the following line to disable auto-commit for manual offset management
		// kgo.DisableAutoCommit(),
	)
	if err != nil { // Check if there's an error initializing the client
		panic(err) // Terminate the program with the error
	}
	defer cl.Close() // Ensure the client is closed properly when the program exits

	// Create a context to manage cancellation or timeouts for fetch operations
	ctx := context.Background()

	// Infinite loop to continuously poll for new Kafka records/messages
	for {
		// PollFetches retrieves messages from the specified Kafka topic(s)
		fetches := cl.PollFetches(ctx)
		fmt.Println("fetched") // Indicate a fetch operation was performed

		// Check if there were any errors during the fetch process
		if errs := fetches.Errors(); len(errs) > 0 {
			// Kafka retries internally for retriable errors, but non-retriable errors
			// are returned here so the user can take action
			fmt.Println("errors", errs) // Log any errors that occurred
			fmt.Println("waiting for Kafka to come up again")
			time.Sleep(5 * time.Second) // Wait 5 seconds before retrying
			continue                    // Continue polling for messages
		}

		// Retrieve records returned by the fetch operation using a record iterator
		iter := fetches.RecordIter()

		// Iterate through all the records fetched
		for !iter.Done() {
			// Get the next record from the iterator
			record := iter.Next()

			// Print the value (message payload) of the record
			fmt.Println(string(record.Value), "from an iterator!")
		}
	}
}
