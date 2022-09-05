package database

import "github.com/lib/pq"

type Comment struct {
	AlbumID,
	AuthorId,
	Content,
	Time string
}

type User struct {
	Id,
	Username,
	Email,
	Bio,
	Salt string
	Bookmarks []string
	Password  []byte
}

type Album struct {
	Id          string         `json:"id"`
	Name        string         `json:"name"`
	AuthorId    string         `json:"authorId"`
	Description string         `json:"description"`
	CreatedAt   string         `json:"createdAt"`
	Likes       uint32         `json:"likes"`
	Tags        pq.StringArray `json:"tags"`
	Files       pq.StringArray `json:"files"`
}

type Session struct {
	Id,
	UserId string
}
