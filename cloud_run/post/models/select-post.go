package models

import (
	"fmt"
	"sort"
)

func SelectPost(posts []Post) (*Post, error) {
	var results []Post
	for i := 0; i < len(posts); i++ {
		if posts[i].ID == "" {
			continue
		}

		if posts[i].AccessToken == "" {
			continue
		}

		if posts[i].AccessSecret == "" {
			continue
		}

		if posts[i].Text == "" {
			continue
		}

		if posts[i].Checked == 0 {
			continue
		}

		results = append(results, posts[i])
	}

	results = CustomSelect(results)

	if len(results) == 0 {
		return nil, fmt.Errorf("posts is empty")
	}

	return &results[one(len(results))], nil
}

func CustomSelect(posts []Post) []Post {
	// Priority sort 降順 高いものが先
	// 選択数が上位10件
	posts = sortPriority(posts)
	posts = _selectSome(10, posts)

	// Count sort 昇順 低いものが先
	// 選択数が上位1件
	posts = sortCount(posts)
	posts = _selectSome(1, posts)

	return posts
}

func sortPriority(posts []Post) []Post {
	// Priority sort 降順 高いものが先
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Priority > posts[j].Priority
	})

	return posts
}

func sortCount(posts []Post) []Post {
	// Priority sort 昇順 低いものが先
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Count < posts[j].Count
	})

	return posts
}

// _selectSome はスライスの先頭から n 個を返す
// n=1であれば先頭1個を返す
func _selectSome(n int, posts []Post) []Post {
	if n > len(posts) {
		return posts
	}

	return posts[:n]
}
