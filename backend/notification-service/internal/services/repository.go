package services

import (
	"context"

	"notification-service/internal/models"
)

type Repository interface {
	GetPreference(context.Context, string) (*models.Preference, error)
	SavePreference(context.Context, models.Preference) error
	CreateRecord(context.Context, models.Record) (bool, error)
	UpdateRecord(context.Context, string, string, string) error
	GetRecord(context.Context, string) (*models.Record, error)
	ListRecords(context.Context, string, int, int) ([]models.Record, int64, error)
}

type EmailSender interface {
	SendWelcome(context.Context, string, string) error
}
