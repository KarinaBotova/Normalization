package server

import (
	"context"
	"encoding/json"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/rs/zerolog/log"

	postgresdriver "github.com/kolya59/easy_normalization/pkg/postgres-driver"
	pb "github.com/kolya59/easy_normalization/proto"
)

func StartServer(projectID, topicName, subName string, done chan interface{}) {
	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal().Msgf("failed to create client: %v", err)
	}

	topic := client.Topic(topicName)

	// Create the topic if it doesn't exist
	exists, err := topic.Exists(ctx)
	if err != nil {
		log.Fatal().Msgf("failed to check topic existence: %v", err)
	}
	if !exists {
		if _, err = client.CreateTopic(ctx, topicName); err != nil {
			log.Fatal().Msgf("failed to create topic: %v", err)
		}
	}

	// создаем подписку
	sub := client.Subscription(subName)
	exists, err = sub.Exists(ctx)
	if err != nil {
		log.Fatal().Msgf("failed to check sub existence: %v", err)
	}
	if !exists {
		if _, err = client.CreateSubscription(ctx, subName, pubsub.SubscriptionConfig{
			Topic:       topic,
			AckDeadline: 10 * time.Second,
		}); err != nil {
			log.Fatal().Msgf("failed to create sub: %v", err)
		}
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		<-done
		ctx.Done()
	}()
	if err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		// Get data
		var Students []models.Student
		if err := json.Unmarshal(msg.Data, &Students); err != nil {
			log.Fatal().Err(err).Msg("Failed to unmarshal Students")
		}
		// Send data in DB
		if err := postgresdriver.SaveStudents(Students); err != nil {
			log.Error().Err(err).Msg("Could not send Students to DB")
		}
		log.Info().Msgf("Students %v was saved via Google Pub/Sub", Students)
		msg.Ack()
	}); err != nil {
		log.Fatal().Msgf("receive error: %v", err)
	}
}
