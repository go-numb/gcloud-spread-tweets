package models

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/go-gota/gota/dataframe"
)

type Post struct {
	UUID string `csv:"uuid" dataframe:"uuid" firestore:"uuid,omitempty" json:"uuid,omitempty"`
	// ID is Twitter/X AccountID
	ID        string `csv:"id" dataframe:"id" firestore:"id,omitempty" json:"id,omitempty"`
	Text      string `csv:"text" dataframe:"text" firestore:"text,omitempty" json:"text,omitempty"`
	File1     string `csv:"file1" dataframe:"file1" firestore:"file_1,omitempty" json:"file_1,omitempty"`
	File2     string `csv:"file2" dataframe:"file2" firestore:"file_2,omitempty" json:"file_2,omitempty"`
	File3     string `csv:"file3" dataframe:"file3" firestore:"file_3,omitempty" json:"file_3,omitempty"`
	File4     string `csv:"file4" dataframe:"file4" firestore:"file_4,omitempty" json:"file_4,omitempty"`
	WithFiles int    `csv:"with_files" dataframe:"with_files" firestore:"with_files,omitempty" json:"with_files,omitempty"`
	Checked   int    `csv:"checked" dataframe:"checked" firestore:"checked,omitempty" json:"checked,omitempty"`
	Priority  int    `csv:"priority" dataframe:"priority" firestore:"priority,omitempty" json:"priority,omitempty"`
	Count     int    `csv:"count" dataframe:"count" firestore:"count,omitempty" json:"count,omitempty"`
	PostURL   string `csv:"post_url" dataframe:"post_url" firestore:"post_url,omitempty" json:"post_url,omitempty"`

	// 以下は、csv, dataframeには含まれない
	IsDelete     bool      `csv:"-" dataframe:"-" firestore:"is_delete,omitempty" json:"-,omitempty"`
	LastPostedAt time.Time `csv:"-" dataframe:"-" firestore:"last_posted_at,omitempty" json:"last_posted_at,omitempty"`
	CreatedAt    time.Time `csv:"-" dataframe:"-" firestore:"created_at,omitempty" json:"created_at,omitempty"`
}

func (p Post) GetID() string {
	return p.ID
}

func (p *Post) SetLastPostedAt() bool {
	if p.LastPostedAt.IsZero() {
		p.LastPostedAt = time.Now()
		return true
	}
	return false
}

// IsLastPostedAt LastPostedAtの時間を比較
// More than the specified time has elapsed.
func (p *Post) IsPastLastPostedAt(minutes int) bool {
	return p.LastPostedAt.Before(time.Now().Add(-time.Duration(minutes) * time.Minute))
}

func (p *Post) SetCreateAt() bool {
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
		return true
	}
	return false
}

// ToURLValues PostID, AccountIDをURLValuesに変換
// この値を使って、Firestore[Accounts,Posts]からデータを取得する
func (p *Post) ToURLValues() url.Values {
	v := url.Values{}
	v.Set("account_id", p.ID)
	v.Set("post_id", p.UUID)
	return v

}

// CheckDupID SpreadID, UserIDの重複を確認
func CheckDupID(id, spreadID string, accounts []Account) error {
	for _, row := range accounts {
		// SpreadIDが重複していないかを確認
		if row.SpreadID == spreadID {
			return errors.New("spread id is already exist")
		}

		// UserIDが重複していないかを確認
		if row.ID == id {
			return errors.New("x_id is already exist")
		}
	}

	return nil
}

func CheckColumns(df dataframe.DataFrame) error {
	str := make([]Post, 1)
	columns := dataframe.LoadStructs(str)
	if columns.Err != nil {
		return columns.Err
	}

	baseColumnNames := columns.Names()
	customerColumnNames := df.Names()

	// カラム名とカラム数が合致するかを確認
	if len(baseColumnNames) != len(customerColumnNames) {
		return fmt.Errorf("columns name and customer columns are not match, base: %v, customer: %v", baseColumnNames, customerColumnNames)
	}
	// base, customerのカラム数は一致している
	// 上で、名称が一致するかを確認
	for i, name := range baseColumnNames {
		if name != customerColumnNames[i] {
			return fmt.Errorf("columns name and customer columns are not match, base: %v, customer: %v", baseColumnNames, customerColumnNames)
		}
	}

	return nil
}

// DfNrowToLastNrow Nrowを末尾に追加するためのNrowを返す
func DfNrowToLastNrow(df dataframe.DataFrame) int {
	// Nrowはheaderを含まないため、+1を追加
	// さらに、空の末尾に追加するため、+1
	return df.Nrow() + 1 + 1
}
