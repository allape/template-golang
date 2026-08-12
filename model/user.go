package model

import (
	"time"

	"github.com/allape/gocrud"
	"github.com/allape/gophorward"
)

type ShareLink struct {
	gocrud.Base
	NonceID     string     `json:"nonceId" gorm:"index"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	ExpireAt    time.Time  `json:"expireAt"`
	ConsumedAt  *time.Time `json:"consumedAt"`
}

type ShareLinkGallery struct {
	ShareID   gocrud.ID `json:"shareId" gorm:"primaryKey"`
	GalleryID gocrud.ID `json:"galleryId" gorm:"primaryKey"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime;<-:create"`
}

type UserGallery struct {
	UserID    gophorward.UserID `json:"userId" gorm:"primaryKey"`
	GalleryID gocrud.ID         `json:"galleryId" gorm:"primaryKey"`
	CreatedAt time.Time         `json:"createdAt" gorm:"autoCreateTime;<-:create"`
}
