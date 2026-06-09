package model

import (
	"github.com/allape/gocrud"
)

type Tag struct {
	gocrud.Base
	Name        string `json:"name"`
	Alias       string `json:"alias"`
	Color       string `json:"color"`
	Description string `json:"description"`
}
