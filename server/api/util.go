package api

import (
	"collections-server/database"
	"database/sql"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/xid"
)

type Execable interface {
	Exec(query string, args ...any) (sql.Result, error)
}

func searchArray[T comparable](array []T, element T) bool {
	for _, e := range array {
		if e == element {
			return true
		}
	}

	return false
}

func createSession(e Execable, user_id string) (string, error) {
	session_id := xid.New().String()

	_, err := e.Exec("INSERT INTO sessions(id, userId) VALUES($1, $2)",
		session_id, user_id)

	return session_id, err
}

func verifySession(ctx *fiber.Ctx) *database.Session {
	id := ctx.Cookies("id")

	if len(id) == 0 {
		return nil
	}

	session := new(database.Session)
	err := db.Get(session, "SELECT * FROM sessions WHERE id = $1", id)

	if err != nil {
		return nil
	}

	return session
}

func validateBody[T any](ctx *fiber.Ctx) (*T, error) {
	body := new(T)

	if err := ctx.BodyParser(body); err != nil {
		return nil, err
	}

	if err := validate.Struct(body); err != nil {
		return nil, err
	}

	return body, nil
}

func validateRequest[T any](ctx *fiber.Ctx) (*T, *database.Session, error) {
	session := verifySession(ctx)

	if session == nil {
		return nil, nil, fiber.NewError(fiber.StatusUnauthorized, "no valid session id found")
	}

	body, err := validateBody[T](ctx)

	return body, session, err
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// eraseAlbum removes the album from the db and filesystem
func eraseAlbum(e Execable, condition, id string) error {

	err := os.RemoveAll("./files/" + id)

	if err != nil {
		return err
	}

	_, err = e.Exec("DELETE FROM albums WHERE "+condition+"= $1", id)

	return err
}
