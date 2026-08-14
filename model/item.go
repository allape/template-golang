package model

import (
	"fmt"
	"image"
	"io"
	"mime"
	"os"
	"path"
	"strings"
	"time"

	"github.com/allape/gocrud"
	"github.com/allape/golang/env"
	"github.com/allape/golang/helper"
	"gorm.io/gorm"
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

func SaveFile(db *gorm.DB, src io.ReadSeeker, size int64, ext string) (filenameSrc, filenameThu string, err error) {
	mimeStr := mime.TypeByExtension(ext)
	extSrc := ext
	extThu := ext

	switch {
	case strings.HasPrefix(mimeStr, "video/"):
		extSrc = ".mp4"
		extThu = ".jpg"
	case strings.HasPrefix(mimeStr, "image/"):
		fallthrough
	default:
		extSrc = ".jpg"
		extThu = ".jpg"
	}

	tmpOriginal, err := os.CreateTemp(os.TempDir(), "gocrud-ori-*"+ext)
	if err != nil {
		err = fmt.Errorf("ori: failed to create temp file")
		return
	}
	defer func() {
		_ = tmpOriginal.Close()
		_ = os.Remove(tmpOriginal.Name())
	}()

	tmpSrc, err := os.CreateTemp(os.TempDir(), "gocrud-src-*"+extSrc)
	if err != nil {
		err = fmt.Errorf("src: failed to create temp file")
		return
	}
	defer func() {
		_ = tmpSrc.Close()
		_ = os.Remove(tmpSrc.Name())
	}()

	tmpThumbnail, err := os.CreateTemp(os.TempDir(), "gocrud-thu-*"+extThu)
	if err != nil {
		err = fmt.Errorf("thu: failed to create temp file")
		return
	}
	defer func() {
		_ = tmpThumbnail.Close()
		_ = os.Remove(tmpThumbnail.Name())
	}()

	n, err := io.Copy(tmpOriginal, src)
	if err != nil {
		err = fmt.Errorf("mpf->ori: copy failed")
		return
	} else if n != size {
		err = fmt.Errorf("mpf->ori: incomplete copy")
		return
	}
	_, _ = tmpOriginal.Seek(0, io.SeekStart)

	err = NewMedia(tmpSrc.Name(), tmpOriginal.Name())
	if err != nil {
		err = fmt.Errorf("ori->src: %v", err)
		return
	}
	_, _ = tmpOriginal.Seek(0, io.SeekStart)
	_, _ = tmpSrc.Seek(0, io.SeekStart)

	err = NewMediaPreview(tmpThumbnail.Name(), tmpOriginal.Name())
	if err != nil {
		err = fmt.Errorf("ori->thu: %v", err)
		return
	}
	_, _ = tmpOriginal.Seek(0, io.SeekStart)
	_, _ = tmpThumbnail.Seek(0, io.SeekStart)

	dareConfig := &gocrud.SaveDareFileConfig{
		BaseFolder: env.StaticFolder,
		Ext:        extSrc,
		MasterKey:  gocrud.SHASum256FromString(env.StaticFileMasterKey),
	}

	dareSrc, err := gocrud.SaveDareFile(tmpSrc, dareConfig)
	if err != nil {
		err = fmt.Errorf("failed to save src: %v", err)
		return
	}

	var objSrc FileObject
	objSrc.FromHttpFile(dareSrc)
	err = db.Model(&objSrc).Save(&objSrc).Error
	if err != nil {
		err = fmt.Errorf("failed to save src object file: %v", err)
		return
	}

	dareConfig.Ext = extThu
	dareThu, err := gocrud.SaveDareFile(tmpThumbnail, dareConfig)
	if err != nil {
		err = fmt.Errorf("failed to save thu: %v", err)
		return
	}

	var objThu FileObject
	objThu.FromHttpFile(dareThu)
	err = db.Model(&objThu).Save(&objThu).Error
	if err != nil {
		err = fmt.Errorf("failed to save thu object file: %v", err)
		return
	}

	filenameSrc = string(dareSrc.Filename)
	filenameThu = string(dareThu.Filename)

	return
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
