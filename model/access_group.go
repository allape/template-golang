package model

import (
	"time"

	"github.com/allape/gocrud"
)

type AccessGroup struct {
	gocrud.Base
	Name string `json:"name"`
}

type AccessGroupGallery struct {
	GroupID   gocrud.ID `json:"groupId" gorm:"primaryKey"`
	GalleryID gocrud.ID `json:"galleryId" gorm:"primaryKey"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime;<-:create"`
}

type AccessGroupUser struct {
	GroupID   gocrud.ID `json:"groupId" gorm:"primaryKey"`
	UserID    gocrud.ID `json:"userId" gorm:"primaryKey"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime;<-:create"`
}
