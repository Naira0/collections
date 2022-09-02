package api

type CreateAccountBody struct {
	Username string `validate:"required,max=40"`
	Password string `validate:"required,min=3,max=40"`
	Email    string `validate:"required,email"`
}

type LoginBody struct {
	// username or email
	Identifier string `validate:"required"`
	Password   string `validate:"required,min=3,max=40"`
}

type UpdateAccountBody struct {
	Username string `validate:"max=40"`
	Bio      string `validate:"max=255"`
	Email    string `validate:"omitempty,email"`
}

type ChangePasswordBody struct {
	OldPassword string `validate:"required,min=3,max=40"`
	NewPassword string `validate:"required,min=3,max=40"`
}

type UploadAlbumBody struct {
	Name        string `validate:"required,min=1,max=40"`
	Description string `validate:"max=255"`
	Tags        []string
}

type AlbumMetaData struct {
	Name,
	Description,
	AuthorId,
	Id string
	Likes uint32
}

type UserMetaData struct {
	Username,
	Bio string
	Profile_pic []byte
}

type PostCommentBody struct {
	Comment string `validate:"required,min=1,max=255"`
	AlbumId string `validate:"required"`
}
