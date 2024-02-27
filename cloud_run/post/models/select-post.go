package models

import "fmt"

func SelectPost(posts []Post) (*Post, error) {

	if len(posts) == 0 {
		return nil, fmt.Errorf("posts is empty")
	}

	return &posts[0], nil
}
