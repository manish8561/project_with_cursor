package services

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"notification-service/internal/models"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type EventConsumer struct {
	reader  *kafka.Reader
	service *NotificationService
	logger  *zap.Logger
}

func NewEventConsumer(brokers, groupID, clientID, topic string, service *NotificationService, logger *zap.Logger) (*EventConsumer, error) {
	if strings.TrimSpace(brokers) == "" {
		return nil, nil
	}
	brokerList := strings.Split(brokers, ",")
	for i := range brokerList {
		brokerList[i] = strings.TrimSpace(brokerList[i])
	}
	if groupID == "" || topic == "" {
		return nil, errors.New("kafka group id and user-created topic are required")
	}
	reader := kafka.NewReader(kafka.ReaderConfig{Brokers: brokerList, GroupID: groupID, Topic: topic, MinBytes: 1, MaxBytes: 10e6, Dialer: &kafka.Dialer{ClientID: clientID}})
	return &EventConsumer{reader: reader, service: service, logger: logger}, nil
}

func (c *EventConsumer) Start(ctx context.Context) {
	if c == nil || c.reader == nil {
		return
	}
	for {
		message, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			c.logger.Error("failed to fetch user-created event", zap.Error(err))
			continue
		}
		var event models.UserEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			c.logger.Error("invalid user-created event", zap.Error(err))
		} else if err := c.service.HandleUserCreated(ctx, event); err != nil {
			c.logger.Error("failed to persist notification outcome", zap.Error(err), zap.String("event_id", event.EventID))
			continue
		}
		if err := c.reader.CommitMessages(ctx, message); err != nil {
			c.logger.Error("failed to commit user-created event", zap.Error(err))
		}
	}
}

func (c *EventConsumer) Close() error {
	if c == nil || c.reader == nil {
		return nil
	}
	return c.reader.Close()
}
