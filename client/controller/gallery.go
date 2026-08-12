package controller

import (
	"net/http"
	"os"
	"path"
	"slices"
	"strconv"

	"github.com/allape/gocrud"
	"github.com/allape/gogger"
	"github.com/allape/golang/asset"
	"github.com/allape/golang/client/helper"
	"github.com/allape/golang/env"
	"github.com/allape/golang/model"
	"github.com/allape/gophorward"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var galleryl = gogger.New("client:controller:gallery")

// TODO cache for data

type GalleryInfo struct {
	Info model.Gallery `json:"info"`
	Tags []model.Tag   `json:"tags"`
	//Cover *model.Item   `json:"cover"`
}

type GalleryDetailPayload struct {
	Gallery     model.Gallery      `json:"gallery"`
	GalleryTags []model.GalleryTag `json:"galleryTags"`
	Items       []model.Item       `json:"items"`
	ItemTags    []model.ItemTag    `json:"itemTags"`
	Tags        []model.Tag        `json:"tags"`
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

		galleryIds := make([]gocrud.ID, len(galleries))
		for i, gallery := range galleries {
			galleryIds[i] = gallery.ID
		}

		var galleryTags []model.GalleryTag
		if err := db.Model(&model.GalleryTag{}).Where("`gallery_id` IN ?", galleryIds).Find(&galleryTags).Error; err != nil {
			galleryl.Error().Printf("failed to get gallery tags: %v", err)
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.InternalServerError(), "failed to get gallery tags [error]")
			return
		}

		tagIds := make([]gocrud.ID, 0, len(galleryTags))
		for _, galleryTag := range galleryTags {
			if !slices.Contains(tagIds, galleryTag.TagID) {
				tagIds = append(tagIds, galleryTag.TagID)
			}
		}

		var tags []model.Tag
		if err := db.Model(&model.Tag{}).Where("`id` IN ?", tagIds).Find(&tags).Error; err != nil {
			galleryl.Error().Printf("failed to get tags: %v", err)
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.InternalServerError(), "failed to get tags [error]")
			return
		}

		simpleGalleries := make([]GalleryInfo, len(galleries))
		for i, gallery := range galleries {
			var innerGalleryTags []model.Tag
			for _, gt := range galleryTags {
				if gt.GalleryID == gallery.ID {
					for _, tag := range tags {
						if tag.ID == gt.TagID {
							innerGalleryTags = append(innerGalleryTags, tag)
							continue
						}
					}
					continue
				}
			}

			//var items []model.Item
			//if err := db.
			//	Model(&model.Item{}).
			//	Joins("JOIN gallery_items ON gallery_items.item_id = items.id").
			//	Where("gallery_items.gallery_id = ?", gallery.ID).
			//	Limit(1).Find(
			//	&items,
			//).Error; err != nil {
			//	galleryl.Error().Printf("failed to get gallery items: %v", err)
			//	gocrud.MakeErrorResponse(context, gocrud.RestCoder.InternalServerError(), "failed to get gallery items [error]")
			//	return
			//}

			simpleGalleries[i] = GalleryInfo{
				Info: gallery,
				Tags: innerGalleryTags,
				//Cover: gocrud.TernaryFunc[*model.Item](func() bool {
				//	return len(items) > 0
				//}, func() *model.Item {
				//	return &items[0]
				//}, func() *model.Item {
				//	return nil
				//}),
			}
		}

		gocrud.MakeOkayDataResponse(context, simpleGalleries)
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

		if ok, err := helper.UserCanAccessTo(user, &gallery, context, db); !ok {
			if err != nil {
				galleryl.Error().Printf("failed to check access for user %s of gallery %d: %v", user.ID, gallery.ID, err)
			}
			Make403Response(context)
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

		var itemTags []model.ItemTag
		if len(items) > 0 {
			var itemIds = make([]gocrud.ID, len(items))
			for i, item := range items {
				itemIds[i] = item.ID
			}

			if err := db.Model(&model.ItemTag{}).Where("`item_id` IN ?", itemIds).Find(&itemTags).Error; err != nil {
				galleryl.Error().Printf("failed to get item tags by gallery id %d: %v", galleryId, err)
				gocrud.MakeErrorResponse(context, gocrud.RestCoder.InternalServerError(), "failed to find item tags [error]")
				return
			}
			payload.ItemTags = itemTags
		}

		var galleryTags []model.GalleryTag
		if err := db.Model(&model.GalleryTag{}).Where("`gallery_id` = ?", galleryId).Find(&galleryTags).Error; err != nil {
			galleryl.Error().Printf("failed to get gallery tags by gallery id %d: %v", galleryId, err)
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.InternalServerError(), "failed to find gallery tags [error]")
			return
		}
		payload.GalleryTags = galleryTags

		var tagIds = make([]gocrud.ID, 0, len(itemTags)+len(galleryTags))
		for _, tag := range itemTags {
			if slices.Contains(tagIds, tag.TagID) {
				continue
			}
			tagIds = append(tagIds, tag.TagID)
		}
		for _, tag := range galleryTags {
			if slices.Contains(tagIds, tag.TagID) {
				continue
			}
			tagIds = append(tagIds, tag.TagID)
		}

		if len(tagIds) > 0 {
			var tags []model.Tag
			if err := db.Model(&model.Tag{}).Where("`id` IN ?", tagIds).Find(&tags).Error; err != nil {
				galleryl.Error().Printf("failed to get tags by gallery id %d: %v", galleryId, err)
				gocrud.MakeErrorResponse(context, gocrud.RestCoder.InternalServerError(), "failed to find tags [error]")
				return
			}
			payload.Tags = tags
		}

		gocrud.MakeOkayDataResponse(context, payload)
	})

	const ErrorMessageHeaderName = "X-Error-Message"
	makeErrorImageResponse := func(context *gin.Context, encodedPngImage []byte, status int, message string) {
		NoCache(context)
		context.Header(ErrorMessageHeaderName, message)
		context.Data(status, asset.MIME, gocrud.Ternary(encodedPngImage == nil, asset.NoImage, encodedPngImage))
	}

	// retrieving the first image of the gallery of :galleryId when :itemId is 0
	group.GET("/image/:galleryId/:itemId", func(context *gin.Context) {
		user, ok := gophorward.GinGetUser(context)
		if !ok {
			//Make401Response(context)
			makeErrorImageResponse(context, asset.DameMan, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}

		galleryId, err := strconv.ParseInt(context.Param("galleryId"), 10, 64)
		if err != nil || galleryId <= 0 {
			makeErrorImageResponse(context, nil, http.StatusBadRequest, "invalid gallery id")
			return
		}

		var gallery model.Gallery
		if err := db.Model(&gallery).Where("`id` = ?", galleryId).First(&gallery).Error; err != nil {
			galleryl.Error().Printf("failed to get gallery: %v", err)
			makeErrorImageResponse(context, nil, http.StatusInternalServerError, "failed to find gallery [error]")
			return
		} else if gallery.ID == 0 {
			makeErrorImageResponse(context, nil, http.StatusNotFound, "failed to find gallery")
			return
		}

		if ok, err := helper.UserCanAccessTo(user, &gallery, context, db); !ok {
			if err != nil {
				galleryl.Error().Printf("failed to check access for user %s of gallery %d: %v", user.ID, gallery.ID, err)
			}
			//Make403Response(context)
			makeErrorImageResponse(context, asset.DameMan, http.StatusForbidden, http.StatusText(http.StatusForbidden))
			return
		}

		itemId, err := strconv.ParseInt(context.Param("itemId"), 10, 64)
		if err != nil || itemId < 0 {
			makeErrorImageResponse(context, nil, http.StatusBadRequest, "invalid item id")
			return
		}

		var items []model.Item

		if itemId == 0 {
			err = db.Model(&model.Item{}).
				Joins("JOIN gallery_items ON gallery_items.item_id = items.id").
				Where("gallery_items.gallery_id = ?", gallery.ID).
				Order("`priority` ASC, `updated_at` DESC").
				Limit(1).
				Find(&items).Error
		} else {
			err = db.Model(&model.Item{}).Where("`id` = ?", itemId).Limit(1).Find(&items).Error
		}

		if err != nil {
			galleryl.Error().Printf("failed to get item by item id %d: %v", itemId, err)
			makeErrorImageResponse(context, nil, http.StatusInternalServerError, "failed to find item [error]")
			return
		}

		if len(items) == 0 {
			makeErrorImageResponse(context, nil, http.StatusNotFound, "failed to find item")
			return
		}

		item := items[0]

		if item.Src == "" {
			galleryl.Error().Printf("item src is empty for %d", item.ID)
			makeErrorImageResponse(context, nil, http.StatusInternalServerError, "src is empty")
			return
		}

		file, err := os.Open(path.Join(env.StaticFolder, item.Src))
		if err != nil {
			galleryl.Error().Printf("failed to open gallery item file by gallery id %d and item id %d: %v", galleryId, itemId, err)
			makeErrorImageResponse(context, nil, http.StatusInternalServerError, "failed to open gallery item file [error]")
			return
		}

		var fileObject model.FileObject
		if err := db.Model(&fileObject).Where("`filename` = ?", item.Src).First(&fileObject).Error; err != nil {
			galleryl.Error().Printf("failed to get file object by item id %d: %v", itemId, err)
			makeErrorImageResponse(context, nil, http.StatusInternalServerError, "failed to find file object [error]")
			return
		}

		serveFunc, err := gocrud.NewDareHttpServeFunc(file, fileObject.HttpFileSystemObjectBase.ToHttpFile())
		if err != nil {
			galleryl.Error().Printf("failed to decrypt file object by item id %d: %v", itemId, err)
			makeErrorImageResponse(context, nil, http.StatusInternalServerError, "failed to decrypt file object [error]")
			return
		}

		serveFunc(context.Writer, context.Request)
	})

	return nil
}
