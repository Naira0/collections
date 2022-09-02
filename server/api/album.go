package api

import (
	"collections-server/database"
	"encoding/json"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"github.com/rs/xid"
)

// routes concerned with albums

func UploadAlbum(ctx *fiber.Ctx) error {

	session := verifySession(ctx)

	if session == nil {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	form, err := ctx.MultipartForm()

	if err != nil {
		return err
	}

	body := new(UploadAlbumBody)
	data, exists := form.Value["data"]

	if !exists {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	if err := json.Unmarshal([]byte(data[0]), body); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	if err := validate.Struct(body); err != nil {
		return err
	}

	delete(form.Value, "data")

	n := len(form.Value)

	if n == 0 {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	id := xid.New().String()
	dir := "./files/" + id + "/"

	os.Mkdir(dir, 0644)

	files := make([]string, n)
	var i int

	for name, bytes := range form.Value {
		err := os.WriteFile(dir+name, []byte(bytes[0]), 0644)

		if err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		files[i] = name
		i++
	}

	_, err = db.Exec("INSERT INTO albums(id, name, authorId, description, tags, files, createdAt) VALUES($1, $2, $3, $4, $5, $6, (SELECT CURRENT_TIMESTAMP))",
		id, body.Name, session.UserId, body.Description, pq.Array(body.Tags), pq.Array(files))

	if err != nil {
		os.RemoveAll(dir)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return nil
}

func DeleteAlbum(ctx *fiber.Ctx) error {

	session := verifySession(ctx)

	if session == nil {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	id := ctx.Params("id")

	_, err := db.Exec("DELETE FROM albums WHERE id = $1", id)

	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return os.RemoveAll("./files/" + id)
}

// NOTE everything having to do with bookmarks is kinda slow because of the dumb way i modeled the data but i cant think of a better way

// this is currently less than efficient, optomize eventually.
func AllBookmarks(ctx *fiber.Ctx) error {
	session := verifySession(ctx)

	if session == nil {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	offset, err := strconv.Atoi(ctx.Query("offset", "0"))

	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	const LIMIT = 12

	var album_ids pq.StringArray

	err = db.Get(&album_ids, "SELECT bookmarks FROM users WHERE id = $1 LIMIT $2 OFFSET $3",
		session.UserId, LIMIT, offset)

	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	albums := make([]AlbumMetaData, len(album_ids))

	for i := 0; i < len(album_ids); i++ {
		err = db.Get(&albums[i], "SELECT name, description, authorId, id, likes from albums where id = $1", album_ids[i])

		if err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}
	}

	return ctx.JSON(fiber.Map{
		"data": albums,
	})
}

func UpdateBookmark(ctx *fiber.Ctx) error {

	session := verifySession(ctx)

	if session == nil {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	id := ctx.Params("id")
	op := ctx.Query("op")

	if !database.CheckRow(db, "albums", "id", id) {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	var err error

	switch op {
	case "set":
		var bookmarks pq.StringArray
		err = db.Get(&bookmarks, "SELECT bookmarks FROM users WHERE id = $1", session.UserId)

		if err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		if searchArray(bookmarks, id) {
			return ctx.SendStatus(fiber.StatusConflict)
		}
		_, err = db.Exec("UPDATE users SET bookmarks = array_append(bookmarks, $1) where id = $2", id, session.UserId)

	case "del":
		_, err = db.Exec("UPDATE users SET bookmarks = array_remove(bookmarks, $1) where id = $2", id, session.UserId)
	default:
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	return err
}

//func PostComment(ctx *fiber.Ctx) error {
//	body, session, err := validateRequest[PostCommentBody](ctx)
//
//	if err != nil {
//		return err
//	}
//
//	if !database.CheckRow(db, "albums", "id", body.AlbumId) {
//		return ctx.SendStatus(fiber.StatusNotFound)
//	}
//
//	_, err = db.Exec("INSERT INTO comments(album_id, author, contents, time) VALUES($1, $2, $3, (SELECT CURRENT_TIMESTAMP))",
//		body.AlbumId, session.Identifier)
//
//	return nil
//}
