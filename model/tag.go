package model

import (
	"github.com/allape/gocrud"
)

type Tag struct {
	gocrud.Base
	Name        string `json:"name"`
	Alias       string `json:"alias"`
	Priority    int64  `json:"priority"`
	Color       string `json:"color"`
	Description string `json:"description"`
}
