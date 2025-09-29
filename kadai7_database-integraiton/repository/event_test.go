package repository

import (
	"testing"
)

func TestEventRepository_ExistsByDateAndTitle(t *testing.T) {
	testCases := []struct {
		name     string
		date     string
		title    string
		expected bool
		wantErr  bool
	}{
		{
			name:     "存在するイベント",
			date:     "2025-09-18",
			title:    "ナカヨシファミリア〜きゅるりんってしてみて×ChumToto〜",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "存在しないイベント",
			date:     "2025-09-19",
			title:    "存在しないイベント",
			expected: false,
			wantErr:  false,
		},
		{
			name:     "空の日付",
			date:     "",
			title:    "テストイベント",
			expected: false,
			wantErr:  false,
		},
	}

	repo, err := NewEventRepository()
	if err != nil {
		t.Fatalf("Repository初期化エラー: %v", err)
	}
	defer repo.Close()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//テスト対象の関数を実行
			got, err := repo.ExistsByDateAndTitle(tc.date, tc.title)

			// エラーチェック: 期待と異なるエラー状況の場合
			if (err != nil && !tc.wantErr) || (err == nil && tc.wantErr) {
				t.Errorf("error = %v, wantErr = %v", err, tc.wantErr)
				return
			}

			// 結果チェック
			if got != tc.expected {
				t.Errorf("got = %v, expected = %v", got, tc.expected)
			}
		})
	}
}
