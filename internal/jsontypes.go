package api

import (
	"time"
)

type UserJSONRecord struct {
	Username string `json:"username"`
	Subscriptions []string `json:"subscriptions"`
}

type ItemJSONRecord struct {
	Slug string `json:"slug"`
	Title string `json:"title"`
	DateUploaded time.Time `json:"date_uploaded"`
	Uploader string `json:"uploader"`
}

type ChannelJSONRecord struct {
	Slug string `json:"slug"`
	Title string `json:"title"`
	Items []*ItemJSONRecord `json:"items"`
}
