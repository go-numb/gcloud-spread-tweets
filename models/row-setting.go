package models

// Set is a struct that represents the setting of a row.
type Set struct {
	UUID string `csv:"-" dataframe:"uuid" firestore:"uuid,omitempty" json:"uuid,omitempty"`

	// 同時刻のスケジュール予約がある場合、全ての投稿を有効にするかどうか
	// why: 複数投稿を同時刻に行う場合、全ての投稿を有効にするかどうかを設定する。年・月・週指定の同時刻投稿を行いたい。
	IsScheduleAllPost bool `csv:"-" dataframe:"is_schedule_all_post" firestore:"is_schedule_all_post,omitempty" json:"is_schedule_all_post,omitempty"`
}
