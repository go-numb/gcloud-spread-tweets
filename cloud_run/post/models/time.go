package models

import "time"

// IsTimeToPost は投稿する時間かどうかを判定する
func IsTimeToPost(t time.Time) bool {
	// 5分ごとに投稿
	return t.Minute()%5 == 0
}
