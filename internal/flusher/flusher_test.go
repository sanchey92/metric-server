package flusher

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/sanchey92/metric-server/internal/flusher/mocks"
)

func TestFlusher_Run(t *testing.T) {
	tests := []struct {
		name           string
		interval       time.Duration
		contextTimeout time.Duration
		expectedError  error
		setupMocks     func(*mocks.MockMemStorage, *mocks.MockPostgresStorage)
	}{
		{
			name:           "successful flush on interval",
			interval:       100 * time.Millisecond,
			contextTimeout: 250 * time.Millisecond,
			expectedError:  nil,
			setupMocks: func(mockMem *mocks.MockMemStorage, mockDB *mocks.MockPostgresStorage) {
				metrics := map[string]float64{
					"cpu":    43.5,
					"memory": 75.0,
				}
				mockMem.EXPECT().Snapshot().Return(metrics).MinTimes(1)
				mockDB.EXPECT().Save(gomock.Any(), metrics).MinTimes(1)
			},
		},
		{
			name:           "empty metrics snapshot",
			interval:       100 * time.Millisecond,
			contextTimeout: 250 * time.Millisecond,
			expectedError:  nil,
			setupMocks: func(mockMem *mocks.MockMemStorage, _ *mocks.MockPostgresStorage) {
				mockMem.EXPECT().Snapshot().Return(map[string]float64{}).MinTimes(1)
			},
		},
		{
			name:           "database error during flush",
			interval:       100 * time.Millisecond,
			contextTimeout: 250 * time.Millisecond,
			expectedError:  errors.New("failed to save metrics: database connection error"),
			setupMocks: func(mockMem *mocks.MockMemStorage, mockDB *mocks.MockPostgresStorage) {
				metrics := map[string]float64{"cpu": 42.5}
				mockMem.EXPECT().Snapshot().Return(metrics).MinTimes(1)
				dbErr := errors.New("database connection error")
				mockDB.EXPECT().Save(gomock.Any(), metrics).Return(dbErr).MinTimes(1)
				mockDB.EXPECT().Save(gomock.Any(), metrics).Return(dbErr).MaxTimes(1)
			},
		},
		{
			name:           "graceful shutdown with successful final flush",
			interval:       1 * time.Second,
			contextTimeout: 50 * time.Millisecond,
			expectedError:  nil,
			setupMocks: func(mockMem *mocks.MockMemStorage, mockDB *mocks.MockPostgresStorage) {
				metrics := map[string]float64{"cpu": 42.5}
				mockMem.EXPECT().Snapshot().Return(metrics).Times(1)
				mockDB.EXPECT().Save(gomock.Any(), metrics).Return(nil).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockMem := mocks.NewMockMemStorage(ctrl)
			mockDB := mocks.NewMockPostgresStorage(ctrl)

			tt.setupMocks(mockMem, mockDB)

			f := New(tt.interval, mockMem, mockDB)

			ctx, cancel := context.WithTimeout(context.Background(), tt.contextTimeout)
			defer cancel()

			err := f.Run(ctx)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
