package database

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
	Id,
	Name,
	AuthorId,
	Description,
	CreatedAt string
	Likes uint32
	Tags,
	Files []string
}

type Session struct {
	Id,
	UserId string
}
