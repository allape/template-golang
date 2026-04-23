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

var iteml = gogger.New("controller:item")

func SetupItemController(group *gin.RouterGroup, db *gorm.DB) error {
	err := gocrud.New(group, db, gocrud.Crud[model.Item]{
		DefaultPageSize: DefaultPageSize,
		SearchHandlers: map[string]gocrud.SearchHandler{
			"in_id":             gocrud.KeywordIDIn("id", gocrud.OverflowedArrayTrimmerFilter[gocrud.ID](DefaultPageSize)),
			"like_name":         gocrud.KeywordLike("name", nil),
			"createBy":          gocrud.KeywordEqual("created_by", nil),
			"orderBy_priority":  gocrud.SortBy("priority"),
			"orderBy_updatedAt": gocrud.SortBy("updated_at"),
			"deleted":           gocrud.NewSoftDeleteSearchHandler(""),
			"orderByDefault": func(db *gorm.DB, values []string, with url.Values) *gorm.DB {
				return db.Order("`priority` DESC, `updated_at` DESC")
			},
			"in_galleryId": func(db *gorm.DB, values []string, with url.Values) *gorm.DB {
				if value, ok := gocrud.PickFirstValuableString(values); ok {
					ids := gocrud.IDsFromCommaSeparatedString(value)
					if len(ids) == 0 {
						return db.Where("1 != 1")
					}
					return db.Where("id IN (SELECT gi.item_id FROM gallery_items gi WHERE gi.gallery_id IN ?)", ids)
				}
				return db
			},
		},
		OnDelete: gocrud.NewSoftDeleteHandler[model.Item](gocrud.RestCoder),
		WillSave: func(record *model.Item, context *gin.Context, db *gorm.DB) {
			record.Name = strings.TrimSpace(record.Name)
			if len(record.Name) > model.MaxItemNameLength {
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

func SetupItemTagController(group *gin.RouterGroup, db *gorm.DB) error {
	return setupDualPrimaryKeyModelController[model.ItemTag](
		group, db,
		"ItemID", "TagID",
		"itemId", "tagId",
		"item_id", "tag_id",
		iteml, DefaultPageSize,
	)
}
