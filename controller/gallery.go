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

	// /items/1?itemIds=1,2,3...
	group.POST("/items/{galleryId}", func(c *gin.Context) {
		galleryId := gocrud.Pick(gocrud.IDsFromCommaSeparatedString(c.Param("galleryId")), 0, 0)
		if galleryId == 0 {
			gocrud.MakeErrorResponse(c, gocrud.RestCoder.BadRequest(), "galleryId is required")
			return
		}

		var gallery model.Gallery
		if err := db.Model(&gallery).Where("id = ?", galleryId).First(&gallery).Error; err != nil {
			galleryl.Error().Printf("failed to find gallery %d: %v", galleryId, err)
			gocrud.MakeErrorResponse(c, gocrud.RestCoder.BadRequest(), "gallery not found [error]")
			return
		} else if gallery.ID == 0 {
			gocrud.MakeErrorResponse(c, gocrud.RestCoder.BadRequest(), "gallery not found")
			return
		}

		itemIds := gocrud.IDsFromCommaSeparatedString(c.Query("itemIds"))
		if len(itemIds) == 0 {
			gocrud.MakeErrorResponse(c, gocrud.RestCoder.BadRequest(), "itemIds is required")
			return
		}

		var items []model.Item
		if err := db.Model(&model.Item{}).Where("id IN ?", itemIds).Find(&items).Error; err != nil {
			galleryl.Error().Printf("failed to find items %v: %v", itemIds, err)
			gocrud.MakeErrorResponse(c, gocrud.RestCoder.InternalServerError(), "failed to find items [error]")
			return
		} else if len(items) == 0 {
			gocrud.MakeErrorResponse(c, gocrud.RestCoder.BadRequest(), "items is empty")
			return
		}

		galleryItems := make([]model.GalleryItem, len(items))
		for i, item := range items {
			galleryItems[i] = model.GalleryItem{
				GalleryID: gallery.ID,
				ItemID:    item.ID,
			}
		}

		res := db.Save(&galleryItems)
		if res.Error != nil {
			galleryl.Error().Printf("failed to save gallery items: %v", res.Error)
			gocrud.MakeErrorResponse(c, gocrud.RestCoder.InternalServerError(), "failed to save gallery items [error]")
			return
		}

		gocrud.MakeOkayDataResponse(c, res.RowsAffected)
	})

	// /items/1?itemIds=1,2,3...
	group.DELETE("/items/{galleryId}", func(c *gin.Context) {
		galleryId := gocrud.Pick(gocrud.IDsFromCommaSeparatedString(c.Param("galleryId")), 0, 0)
		if galleryId == 0 {
			gocrud.MakeErrorResponse(c, gocrud.RestCoder.BadRequest(), "galleryId is required")
			return
		}

		var gallery model.Gallery
		if err := db.Model(&gallery).Where("id = ?", galleryId).First(&gallery).Error; err != nil {
			galleryl.Error().Printf("failed to find gallery %d: %v", galleryId, err)
			gocrud.MakeErrorResponse(c, gocrud.RestCoder.BadRequest(), "gallery not found [error]")
			return
		} else if gallery.ID == 0 {
			gocrud.MakeErrorResponse(c, gocrud.RestCoder.BadRequest(), "gallery not found")
			return
		}

		itemIds := gocrud.IDsFromCommaSeparatedString(c.Query("itemIds"))
		if len(itemIds) == 0 {
			gocrud.MakeErrorResponse(c, gocrud.RestCoder.BadRequest(), "itemIds is required")
			return
		}

		res := db.Model(&model.GalleryItem{}).Delete("gallery_id = ? AND item_id IN ?", gallery.ID, itemIds)
		if res.Error != nil {
			galleryl.Error().Printf("failed to delete gallery items: %v", res.Error)
			gocrud.MakeErrorResponse(c, gocrud.RestCoder.InternalServerError(), "failed to delete gallery items [error]")
			return
		}

		gocrud.MakeOkayDataResponse(c, res.RowsAffected)
	})

	group.DELETE("/empty-gallery/{galleryId}", func(c *gin.Context) {
		galleryId := gocrud.Pick(gocrud.IDsFromCommaSeparatedString(c.Param("galleryId")), 0, 0)
		if galleryId == 0 {
			gocrud.MakeErrorResponse(c, gocrud.RestCoder.BadRequest(), "galleryId is required")
			return
		}

		var gallery model.Gallery
		if err := db.Model(&gallery).Where("id = ?", galleryId).First(&gallery).Error; err != nil {
			galleryl.Error().Printf("failed to find gallery %d: %v", galleryId, err)
			gocrud.MakeErrorResponse(c, gocrud.RestCoder.BadRequest(), "gallery not found [error]")
			return
		} else if gallery.ID == 0 {
			gocrud.MakeErrorResponse(c, gocrud.RestCoder.BadRequest(), "gallery not found")
			return
		}

		res := db.Model(&model.GalleryItem{}).Delete("gallery_id = ?", gallery.ID)
		if res.Error != nil {
			galleryl.Error().Printf("failed to delete gallery item for id %d: %v", galleryId, res.Error)
			gocrud.MakeErrorResponse(c, gocrud.RestCoder.InternalServerError(), "failed to delete gallery item [error]")
			return
		}

		gocrud.MakeOkayDataResponse(c, res.RowsAffected)
	})

	return nil
}
