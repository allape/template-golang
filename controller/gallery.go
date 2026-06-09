package controller

import (
	"strings"

	"github.com/allape/gocrud"
	"github.com/allape/golang/model"
	"github.com/allape/gophorward"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var galleryL = l.New("gallery")

func SetupGalleryController(group *gin.RouterGroup, db *gorm.DB) error {
	err := gocrud.Setup(group, db, galleryL.New("crud"), &gocrud.Crud[model.Gallery]{
		EnableGetAll: true,
		SearchHandlers: gocrud.BaseSearchHandlers(gocrud.SearchHandlers{
			"like_name": gocrud.KeywordLike("name", nil),
			"isPublic":  gocrud.KeywordEqual("is_public", nil),
			"createBy":  gocrud.KeywordEqual("created_by", nil),
		}),
		WillSave: func(record *model.Gallery, context *gin.Context, db *gorm.DB) {
			record.Name = strings.TrimSpace(record.Name)
			if record.Name == "" {
				gocrud.MakeErrorResponse(context, gocrud.RestCoder.BadRequest(), "name cannot be empty")
				return
			} else if len(record.Name) > model.MaxGalleryNameLength {
				gocrud.MakeErrorResponse(context, gocrud.RestCoder.BadRequest(), "name too long")
				return
			}

			record.CreatedBy = gophorward.UserID(strings.TrimSpace(string(record.CreatedBy)))
			if record.CreatedBy == "" {
				record.CreatedBy = "0"
			}
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func SetupGalleryItemController(group *gin.RouterGroup, db *gorm.DB) error {
	return gocrud.SetupDualPrimaryKeyModelController[model.GalleryItem](
		group, db, galleryL.New("item"),
		"GalleryID", "ItemID",
		"gallery_id", "item_id",
	)
}
