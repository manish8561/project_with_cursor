package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"notification-service/internal/models"
)

type memoryRepository struct {
	preference *models.Preference
	records    map[string]models.Record
	total      int64
}

func newMemoryRepository() *memoryRepository {
	return &memoryRepository{records: make(map[string]models.Record)}
}

func (r *memoryRepository) GetPreference(_ context.Context, userID string) (*models.Preference, error) {
	if r.preference == nil {
		return &models.Preference{UserID: userID, EmailEnabled: true}, nil
	}
	return r.preference, nil
}

func (r *memoryRepository) SavePreference(_ context.Context, preference models.Preference) error {
	r.preference = &preference
	return nil
}

func (r *memoryRepository) CreateRecord(_ context.Context, record models.Record) (bool, error) {
	if _, exists := r.records[record.EventID]; exists {
		return false, nil
	}
	r.records[record.EventID] = record
	return true, nil
}

func (r *memoryRepository) UpdateRecord(_ context.Context, eventID, status, summary string) error {
	record, exists := r.records[eventID]
	if !exists {
		return errors.New("record not found")
	}
	record.Status = status
	record.ErrorSummary = summary
	record.UpdatedAt = time.Now().UTC()
	r.records[eventID] = record
	return nil
}

func (r *memoryRepository) GetRecord(_ context.Context, eventID string) (*models.Record, error) {
	record, exists := r.records[eventID]
	if !exists {
		return nil, errors.New("record not found")
	}
	return &record, nil
}

func (r *memoryRepository) ListRecords(_ context.Context, _ string, page, size int) ([]models.Record, int64, error) {
	return []models.Record{{ID: "record-1", Status: "sent"}}, r.total, nil
}

type recordingEmailSender struct {
	calls int
	err   error
}

func (s *recordingEmailSender) SendWelcome(context.Context, string, string) error {
	s.calls++
	return s.err
}

func testEvent() models.UserEvent {
	return models.UserEvent{EventID: "event-1", EventType: "user.created.v1", UserID: "user-1", Email: "user@example.com", Name: "Ada"}
}

func TestHandleUserCreatedHonorsDisabledPreference(t *testing.T) {
	repository := newMemoryRepository()
	repository.preference = &models.Preference{UserID: "user-1", EmailEnabled: false}
	sender := &recordingEmailSender{}
	service := NewNotificationService(repository, sender)

	if err := service.HandleUserCreated(context.Background(), testEvent()); err != nil {
		t.Fatal(err)
	}
	if sender.calls != 0 {
		t.Fatalf("expected no email sends, got %d", sender.calls)
	}
	if got := repository.records["event-1"].Status; got != "skipped" {
		t.Fatalf("expected skipped status, got %q", got)
	}
}

func TestHandleUserCreatedSendsOnceForDuplicateEvent(t *testing.T) {
	repository := newMemoryRepository()
	sender := &recordingEmailSender{}
	service := NewNotificationService(repository, sender)

	for range 2 {
		if err := service.HandleUserCreated(context.Background(), testEvent()); err != nil {
			t.Fatal(err)
		}
	}
	if sender.calls != 1 {
		t.Fatalf("expected one email send, got %d", sender.calls)
	}
	if got := repository.records["event-1"].Status; got != "sent" {
		t.Fatalf("expected sent status, got %q", got)
	}
}

func TestHandleUserCreatedRetriesAndRecordsFailure(t *testing.T) {
	repository := newMemoryRepository()
	sender := &recordingEmailSender{err: errors.New("smtp unavailable")}
	service := NewNotificationService(repository, sender)

	if err := service.HandleUserCreated(context.Background(), testEvent()); err != nil {
		t.Fatal(err)
	}
	if sender.calls != 3 {
		t.Fatalf("expected three bounded attempts, got %d", sender.calls)
	}
	if got := repository.records["event-1"].Status; got != "failed" {
		t.Fatalf("expected failed status, got %q", got)
	}
}

func TestHistoryNormalizesPagination(t *testing.T) {
	service := NewNotificationService(newMemoryRepository(), nil)
	response, err := service.History(context.Background(), "user-1", 0, 500)
	if err != nil {
		t.Fatal(err)
	}
	if response.Page != 1 || response.Size != 20 {
		t.Fatalf("expected page 1 size 20, got page %d size %d", response.Page, response.Size)
	}
}
