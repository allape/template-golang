package controller

import (
	"net/http"
	"path"
	"slices"

	"github.com/allape/gocrud"
	"github.com/allape/gogger"
	"github.com/allape/golang/asset"
	"github.com/allape/golang/env"
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
		if err := db.Model(&model.Gallery{}).Where("`created_by` = ?", user.ID).Find(&galleries).Error; err != nil {
			galleryl.Error().Printf("failed to get gallery: %v", err)
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.InternalServerError(), "failed to get gallery [error]")
			return
		}

		gocrud.MakeOkayDataResponse(context, galleries)
	})

	group.GET("/detail/:galleryId", func(context *gin.Context) {
		user, ok := gophorward.GinGetUser(context)
		if !ok {
			Make401Response(context)
			return
		}

		galleryId := gocrud.Pick(gocrud.IDsFromCommaSeparatedString(context.Param("galleryId")), 0, 0)
		if galleryId == 0 {
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.BadRequest(), "invalid gallery id")
			return
		}

		var payload GalleryDetailPayload

		var gallery model.Gallery
		if err := db.Model(&gallery).Where("`id` = ?", galleryId).First(&gallery).Error; err != nil {
			galleryl.Error().Printf("failed to get gallery by %d: %v", galleryId, err)
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
			"`id` IN (SELECT `gi`.`item_id` FROM `gallery_items` `gi` WHERE `gi`.`gallery_id` = ?)",
			galleryId,
		).Find(&items).Error; err != nil {
			galleryl.Error().Printf("failed to get item by gallery id %d: %v", galleryId, err)
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
		if err := db.Model(&model.ItemTag{}).Where("`item_id` IN ?", itemIds).Find(&itemTags).Error; err != nil {
			galleryl.Error().Printf("failed to get item tags by gallery id %d: %v", galleryId, err)
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
		if err := db.Model(&model.Tag{}).Where("`id` IN ?", tagIds).Find(&tags).Error; err != nil {
			galleryl.Error().Printf("failed to get tags by gallery id %d: %v", galleryId, err)
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.InternalServerError(), "failed to find tags [error]")
			return
		}

		payload.Tags = tags

		gocrud.MakeOkayDataResponse(context, payload)
	})

	// :itemId can be empty for retrieving the first image of the gallery of :galleryId
	group.GET("/image/:galleryId/:itemId", func(context *gin.Context) {
		user, ok := gophorward.GinGetUser(context)
		if !ok {
			Make401Response(context)
			return
		}

		galleryId := gocrud.Pick(gocrud.IDsFromCommaSeparatedString(context.Param("galleryId")), 0, 0)
		if galleryId == 0 {
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.BadRequest(), "invalid gallery id")
			return
		}

		var gallery model.Gallery
		if err := db.Model(&gallery).Where("`id` = ?", galleryId).First(&gallery).Error; err != nil {
			galleryl.Error().Printf("failed to get gallery: %v", err)
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.InternalServerError(), "failed to find gallery [error]")
			return
		} else if gallery.ID == 0 {
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.BadRequest(), "failed to find gallery")
			return
		}

		if !gallery.IsPublic && gallery.CreatedBy != user.ID {
			//gocrud.MakeErrorResponse(context, gocrud.RestCoder.FromStatus(http.StatusForbidden), http.StatusText(http.StatusForbidden))
			context.Data(http.StatusForbidden, asset.MIME, asset.DameMan)
			return
		}

		var galleryItem model.GalleryItem

		itemId := gocrud.Pick(gocrud.IDsFromCommaSeparatedString(context.Param("itemId")), 0, 0)

		//if itemId == 0 {
		//	gocrud.MakeErrorResponse(context, gocrud.RestCoder.BadRequest(), "invalid item id")
		//	return
		//}

		if err := gocrud.TernaryFunc(
			func() bool {
				return itemId == 0
			},
			func() error {
				return db.Model(&galleryItem).Where("`gallery_id` = ?", galleryId).Order("`created_at` ASC").First(&galleryItem).Error
			},
			func() error {
				return db.Model(&galleryItem).Where("`gallery_id` = ? AND `item_id` = ?", galleryId, itemId).First(&galleryItem).Error
			},
		); err != nil {
			galleryl.Error().Printf("failed to get gallery item by gallery id %d and item id %d: %v", galleryId, itemId, err)
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.InternalServerError(), "failed to find gallery item [error]")
			return
		} else if galleryItem.ItemID == 0 {
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.BadRequest(), "failed to find gallery item")
			return
		}

		var item model.Item
		if err := db.Model(&item).Where("`id` = ?", galleryItem.ItemID).First(&item).Error; err != nil {
			galleryl.Error().Printf("failed to get item by item id %d: %v", galleryItem.ItemID, err)
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.InternalServerError(), "failed to find item [error]")
			return
		} else if item.ID == 0 {
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.BadRequest(), "failed to find item")
			return
		}

		if item.Src == "" {
			galleryl.Error().Printf("item src is empty for %d", item.ID)
			context.Data(http.StatusNotFound, asset.MIME, asset.NoImage)
			return
		}

		context.File(path.Join(env.StaticFolder, item.Src))
	})

	return nil
}
