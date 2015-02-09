package api

import (
	"time"
)

type UserDBRecord struct {
	Username string "_id,omitempty"
	Subscriptions []string "subscriptions"
}

type ItemDBRecord struct {
	Slug string "_id,omitempty"
	Title string "title"
	DateUploaded time.Time "date_uploaded"
	Data string "data"
	Uploader string "uploader"
}

type ChannelDBRecord struct {
	Slug string "_id,omitempty"
	Title string "title"
	Items []ItemDBRecord "items"
}
