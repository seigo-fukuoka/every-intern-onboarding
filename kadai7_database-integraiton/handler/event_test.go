package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"kadai7_database-integration/mocks"
	"kadai7_database-integration/repository"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
)

func TestEventHandler_GetEvents(t *testing.T) {
	testCases := []struct {
		name           string
		mockEvents     []repository.Event
		mockError      error
		expectedStatus int
		expectedEvents []repository.Event
		wantErr        bool
	}{
		{
			name: "正常系_データが1件",
			mockEvents: []repository.Event{
				{ID: 1, Title: "テストイベント1", Date: "2025-01-01"},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedEvents: []repository.Event{
				{ID: 1, Title: "テストイベント1", Date: "2025-01-01"},
			},
			wantErr: false,
		},
		{
			name:           "異常系_データベースエラー",
			mockEvents:     nil,
			mockError:      fmt.Errorf("データベース接続エラー"),
			expectedStatus: http.StatusInternalServerError,
			expectedEvents: nil,
			wantErr:        true,
		},
		{
			name:           "境界値_空のデータ",
			mockEvents:     []repository.Event{},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedEvents: []repository.Event{},
			wantErr:        false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// ここからあなたが実装
			// 1. Service層のMockを作成
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockService := mocks.NewMockEventServiceInterface(mockCtrl)

			mockService.EXPECT().GetAllEvents().Return(tc.mockEvents, tc.mockError)
			// 2. Echo Contextを設定
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/events", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			// 3. Handler層のメソッドを実行
			handler := NewEventHandler(mockService)
			handler.GetEvents(c)
			// 4. HTTPステータスとJSONを検証
			if rec.Code != tc.expectedStatus {
				t.Errorf("HTTP status = %d, want = %d", rec.Code, tc.expectedStatus)
			}
			// JSON形式をチェック
			var actualEvents []repository.Event
			json.Unmarshal(rec.Body.Bytes(), &actualEvents)

			// データの件数チェック
			if len(actualEvents) != len(tc.expectedEvents) {
				t.Errorf("events count = %d, want = %d", len(actualEvents), len(tc.expectedEvents))
			}

			// データの内容チェック
			for i, expected := range tc.expectedEvents {
				if i >= len(actualEvents) {
					t.Errorf("actualEvents[%d] is missing", i)
					continue
				}
				actual := actualEvents[i]
				if actual.ID != expected.ID || actual.Title != expected.Title || actual.Date != expected.Date {
					t.Errorf("actualEvents[%d] = %+v, want = %+v", i, actual, expected)
				}
			}
		})
	}
}
