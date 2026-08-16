package controller

import (
	"fmt"
	"strings"

	"github.com/allape/gocrud"
	"github.com/allape/golang/model"
	"github.com/allape/gophorward"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var galleryL = l.New("gallery")

func SetupGalleryController(group *gin.RouterGroup, db *gorm.DB) error {
	NameDuplicateCheckFunc, err := gocrud.NewDuplicateFieldCheckFunc[model.Gallery](db, galleryL, "Name")
	if err != nil {
		return err
	}

	err = gocrud.Setup(group, db, galleryL.New("crud"), &gocrud.Crud[model.Gallery]{
		EnableGetAll: true,
		SearchHandlers: gocrud.BaseSearchHandlers(gocrud.SearchHandlers{
			"like_name":     gocrud.KeywordLike("name", nil),
			"like_keywords": gocrud.KeywordLike("keywords", nil),
			"keywords": func(db *gorm.DB, values []string, context *gin.Context) (*gorm.DB, error) {
				if value, ok := gocrud.PickFirstValuableString(values); ok {
					likeValue := fmt.Sprintf("%%%s%%", value)
					return db.Where(
						"`keywords` LIKE ? OR `name` LIKE ?",
						likeValue,
						likeValue,
					), nil
				}
				return db, nil
			},
			"isPublic": gocrud.KeywordEqual("is_public", nil),
			"createBy": gocrud.KeywordEqual("created_by", nil),
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

			err = NameDuplicateCheckFunc(context, record)
			if err != nil {
				return
			}

			record.Keywords = strings.TrimSpace(record.Keywords)
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
	return gocrud.SetupM2MConnectorController[model.GalleryItem](
		group, db, galleryL.New("item"),
		"GalleryID", "ItemID",
		nil,
	)
}

func SetupGalleryTagController(group *gin.RouterGroup, db *gorm.DB) error {
	return gocrud.SetupM2MConnectorController[model.GalleryTag](
		group, db, galleryL.New("tag"),
		"GalleryID", "TagID",
		nil,
	)
}
