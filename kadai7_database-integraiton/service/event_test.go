package service

import (
	"fmt"
	"kadai7_database-integration/mocks"
	"kadai7_database-integration/repository"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestEventService_GetAllEvents(t *testing.T) {
	testCases := []struct {
		name           string
		mockEvents     []repository.Event
		mockError      error
		expectedEvents []repository.Event
		expectedError  error
		wantErr        bool
	}{
		// テストケース1: 正常系
		{
			name: "正常系",
			mockEvents: []repository.Event{
				{ID: 1, Title: "テストイベント1", Date: "2025-01-01"},
			},
			expectedEvents: []repository.Event{
				{ID: 1, Title: "テストイベント1", Date: "2025-01-01"},
			},
			expectedError: nil,
			wantErr:       false,
		},
		// テストケース2: 異常系
		{
			name:           "異常系_データベースエラー",
			mockEvents:     nil,
			mockError:      fmt.Errorf("データベース接続エラー"),
			expectedEvents: nil,
			expectedError:  fmt.Errorf("データベース接続エラー"),
			wantErr:        true,
		},
		// テストケース3: 境界値
		{
			name:           "データが空",
			mockEvents:     []repository.Event{},
			mockError:      nil,
			expectedEvents: []repository.Event{},
			expectedError:  nil,
			wantErr:        false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// mockの準備
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockRepo := mocks.NewMockEventRepositoryInterface(mockCtrl)

			// mockの動作を設定、mockイベントを返すように設定
			mockRepo.EXPECT().GetAll().Return(tc.mockEvents, tc.mockError)

			// Service層をテスト
			service := NewEventService(mockRepo)
			actualEvents, err := service.GetAllEvents()

			// エラーが発生してないかチェック、期待と異なるエラー状況の場合
			if (err != nil && !tc.wantErr) || (err == nil && tc.wantErr) {
				t.Errorf("error = %v, wantErr = %v", err, tc.wantErr)
			}
			// エラーメッセージチェック
			if tc.wantErr && err != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("error message = %v, want = %v", err.Error(), tc.expectedError)
			}
			// データ件数チェック
			if len(actualEvents) != len(tc.expectedEvents) {
				t.Errorf("expected event count = %d, got = %d", len(tc.expectedEvents), len(actualEvents))
			}
			// データ内容チェック
			for i, expected := range tc.expectedEvents {
				if i >= len(actualEvents) {
					t.Errorf("expected event count = %d, got = %d", len(tc.expectedEvents), len(actualEvents))
					continue
				}
				actual := actualEvents[i]
				if actual.ID != expected.ID || actual.Title != expected.Title || actual.Date != expected.Date {
					t.Errorf("acutualEvents[%d] = %+v, want = %+v", i, actual, expected)
				}
			}

		})
	}
}
