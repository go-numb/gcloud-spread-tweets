package models

// GetUniqueExtractor は一意なキーを取得する関数型
type GetUniqueExtractor[T any] func(T) string

// GetUniqueKeys は汎用的な型 T のスライスに対応し、ユニークなキーのリストを返します。
// extractor パラメータは、各要素からキーを抽出する関数です。
func GetUniqueKeys[T any](data []T, extractor GetUniqueExtractor[T]) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range data {
		id := extractor(entry)
		if !keys[id] {
			keys[id] = true
			list = append(list, id)
		}
	}

	return list
}

// 配列の重複チェック
func CheckDuplicate(str []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range str {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
