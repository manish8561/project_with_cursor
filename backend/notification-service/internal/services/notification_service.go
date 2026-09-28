package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"notification-service/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationService struct {
	repository Repository
	email      EmailSender
}

func NewNotificationService(repository Repository, email EmailSender) *NotificationService {
	return &NotificationService{repository: repository, email: email}
}

func (s *NotificationService) Preference(ctx context.Context, userID string) (*models.Preference, error) {
	return s.repository.GetPreference(ctx, userID)
}

func (s *NotificationService) UpdatePreference(ctx context.Context, userID string, enabled bool) (*models.Preference, error) {
	preference := models.Preference{UserID: userID, EmailEnabled: enabled, UpdatedAt: time.Now().UTC()}
	if err := s.repository.SavePreference(ctx, preference); err != nil {
		return nil, err
	}
	return &preference, nil
}

func (s *NotificationService) History(ctx context.Context, userID string, page, size int) (*models.HistoryResponse, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	items, total, err := s.repository.ListRecords(ctx, userID, page, size)
	if err != nil {
		return nil, err
	}
	return &models.HistoryResponse{Items: items, Total: total, Page: page, Size: size}, nil
}

func (s *NotificationService) HandleUserCreated(ctx context.Context, event models.UserEvent) error {
	event.EventType = strings.TrimSpace(event.EventType)
	if event.EventType != "user.created.v1" || event.EventID == "" || event.UserID == "" || event.Email == "" {
		return errors.New("invalid user-created event")
	}
	preference, err := s.repository.GetPreference(ctx, event.UserID)
	if err != nil {
		return err
	}
	_, shouldDeliver, err := s.ensureRecord(ctx, event)
	if err != nil || !shouldDeliver {
		return err
	}
	if !preference.EmailEnabled {
		return s.repository.UpdateRecord(ctx, event.EventID, "skipped", "")
	}
	return s.deliverWelcome(ctx, event)
}

func (s *NotificationService) ensureRecord(ctx context.Context, event models.UserEvent) (models.Record, bool, error) {
	now := time.Now().UTC()
	record := models.Record{
		ID: primitive.NewObjectID().Hex(), EventID: event.EventID, UserID: event.UserID,
		Type: "welcome", Channel: "email", Status: "queued", CreatedAt: now,
		UpdatedAt: now, ExpireAt: now.Add(90 * 24 * time.Hour),
	}
	created, err := s.repository.CreateRecord(ctx, record)
	if err != nil {
		return models.Record{}, false, err
	}
	if !created {
		existing, getErr := s.repository.GetRecord(ctx, event.EventID)
		if getErr != nil {
			return models.Record{}, false, getErr
		}
		if existing.Status == "sent" || existing.Status == "skipped" {
			return *existing, false, nil
		}
		record = *existing
	}
	return record, true, nil
}

func (s *NotificationService) deliverWelcome(ctx context.Context, event models.UserEvent) error {
	var err error
	if s.email == nil {
		err = errors.New("email sender is unavailable")
	} else {
		for attempt := 0; attempt < 3; attempt++ {
			err = s.email.SendWelcome(ctx, event.Email, event.Name)
			if err == nil || ctx.Err() != nil {
				break
			}
			if attempt < 2 {
				select {
				case <-ctx.Done():
					err = ctx.Err()
				case <-time.After(time.Duration(attempt+1) * 250 * time.Millisecond):
				}
			}
		}
	}
	if err != nil {
		_ = s.repository.UpdateRecord(ctx, event.EventID, "failed", "delivery failed")
		return nil
	}
	return s.repository.UpdateRecord(ctx, event.EventID, "sent", "")
}
