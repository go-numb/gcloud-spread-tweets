package models

import (
	"errors"
	"fmt"

	"github.com/go-gota/gota/dataframe"
)

type Post struct {
	UUID      string `csv:"uuid" dataframe:"uuid" json:"uuid,omitempty"`
	ID        string `csv:"x_id" dataframe:"x_id" json:"id,omitempty"`
	Text      string `csv:"x_text" dataframe:"x_text" json:"text,omitempty"`
	File1     string `csv:"x_file1" dataframe:"x_file1" json:"file_1,omitempty"`
	File2     string `csv:"x_file2" dataframe:"x_file2" json:"file_2,omitempty"`
	File3     string `csv:"x_file3" dataframe:"x_file3" json:"file_3,omitempty"`
	File4     string `csv:"x_file4" dataframe:"x_file4" json:"file_4,omitempty"`
	WithFiles int    `csv:"with_files" dataframe:"with_files" json:"with_files,omitempty"`
	Checked   int    `csv:"checked" dataframe:"checked" json:"checked,omitempty"`
	Priority  int    `csv:"priority" dataframe:"priority" json:"priority,omitempty"`
	Count     int    `csv:"count" dataframe:"count" json:"count,omitempty"`
	PostURL   string `csv:"post_url" dataframe:"post_url" json:"post_url,omitempty"`
	LastDate  string `csv:"last_date" dataframe:"last_date" json:"last_date,omitempty"`
}

func (p Post) GetID() string {
	return p.ID
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
