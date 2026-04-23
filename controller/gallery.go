package controller

import (
	"net/url"
	"strings"

	"github.com/allape/gocrud"
	"github.com/allape/gogger"
	"github.com/allape/golang/model"
	"github.com/allape/gophorward"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var galleryl = gogger.New("controller:gallery")

func SetupGalleryController(group *gin.RouterGroup, db *gorm.DB) error {
	err := gocrud.New(group, db, gocrud.Crud[model.Gallery]{
		DefaultPageSize: DefaultPageSize,
		EnableGetAll:    true,
		WillGetAll:      willGetAllInID,
		SearchHandlers: map[string]gocrud.SearchHandler{
			"in_id":             gocrud.KeywordIDIn("id", gocrud.OverflowedArrayTrimmerFilter[gocrud.ID](DefaultPageSize)),
			"like_name":         gocrud.KeywordLike("name", nil),
			"isPublic":          gocrud.KeywordEqual("is_public", nil),
			"createBy":          gocrud.KeywordEqual("created_by", nil),
			"orderBy_priority":  gocrud.SortBy("priority"),
			"orderBy_updatedAt": gocrud.SortBy("updated_at"),
			"orderByDefault": func(db *gorm.DB, values []string, with url.Values) *gorm.DB {
				return db.Order("`priority` DESC, `updated_at` DESC")
			},
			"deleted": gocrud.NewSoftDeleteSearchHandler(""),
		},
		OnDelete: gocrud.NewSoftDeleteHandler[model.Gallery](gocrud.RestCoder),
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
	return setupDualPrimaryKeyModelController[model.GalleryItem](
		group, db,
		"GalleryID", "ItemID",
		"galleryId", "itemId",
		"gallery_id", "item_id",
		galleryl, DefaultPageSize,
	)
}
