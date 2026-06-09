package controller

import (
	"strings"

	"github.com/allape/gocrud"
	"github.com/allape/golang/model"
	"github.com/allape/gophorward"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var itemL = l.New("item")

func SetupItemController(group *gin.RouterGroup, db *gorm.DB) error {
	err := gocrud.Setup(group, db, itemL.New("crud"), &gocrud.Crud[model.Item]{
		SearchHandlers: gocrud.BaseSearchHandlers(gocrud.SearchHandlers{
			"like_name": gocrud.KeywordLike("name", nil),
			"createBy":  gocrud.KeywordEqual("created_by", nil),
			"in_galleryId": func(db *gorm.DB, values []string, _ *gin.Context) (*gorm.DB, error) {
				if value, ok := gocrud.PickFirstValuableString(values); ok {
					ids := gocrud.IDsFromCommaSeparatedString(value)
					if len(ids) == 0 {
						return db, gocrud.NotArrayError
					}
					return db.Where("id IN (SELECT gi.item_id FROM gallery_items gi WHERE gi.gallery_id IN ?)", ids), nil
				}
				return db, nil
			},
		}),
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
	return gocrud.SetupDualPrimaryKeyModelController[model.ItemTag](
		group, db, itemL.New("tag"),
		"ItemID", "TagID",
		"item_id", "tag_id",
	)
}
