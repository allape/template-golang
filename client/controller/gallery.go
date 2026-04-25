package controller

import (
	"net/http"
	"slices"

	"github.com/allape/gocrud"
	"github.com/allape/gogger"
	"github.com/allape/golang/model"
	"github.com/allape/gophorward"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var galleryl = gogger.New("client:controller:gallery")

type GalleryDetailPayload struct {
	Gallery  model.Gallery   `json:"gallery"`
	Items    []model.Item    `json:"items"`
	ItemTags []model.ItemTag `json:"itemTags"`
	Tags     []model.Tag     `json:"tags"`
}

func SetupGalleryController(group *gin.RouterGroup, db *gorm.DB) error {
	group.GET("/all", func(context *gin.Context) {
		user, ok := gophorward.GinGetUser(context)
		if !ok {
			Make401Response(context)
			return
		}

		var galleries []model.Gallery
		if err := db.Model(&model.Gallery{}).Where("user_id = ?", user.ID).Find(&galleries).Error; err != nil {
			galleryl.Error().Printf("failed to get gallery: %v", err)
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.InternalServerError(), "failed to get gallery [error]")
			return
		}

		gocrud.MakeOkayDataResponse(context, galleries)
	})

	group.GET("/:id", func(context *gin.Context) {
		user, ok := gophorward.GinGetUser(context)
		if !ok {
			Make401Response(context)
			return
		}

		id := gocrud.Pick(gocrud.IDsFromCommaSeparatedString(context.Param("id")), 0, 0)
		if id == 0 {
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.BadRequest(), "invalid id")
			return
		}

		var payload GalleryDetailPayload

		var gallery model.Gallery
		if err := db.Model(&gallery).Where("id = ?", id).First(&gallery).Error; err != nil {
			galleryl.Error().Printf("failed to get gallery by %d: %v", id, err)
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.InternalServerError(), "failed to find gallery [error]")
			return
		} else if gallery.ID == 0 {
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.BadRequest(), "failed to find gallery")
			return
		}

		if !gallery.IsPublic && gallery.CreatedBy != user.ID {
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.FromStatus(http.StatusForbidden), http.StatusText(http.StatusForbidden))
			return
		}

		payload.Gallery = gallery

		var items []model.Item
		if err := db.Model(&model.Item{}).Where(
			"id IN (SELECT gi.item_id FROM gallery_items gi WHERE gi.gallery_id = ?)",
			id,
		).Find(&items).Error; err != nil {
			galleryl.Error().Printf("failed to get item by gallery id %d: %v", id, err)
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.InternalServerError(), "failed to find items [error]")
			return
		}

		payload.Items = items

		if len(items) == 0 {
			gocrud.MakeOkayDataResponse(context, payload)
			return
		}

		var itemIds = make([]gocrud.ID, len(items))
		for i, item := range items {
			itemIds[i] = item.ID
		}

		var itemTags []model.ItemTag
		if err := db.Model(&model.ItemTag{}).Where("item_id IN ?", itemIds).Find(&itemTags).Error; err != nil {
			galleryl.Error().Printf("failed to get item tags by gallery id %d: %v", id, err)
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.InternalServerError(), "failed to find item tags [error]")
			return
		}

		payload.ItemTags = itemTags

		if len(itemTags) == 0 {
			gocrud.MakeOkayDataResponse(context, payload)
			return
		}

		var tagIds = make([]gocrud.ID, 0, len(itemTags))
		for _, tag := range itemTags {
			if slices.Contains(tagIds, tag.TagID) {
				continue
			}
			tagIds = append(tagIds, tag.TagID)
		}

		var tags []model.Tag
		if err := db.Model(&model.Tag{}).Where("id IN ?", tagIds).Find(&tags).Error; err != nil {
			galleryl.Error().Printf("failed to get tags by gallery id %d: %v", id, err)
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.InternalServerError(), "failed to find tags [error]")
			return
		}

		payload.Tags = tags

		gocrud.MakeOkayDataResponse(context, payload)
	})

	return nil
}
