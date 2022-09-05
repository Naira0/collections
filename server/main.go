package main

import (
	"collections-server/api"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"strings"
)

/*
TODO
	- create thumbnails for albums
	- notification system
	- implement comments
	- implement likes
	- implement search
*/

type Config struct {
	Port string
}

func newConfig() *Config {
	bytes, err := os.ReadFile("./config.json")

	default_config := func() *Config {
		return &Config{
			Port: ":3000",
		}
	}

	if err != nil {
		return default_config()
	}

	config := new(Config)

	err = json.Unmarshal(bytes, config)

	if err != nil {
		return default_config()
	}

	if !strings.HasPrefix(config.Port, ":") {
		config.Port = ":" + config.Port
	}

	return config
}

func main() {

	config := newConfig()

	os.Mkdir("./files/profiles/", 0750)

	app := fiber.New()
	user := app.Group("/user")

	user.Get("/get/", api.GetAccount)
	user.Post("/create", api.CreateAccount)
	user.Post("/login", api.Login)
	user.Post("/logout", api.Logout)
	user.Delete("/delete", api.DeleteAccount)
	user.Patch("/update", api.UpdateAccount)
	user.Patch("/change_password", api.ChangePassword)
	user.Patch("/update_bookmark/:id", api.UpdateBookmark)
	user.Get("/all_bookmarks", api.AllBookmarks)
	user.Patch("/set_profile_pic/", api.SetProfilePic)

	app.Static("/files", "./files/")

	album := app.Group("/album")

	album.Post("/upload", api.UploadAlbum)
	album.Delete("/delete/:id", api.DeleteAlbum)
	album.Get("/get/:id", api.GetAlbum)

	app.Post("/test", api.Test)

	log.Fatal(app.Listen(config.Port))
}
