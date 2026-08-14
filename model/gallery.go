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
	Keywords    string            `json:"keywords"`
	IsPublic    bool              `json:"isPublic"`
	Description string            `json:"description"`
	CreatedBy   gophorward.UserID `json:"createdBy"`
	Enabled     bool              `json:"enabled"`
}

type GalleryItem struct {
	GalleryID gocrud.ID `json:"galleryId" gorm:"primaryKey"`
	ItemID    gocrud.ID `json:"itemId" gorm:"primaryKey"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime;<-:create"`
}

type GalleryTag struct {
	GalleryID gocrud.ID `json:"galleryId" gorm:"primaryKey"`
	TagID     gocrud.ID `json:"tagId" gorm:"primaryKey"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime;<-:create"`
}
