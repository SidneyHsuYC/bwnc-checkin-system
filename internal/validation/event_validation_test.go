package validation

import (
	"testing"
	"time"

	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/models"
)

func TestValidateEvent(t *testing.T) {
	futureTime := time.Now().Add(24 * time.Hour)

	tests := []struct {
		name    string
		event   *models.Event
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid event",
			event: &models.Event{
				EventName: "Test Event",
				EventTime: futureTime,
				EventType: "Workshop",
			},
			wantErr: false,
		},
		{
			name: "missing event name",
			event: &models.Event{
				EventTime: futureTime,
				EventType: "Workshop",
			},
			wantErr: true,
			errMsg:  "event_name is required",
		},
		{
			name: "missing event type",
			event: &models.Event{
				EventName: "Test Event",
				EventTime: futureTime,
			},
			wantErr: true,
			errMsg:  "event_type is required",
		},
		{
			name: "whitespace only event name",
			event: &models.Event{
				EventName: "   ",
				EventTime: futureTime,
				EventType: "Workshop",
			},
			wantErr: true,
			errMsg:  "event_name is required",
		},
		{
			name: "zero time",
			event: &models.Event{
				EventName: "Test Event",
				EventTime: time.Time{},
				EventType: "Workshop",
			},
			wantErr: true,
			errMsg:  "event_time is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEvent(tt.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateEvent() error = %v, want error containing %v", err, tt.errMsg)
				}
			}
		})
	}
}

func TestSanitizeEvent(t *testing.T) {
	eventTime := time.Now().Add(24 * time.Hour)

	tests := []struct {
		name     string
		input    *models.Event
		expected *models.Event
	}{
		{
			name: "trim whitespace from all fields",
			input: &models.Event{
				EventName: "  Test Event  ",
				EventType: "  Workshop  ",
				EventTime: eventTime,
			},
			expected: &models.Event{
				EventName: "Test Event",
				EventType: "Workshop",
				EventTime: eventTime,
			},
		},
		{
			name: "no changes needed",
			input: &models.Event{
				EventName: "Test Event",
				EventType: "Workshop",
				EventTime: eventTime,
			},
			expected: &models.Event{
				EventName: "Test Event",
				EventType: "Workshop",
				EventTime: eventTime,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SanitizeEvent(tt.input)
			if tt.input.EventName != tt.expected.EventName {
				t.Errorf("EventName = %q, want %q", tt.input.EventName, tt.expected.EventName)
			}
			if tt.input.EventType != tt.expected.EventType {
				t.Errorf("EventType = %q, want %q", tt.input.EventType, tt.expected.EventType)
			}
		})
	}
}
