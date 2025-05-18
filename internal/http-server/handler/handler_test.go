package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/sanchey92/metric-server/internal/http-server/handler/mocks"
	"github.com/sanchey92/metric-server/internal/models"
)

func TestHandler_HandleMetrics(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		setupMock      func(*mocks.MockMemStorage)
	}{
		{
			name: "valid_metrics",
			requestBody: []models.Metric{
				{Name: "cpu", Value: 42.5},
				{Name: "memory", Value: 75.0},
			},
			expectedStatus: http.StatusOK,
			setupMock: func(m *mocks.MockMemStorage) {
				m.EXPECT().Set("cpu", 42.5).Times(1)
				m.EXPECT().Set("memory", 75.0).Times(1)
			},
		},
		{
			name:           "invalid json",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(_ *mocks.MockMemStorage) {},
		},
		{
			name:           "empty array",
			requestBody:    []models.Metric{},
			expectedStatus: http.StatusOK,
			setupMock:      func(_ *mocks.MockMemStorage) {},
		},
		{
			name: "null value",
			requestBody: []models.Metric{
				{Name: "null_metric", Value: 0},
			},
			expectedStatus: http.StatusOK,
			setupMock: func(m *mocks.MockMemStorage) {
				m.EXPECT().Set("null_metric", 0.0).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := mocks.NewMockMemStorage(ctrl)
			tt.setupMock(mockStorage)

			handler := New(mockStorage)

			var body []byte
			var err error

			switch v := tt.requestBody.(type) {
			case string:
				body = []byte(v)
			default:
				body, err = json.Marshal(tt.requestBody)
				require.NoError(t, err, "Failed to marshal request body")
			}

			r := httptest.NewRequest(http.MethodPost, "/update", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.HandleMetrics(w, r)

			require.Equal(t, tt.expectedStatus, w.Code, "HTTP status should match expected")
		})
	}
}
