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
					return db.Where("id IN (SELECT id FROM gallery_items gi WHERE gi.gallery_id IN ?)", ids)
				}
				return db
			},
		},
		OnDelete: gocrud.NewSoftDeleteHandler[model.Item](gocrud.RestCoder),
		WillSave: func(record *model.Item, context *gin.Context, db *gorm.DB) {
			record.Name = strings.TrimSpace(record.Name)
			if record.Name == "" {
				gocrud.MakeErrorResponse(context, gocrud.RestCoder.BadRequest(), "name cannot be empty")
				return
			} else if len(record.Name) > model.MaxItemNameLength {
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

	// /tags?itemIds=1,2,3,...
	group.GET("/tags", func(c *gin.Context) {
		itemIds := gocrud.IDsFromCommaSeparatedString(c.Query("itemIds"))
		if len(itemIds) < 0 {
			gocrud.MakeErrorResponse(c, gocrud.RestCoder.BadRequest(), "itemIds is required")
			return
		}

		var items []model.Item
		if err := db.Model(&model.Item{}).Where("`id` IN ?", itemIds).Find(&items).Error; err != nil {
			iteml.Error().Printf("failed to find item %v: %v", itemIds, err)
			gocrud.MakeErrorResponse(c, gocrud.RestCoder.BadRequest(), "item not found [error]")
			return
		} else if len(items) == 0 {
			gocrud.MakeErrorResponse(c, gocrud.RestCoder.NotFound(), "item not found")
			return
		}

		parsedItemIds := make([]gocrud.ID, len(items))
		for index, item := range items {
			parsedItemIds[index] = item.ID
		}

		var tags []model.Tag
		if //goland:noinspection SqlNoDataSourceInspection
		err := db.Raw("SELECT t.* FROM `tags` t WHERE t.`id` IN (SELECT it.tag_id FROM `item_tags` it WHERE it.item_id IN ?)", parsedItemIds).Scan(&tags).Error; err != nil {
			iteml.Error().Printf("failed to find tags for item ids %v: %v", parsedItemIds, err)
			gocrud.MakeErrorResponse(c, gocrud.RestCoder.InternalServerError(), "tag not found [error]")
			return
		}

		gocrud.MakeOkayDataResponse(c, tags)
	})

	// /tags/1?tagIds=1,2,3,...
	group.POST("/tags/{itemId}", func(c *gin.Context) {
		itemId := gocrud.Pick(gocrud.IDsFromCommaSeparatedString(c.Param("itemId")), 0, 0)
		if itemId == 0 {
			gocrud.MakeErrorResponse(c, gocrud.RestCoder.BadRequest(), "itemId is required")
			return
		}

		var item model.Item
		if err := db.Model(&item).Where("`id` = ?", itemId).First(&item).Error; err != nil {
			iteml.Error().Printf("failed to find item %d: %v", itemId, err)
			gocrud.MakeErrorResponse(c, gocrud.RestCoder.InternalServerError(), "item not found [error]")
			return
		} else if item.ID == 0 {
			gocrud.MakeErrorResponse(c, gocrud.RestCoder.NotFound(), "item not found")
			return
		}

		tagIds := gocrud.IDsFromCommaSeparatedString(c.Query("tagIds"))

		count := int64(0)

		err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&model.ItemTag{}).Delete("item_id = ?", itemId).Error; err != nil {
				return err
			}

			if len(tagIds) == 0 {
				return nil
			}

			var parsedTags []model.Tag
			if err := tx.Model(&model.Tag{}).Where("`id` IN ?", tagIds).Find(&parsedTags).Error; err != nil {
				return err
			}

			if len(parsedTags) == 0 {
				return nil
			}

			itemTags := make([]model.ItemTag, len(parsedTags))
			for i, tag := range parsedTags {
				itemTags[i] = model.ItemTag{
					ItemID: item.ID,
					TagID:  tag.ID,
				}
			}

			res := tx.Save(&itemTags)
			if res.Error != nil {
				return res.Error
			}

			count = res.RowsAffected

			return nil
		})
		if err != nil {
			iteml.Error().Printf("failed to save tags for item %d: %v", itemId, err)
			gocrud.MakeErrorResponse(c, gocrud.RestCoder.InternalServerError(), "failed to save item tags [error]")
			return
		}

		gocrud.MakeOkayDataResponse(c, count)
	})

	return nil
}
