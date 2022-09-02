package api

// routes related to a user account

import (
	"bytes"
	"collections-server/database"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

func Test(ctx *fiber.Ctx) error {
	_, err := os.ReadFile("cc4eg0dfacp3oq1v1gkg")
	return err
}

// TODO make it check if email already exists
func CreateAccount(ctx *fiber.Ctx) error {

	body, err := validateBody[CreateAccountBody](ctx)

	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	tx, _ := db.Begin()

	salt := xid.New().String()

	hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password+salt), bcrypt.DefaultCost)

	user_id := xid.New().String()

	_, err = tx.Exec("INSERT INTO users(id, username, email, bio, salt, password) VALUES($1, $2, $3, $4, $5, $6)",
		user_id, body.Username, body.Email, "", salt, hash)

	if err != nil {
		return ctx.SendStatus(fiber.StatusConflict)
	}

	session_id, _ := createSession(tx, user_id)

	err = tx.Commit()

	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	ctx.Cookie(&fiber.Cookie{
		Name:  "id",
		Value: session_id,
	})

	return nil
}

func DeleteAccount(ctx *fiber.Ctx) error {

	session := verifySession(ctx)

	if session == nil {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	tx, _ := db.Begin()

	tx.Exec("DELETE FROM users WHERE id = $1", session.UserId)
	tx.Exec("DELETE FROM sessions WHERE userid = $1", session.UserId)

	var album_ids pq.StringArray
	db.Select(&album_ids, "SELECT id FROM albums WHERE authorId = $1", session.UserId)

	tx.Exec("DELETE FROM albums WHERE authorId = $1", session.UserId)

	err := eraseAlbum(tx, "authorId", session.UserId)

	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	ctx.Cookie(&fiber.Cookie{
		Name:  "id",
		Value: "invalid",
	})

	if tx.Commit() != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	for _, id := range album_ids {
		os.RemoveAll("./files/" + id)
	}

	return nil
}

// NOTE currently allows multiple sessions, be aware of any potential bugs or vulnerabilites that could come from this
func Login(ctx *fiber.Ctx) error {

	body, err := validateBody[LoginBody](ctx)

	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	user := new(database.User)
	err = db.Get(user, "SELECT id, username, salt, password FROM users WHERE username = $1 or email = $1", body.Identifier)

	if err != nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	id_match := body.Identifier == user.Email || body.Identifier == user.Username

	pass_match := bcrypt.CompareHashAndPassword(user.Password, []byte(body.Password+user.Salt))

	if !id_match || pass_match != nil {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	session_id, err := createSession(db, user.Id)

	if err != nil {
		return err
	}

	ctx.Cookie(&fiber.Cookie{
		Name:  "id",
		Value: session_id,
	})

	return nil
}

func Logout(ctx *fiber.Ctx) error {
	session := verifySession(ctx)

	if session == nil {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	_, err := db.Exec("DELETE FROM sessions WHERE id = $1", session.Id)

	ctx.Cookie(&fiber.Cookie{
		Name:  "id",
		Value: "invalid",
	})

	return err
}

func UpdateAccount(ctx *fiber.Ctx) error {

	body, session, err := validateRequest[UpdateAccountBody](ctx)

	if err != nil {
		return err
	}

	tx, _ := db.Begin()

	var query string

	value := reflect.ValueOf(*body)
	fields := reflect.VisibleFields(value.Type())

	var query_values [4]any
	var count int

	// argument to Exec must be a variadic so expanding the array plus this as an additional argument would be invalid
	query_values[0] = session.UserId

	for i, f := range fields {

		switch f.Index[0] {
		case 0:
			if len(body.Username) == 0 {
				continue
			}
		case 1:
			if len(body.Bio) == 0 {
				continue
			}
		case 2:
			if len(body.Email) == 0 {
				continue
			}
		}

		count++

		n := strconv.Itoa(count + 1)
		name := f.Name

		query += strings.ToLower(name[0:1]) + name[1:] + " = $" + n + ", "

		query_values[count] = value.Field(i).Interface()
	}

	if count == 0 {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	query = query[0 : len(query)-2]

	tx.Exec("UPDATE users SET "+query+" WHERE id = $1", query_values[:count+1]...)

	return tx.Commit()
}

func SetProfilePic(ctx *fiber.Ctx) error {

	session := verifySession(ctx)

	if session == nil {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	data := ctx.Body()

	// fails if body is larger than 5mbs
	if unsafe.Sizeof(data) >= 5e+6 {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	// makes sure the image is a valid png or jpg
	_, _, err := image.Decode(bytes.NewReader(data))

	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	err = os.WriteFile("./files/profiles/"+session.UserId, data, 0644)

	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return nil
}

func ChangePassword(ctx *fiber.Ctx) error {

	body, session, err := validateRequest[ChangePasswordBody](ctx)

	if err != nil {
		return err
	}

	user := new(database.User)
	err = db.Get(user, "SELECT id, password, salt FROM users WHERE id = $1", session.UserId)

	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(body.OldPassword+user.Salt))

	if err != nil {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	salt := xid.New().String()
	hash, _ := bcrypt.GenerateFromPassword([]byte(body.NewPassword+salt), bcrypt.DefaultCost)

	_, err = db.Exec("UPDATE users SET password = $1, salt = $2 WHERE id = $3", hash, salt, user.Id)

	return err
}

/*
GetAccount retrieves metadata on a user with a given user id or username
expects plain string

TODO send thumbnails
*/
func GetAccount(ctx *fiber.Ctx) error {

	identifier := string(ctx.Body())
	user := new(database.User)

	err := db.Get(user, "SELECT id, username, bio FROM users WHERE username = $1 or id = $1", identifier)

	if err != nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	albums := new([]AlbumMetaData)
	err = db.Select(albums, "SELECT name, description, authorID, id, likes from albums where authorId = $1", user.Id)

	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"username": user.Username,
		"id":       user.Id,
		"bio":      user.Bio,
		"albums":   albums,
	})
}
