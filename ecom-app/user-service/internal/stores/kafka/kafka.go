package kafka

import (
	"context"
	"fmt"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"os"
	"time"
)

type Conf struct {
	client *kgo.Client
	admin  *kadm.Client
}

func NewConf() (*Conf, error) {
	host := os.Getenv("KAFKA_HOST")
	port := os.Getenv("KAFKA_PORT")

	if host == "" || port == "" {
		return nil, fmt.Errorf("kafka host or port is empty")
	}
	connString := fmt.Sprintf("%s:%s", host, port)
	var err error
	var client *kgo.Client
	for i := 1; i < 8; i++ {

		var backoff time.Duration = 2
		client, err = kgo.NewClient(
			kgo.SeedBrokers(connString),

			//ProducerLinger sets how long individual topic partitions will linger waiting for more records
			//before triggering a request to be built.
			kgo.ProducerLinger(0),
		)
		if err != nil {
			fmt.Printf("kafka client error: %v", err)
			time.Sleep(backoff * time.Second)
			backoff *= 2
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
		defer cancel()
		err = client.Ping(ctx)

		if err != nil {
			fmt.Printf("kafka client error: %v", err)
			time.Sleep(backoff * time.Second)
			backoff *= 2
			continue
		}

		break
	}

	if err != nil {
		return nil, fmt.Errorf("kafka client error: %v", err)
	}
	admin := kadm.NewClient(client)
	return &Conf{
		client: client,
		admin:  admin,
	}, nil
}
