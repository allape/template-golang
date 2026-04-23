package model

import (
	"time"

	"github.com/allape/gocrud"
	"github.com/allape/gophorward"
)

const MaxItemNameLength = 50

type Item struct {
	gocrud.Base
	Name        string            `json:"name"`
	Src         string            `json:"src"`
	Priority    int64             `json:"priority"`
	Description string            `json:"description"`
	CreatedBy   gophorward.UserID `json:"createdBy"`
}

type ItemTag struct {
	ItemID    gocrud.ID `json:"itemId" gorm:"primaryKey"`
	TagID     gocrud.ID `json:"tagId" gorm:"primaryKey"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime;<-:create"`
}
