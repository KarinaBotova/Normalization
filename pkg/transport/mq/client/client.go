package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"

	pb "github.com/kolya59/easy_normalization/proto"
)

type Client struct {
	topic *pubsub.Topic
}

func NewClient(projectID, topicName string) (*Client, error) {
	ctx := context.Background()
	//создание клиента очереди
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}
	// из клиента получаем тему
	topic := client.Topic(topicName)

	// есть ли такая тема, проверка
	exists, err := topic.Exists(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check topic existense: %v", err)
	}
	if !exists { //если нет этой темы, то создаем ее
		if _, err = client.CreateTopic(ctx, topicName); err != nil {
			return nil, fmt.Errorf("failed to create topic")
		}
	}

	return &Client{topic: topic}, nil
}
//очередь
func (c *Client) SendStudents(Students []models.Student) error {
	data, err := json.Marshal(Students)
	if err != nil {
		return err
	}
	msg := &pubsub.Message{Data: data}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := c.topic.Publish(ctx, msg).Get(ctx); err != nil {
		return err
	}

	return nil
}
