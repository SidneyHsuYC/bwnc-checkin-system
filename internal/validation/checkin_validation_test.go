package validation

import (
	"testing"

	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/models"
)

func TestValidateCheckin(t *testing.T) {
	tests := []struct {
		name    string
		checkin *models.Checkin
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid checkin",
			checkin: &models.Checkin{
				StudentID: 1,
				EventID:   1,
			},
			wantErr: false,
		},
		{
			name: "missing student ID",
			checkin: &models.Checkin{
				EventID: 1,
			},
			wantErr: true,
			errMsg:  "student ID is required",
		},
		{
			name: "missing event ID",
			checkin: &models.Checkin{
				StudentID: 1,
			},
			wantErr: true,
			errMsg:  "event ID is required",
		},
		{
			name: "zero student ID",
			checkin: &models.Checkin{
				StudentID: 0,
				EventID:   1,
			},
			wantErr: true,
			errMsg:  "student ID is required",
		},
		{
			name: "zero event ID",
			checkin: &models.Checkin{
				StudentID: 1,
				EventID:   0,
			},
			wantErr: true,
			errMsg:  "event ID is required",
		},
		{
			name: "negative student ID",
			checkin: &models.Checkin{
				StudentID: -1,
				EventID:   1,
			},
			wantErr: true,
			errMsg:  "student ID is required",
		},
		{
			name: "negative event ID",
			checkin: &models.Checkin{
				StudentID: 1,
				EventID:   -1,
			},
			wantErr: true,
			errMsg:  "event ID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCheckin(tt.checkin)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCheckin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateCheckin() error = %v, want error containing %v", err, tt.errMsg)
				}
			}
		})
	}
}
