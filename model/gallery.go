package model

import (
	"time"

	"github.com/allape/gocrud"
	"github.com/allape/gophorward"
)

const MaxGalleryNameLength = 50

type Gallery struct {
	gocrud.Base
	Name        string            `json:"name"`
	IsPublic    bool              `json:"isPublic"`
	Priority    int64             `json:"priority"`
	Description string            `json:"description"`
	CreatedBy   gophorward.UserID `json:"createdBy"`
}

type GalleryItem struct {
	GalleryID gocrud.ID `json:"galleryId" gorm:"primaryKey"`
	ItemID    gocrud.ID `json:"itemId" gorm:"primaryKey"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime;<-:create"`
}
