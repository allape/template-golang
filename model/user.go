package model

import (
	"github.com/allape/gocrud"
)

type User struct {
	gocrud.Base
	Name        string `json:"name"`
	Description string `json:"description"`
}
