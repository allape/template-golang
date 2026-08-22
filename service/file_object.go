package service

import (
	"fmt"
	"io"
	"mime"
	"os"
	"strings"

	"github.com/allape/gocrud"
	"github.com/allape/golang/env"
	"github.com/allape/golang/model"
	"gorm.io/gorm"
)

var fileObjectL = l.New("fileobject")

var FileObjectHandler *gocrud.HttpFileSystemObjectHandler

func SetupFileObjectHandlerService(db *gorm.DB) error {
	var err error
	FileObjectHandler, err = gocrud.NewHttpFileSystemObjectHandler[model.FileObject](
		db, fileObjectL.New("handler"),
		env.StaticFolder,
		gocrud.SHASum256FromString(env.StaticFileMasterKey),
		gocrud.SHASum256FromString(env.StaticFileHashSalt),
		"",
	)
	if err != nil {
		return err
	}
	return nil
}

func SaveFile(src io.ReadSeeker, size int64, ext string) (filenameSrc, filenameThu string, err error) {
	mimeStr := mime.TypeByExtension(ext)
	extSrc := ext
	extThu := ext

	switch {
	case ext == ".gif":
		extSrc = ".gif"
		extThu = ".jpg"
	case ext == ".webp":
		extSrc = ".webp"
		extThu = ".webp"
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

	err = model.NewMedia(tmpSrc.Name(), tmpOriginal.Name())
	if err != nil {
		err = fmt.Errorf("ori->src: %v", err)
		return
	}
	_, _ = tmpOriginal.Seek(0, io.SeekStart)
	_, _ = tmpSrc.Seek(0, io.SeekStart)

	err = model.NewMediaPreview(tmpThumbnail.Name(), tmpOriginal.Name())
	if err != nil {
		err = fmt.Errorf("ori->thu: %v", err)
		return
	}
	_, _ = tmpOriginal.Seek(0, io.SeekStart)
	_, _ = tmpThumbnail.Seek(0, io.SeekStart)

	dareSrc, err := FileObjectHandler.Save(tmpSrc, extSrc, 0, "")
	if err != nil {
		err = fmt.Errorf("failed to save src: %v", err)
		return
	}

	dareThu, err := FileObjectHandler.Save(tmpThumbnail, extThu, 0, "")
	if err != nil {
		err = fmt.Errorf("failed to save thu: %v", err)
		return
	}

	filenameSrc = string(dareSrc.Name)
	filenameThu = string(dareThu.Name)

	return
}
