package api

import (
	"collections-server/database"

	"github.com/go-playground/validator"
)

var db = database.New()
var validate = validator.New()
