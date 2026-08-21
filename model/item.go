package model

import (
	"fmt"
	"image"
	"mime"
	"path"
	"strings"
	"time"

	"github.com/allape/gocrud"
	"github.com/allape/golang/helper"
)

type Item struct {
	gocrud.Base
	Name        string `json:"name"`
	Src         string `json:"src"`
	Thumbnail   string `json:"thumbnail"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

type ItemTag struct {
	ItemID    gocrud.ID `json:"itemId" gorm:"primaryKey"`
	TagID     gocrud.ID `json:"tagId" gorm:"primaryKey"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime;<-:create"`
}

func NewMediaPreview(dst, src string) error {
	mimeStr := mime.TypeByExtension(path.Ext(dst))

	switch {
	case strings.HasPrefix(mimeStr, "image/"):
		ext := strings.ToLower(path.Ext(dst))
		switch ext {
		case ".gif":
			_, err := helper.FFMpegVideoSampleImage(dst, src, 0.5, image.Point{X: 2, Y: 2})
			if err != nil {
				return err
			}
		case ".raw":
			fallthrough
		case ".arw":
			err := helper.ExifToolPreview(dst, src)
			if err != nil {
				return err
			}
		default:
			_, err := helper.FFMpegScale(dst, src, 0.2)
			if err != nil {
				return err
			}
		}
	case strings.HasPrefix(mimeStr, "video/"):
		_, err := helper.FFMpegVideoSampleImage(dst, src, 0.2, image.Point{X: 10, Y: 10})
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("filetype %s is not supported", mimeStr)
	}

	return nil
}

func NewMedia(dst, src string) error {
	ext := strings.ToLower(path.Ext(dst))
	switch ext {
	case ".raw":
		fallthrough
	case ".arw":
		err := helper.ExifToolPreview(dst, src)
		if err != nil {
			return err
		}
	default:
		_, err := helper.FFMpegScale(dst, src, 0.5)
		if err != nil {
			return err
		}
	}

	return nil
}
