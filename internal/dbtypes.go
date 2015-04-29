package api

import (
	"time"
)

type UserDBRecord struct {
	Username      string   "_id,omitempty"
	Subscriptions []string "subscriptions"
}

func (self *UserDBRecord) ToJSON() *UserJSONRecord {
	return &UserJSONRecord{
		Username:      self.Username,
		Subscriptions: self.Subscriptions,
	}
}

type ItemDBRecord struct {
	Slug         string    "_id,omitempty"
	Title        string    "title"
	DateUploaded time.Time "date_uploaded"
	Data         string    "data"
	Uploader     string    "uploader"
}

func (self *ItemDBRecord) ToJSON() *ItemJSONRecord {
	return &ItemJSONRecord{
		Slug:         self.Slug,
		Title:        self.Title,
		DateUploaded: self.DateUploaded,
		Uploader:     self.Uploader,
	}
}

type ChannelDBRecord struct {
	Slug  string         "_id,omitempty"
	Title string         "title"
	Owner string         "owner"
	Items []ItemDBRecord "items"
}

func (self *ChannelDBRecord) ToJSON() *ChannelJSONRecord {
	items := make([]*ItemJSONRecord, 0)
	for _, item := range self.Items {
		items = append(items, item.ToJSON())
	}
	return &ChannelJSONRecord{
		Slug:  self.Slug,
		Title: self.Title,
		Items: items,
	}
}
