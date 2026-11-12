package models

import (
	"fmt"
	errs "room-booking/internal/errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUpdateStatus(t *testing.T) {
	tests := []struct {
		name           string
		hourPrice      int64
		amount         int64
		from           time.Time
		to             time.Time
		actualStatus   BookingStatus
		expectedStatus BookingStatus
	}{
		{
			name:           "status confirmed",
			hourPrice:      100,
			amount:         300,
			from:           time.Date(2026, 12, 01, 10, 0, 0, 0, time.UTC),
			to:             time.Date(2026, 12, 01, 12, 0, 0, 0, time.UTC),
			expectedStatus: BookingStatusConfirmed,
		},
		{
			name:           "status pendind",
			hourPrice:      100,
			amount:         0,
			from:           time.Date(2026, 12, 01, 10, 0, 0, 0, time.UTC),
			to:             time.Date(2026, 12, 01, 12, 0, 0, 0, time.UTC),
			actualStatus:   BookingStatusConfirmed,
			expectedStatus: BookingStatusPending,
		},
		{
			name:           "status not changed",
			hourPrice:      100,
			amount:         100,
			from:           time.Date(2026, 12, 01, 10, 0, 0, 0, time.UTC),
			to:             time.Date(2026, 12, 01, 12, 0, 0, 0, time.UTC),
			actualStatus:   BookingStatusPending,
			expectedStatus: BookingStatusPending,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			booking := Booking{
				FromTime: tc.from,
				ToTime:   tc.to,
				Amount:   tc.amount,
			}

			booking.UpdateStatus(tc.hourPrice)

			require.Equal(t, tc.expectedStatus, booking.Status)
			if tc.actualStatus == BookingStatusConfirmed {
				require.NotEqual(t, tc.actualStatus, booking.Status)
			}
		})
	}
}

func TestCompletedStatus(t *testing.T) {
	tests := []struct {
		name           string
		actualStatus   BookingStatus
		expectedStatus BookingStatus
	}{
		{
			name:           "status changed to Completed",
			actualStatus:   BookingStatusConfirmed,
			expectedStatus: BookingStatusCompleted,
		},
		{
			name:         "invalid status transition",
			actualStatus: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			booking := Booking{
				Status: tc.actualStatus,
			}
			err := booking.Completed()
			if tc.actualStatus == BookingStatusConfirmed {
				require.Equal(t, tc.expectedStatus, booking.Status)
				require.NoError(t, err)
			}
			if tc.actualStatus == "" {
				require.Error(t, errs.ErrInvalidStatusTransition)
			}
		})
	}
}

func TestValidateTitle(t *testing.T) {
	tests := []struct {
		name        string
		actualValue string
		wantError   error
	}{
		{
			name:        "not null value",
			actualValue: "Booking title",
			wantError:   nil,
		},
		{
			name:        "null value",
			actualValue: "null",
			wantError:   errs.ErrInvalidTitle,
		},
		{
			name:        "empty value",
			actualValue: "",
			wantError:   errs.ErrInvalidTitle,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			booking := Booking{
				Title: tc.actualValue,
			}
			err := booking.validateTitle()
			if tc.wantError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.wantError)
			} else {
				require.NoError(t, err)
			}

		})
	}
}
func TestValidateDescription(t *testing.T) {
	tests := []struct {
		name        string
		actualValue string
		wantErr     error
	}{
		{
			name:        "not null value",
			actualValue: "Booking description",
			wantErr:     nil,
		},
		{
			name:        "null value",
			actualValue: "null",
			wantErr:     errs.ErrInvalidDescription,
		},
		{
			name:        "empty value",
			actualValue: "",
			wantErr:     errs.ErrInvalidDescription,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			booking := Booking{
				Description: tc.actualValue,
			}
			err := booking.validateDescription()

			if tc.wantErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}

		})
	}

}
func TestValidateTimeRange(t *testing.T) {
	tests := []struct {
		name    string
		from    time.Time
		to      time.Time
		wantErr error
	}{
		{
			name:    "valid range",
			from:    time.Date(2026, 12, 01, 10, 0, 0, 0, time.UTC),
			to:      time.Date(2026, 12, 01, 11, 0, 0, 0, time.UTC),
			wantErr: nil,
		},
		{
			name:    "invalid range",
			from:    time.Date(2026, 12, 01, 10, 0, 0, 0, time.UTC),
			to:      time.Date(2026, 12, 01, 9, 0, 0, 0, time.UTC),
			wantErr: errs.ErrInvalidTimeRange,
		},
		{
			name:    "equal range",
			from:    time.Date(2026, 12, 01, 10, 0, 0, 0, time.UTC),
			to:      time.Date(2026, 12, 01, 10, 0, 0, 0, time.UTC),
			wantErr: errs.ErrInvalidTimeRange,
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {
			booking := Booking{
				FromTime: tc.from,
				ToTime:   tc.to,
			}

			err := booking.validateTimeRange()

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidatMinMeetingTime(t *testing.T) {
	tests := []struct {
		name    string
		from    time.Time
		to      time.Time
		wantErr error
	}{
		{
			name:    "valid min time",
			from:    time.Date(2026, 12, 01, 10, 0, 0, 0, time.UTC),
			to:      time.Date(2026, 12, 01, 11, 0, 0, 0, time.UTC),
			wantErr: nil,
		},
		{
			name:    "invalid min time",
			from:    time.Date(2026, 12, 01, 10, 0, 0, 0, time.UTC),
			to:      time.Date(2026, 12, 01, 10, 27, 0, 0, time.UTC),
			wantErr: errs.ErrMinTimeLack,
		},
		{
			name:    "exact min time",
			from:    time.Date(2026, 12, 01, 10, 0, 0, 0, time.UTC),
			to:      time.Date(2026, 12, 01, 10, 30, 0, 0, time.UTC),
			wantErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			booking := Booking{
				FromTime: tc.from,
				ToTime:   tc.to,
			}

			err := booking.validateMinMeetingTime()

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}

		})
	}
}

func TestValidateMaxMeetingTime(t *testing.T) {
	tests := []struct {
		name    string
		from    time.Time
		to      time.Time
		wantErr error
	}{
		{
			name:    "valid max time",
			from:    time.Date(2026, 12, 01, 9, 0, 0, 0, time.UTC),
			to:      time.Date(2026, 12, 01, 16, 0, 0, 0, time.UTC),
			wantErr: nil,
		},
		{
			name:    "invalid max time",
			from:    time.Date(2026, 12, 01, 9, 0, 0, 0, time.UTC),
			to:      time.Date(2026, 12, 01, 18, 0, 0, 0, time.UTC),
			wantErr: errs.ErrMaxTimeExceed,
		},
		{
			name:    "exact time",
			from:    time.Date(2026, 12, 01, 9, 0, 0, 0, time.UTC),
			to:      time.Date(2026, 12, 01, 17, 0, 0, 0, time.UTC),
			wantErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			booking := Booking{
				FromTime: tc.from,
				ToTime:   tc.to,
			}

			err := booking.validateMaxMeetingTime()

			if tc.wantErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateBackdate(t *testing.T) {
	tests := []struct {
		name    string
		from    time.Time
		wantErr error
	}{
		{
			name:    "valid date",
			from:    time.Date(2026, 12, 01, 10, 0, 0, 0, time.UTC),
			wantErr: nil,
		},
		{
			name:    "back date",
			from:    time.Date(2026, 11, 11, 10, 0, 0, 0, time.UTC),
			wantErr: errs.ErrInvalidBookingTime,
		},
		{
			name:    "valid date + 1 minute",
			from:    time.Now().Add(60 * time.Second),
			wantErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			booking := Booking{
				FromTime: tc.from,
			}
			err := booking.validateBackdate()

			if tc.wantErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.wantErr)
			} else {
				fmt.Println("Booking: ", booking.FromTime, " tc.from: ", tc.from)
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateParticipantsCount(t *testing.T) {
	tests := []struct {
		name     string
		actualPC int
		wantErr  error
	}{
		{
			name:     "valid participants count",
			actualPC: 10,
			wantErr:  nil,
		},
		{
			name:     "zero participants count",
			actualPC: 0,
			wantErr:  errs.ErrInvalidParticipantsCount,
		},
		{
			name:     "negative participants count",
			actualPC: -2,
			wantErr:  errs.ErrInvalidParticipantsCount,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			booking := Booking{
				ParticipantsCount: tc.actualPC,
			}

			err := booking.validateParticipantsNumber()

			if tc.wantErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateAmount(t *testing.T) {
	tests := []struct {
		name         string
		actualAmount int64
		wantError    error
	}{
		{
			name:         "valid amount",
			actualAmount: 200,
			wantError:    nil,
		},
		{
			name:         "zero mount",
			actualAmount: 0,
			wantError:    nil,
		},
		{
			name:         "negative amount",
			actualAmount: -200,
			wantError:    errs.ErrInvalidAmount,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			booking := Booking{
				Amount: tc.actualAmount,
			}

			err := booking.validateAmount()

			if tc.wantError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.wantError)
			} else {
				require.NoError(t, err)
			}

		})
	}
}
