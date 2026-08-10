package controller

import (
	"strconv"
	"strings"

	"github.com/allape/gocrud"
	"github.com/allape/golang/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var shareLinkL = l.New("sharelink")

func SetupShareLinkController(group *gin.RouterGroup, db *gorm.DB) error {
	err := gocrud.Setup(group, db, shareLinkL.New("crud"), &gocrud.Crud[model.ShareLink]{
		SearchHandlers: gocrud.BaseSearchHandlers(gocrud.SearchHandlers{
			"like_name": gocrud.KeywordLike("name", nil),
		}),
		WillSave: func(record *model.ShareLink, context *gin.Context, db *gorm.DB) {
			record.Name = strings.TrimSpace(record.Name)
			if record.Name == "" {
				gocrud.MakeErrorResponse(context, gocrud.RestCoder.BadRequest(), "name is required")
				return
			}

			if record.ID == 0 {
				record.NonceID = strings.ReplaceAll(uuid.NewString(), "-", "")
			}
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func SetupShareLinkGalleryController(group *gin.RouterGroup, db *gorm.DB) error {
	return gocrud.SetupM2MConnectorController[model.ShareLinkGallery](
		group, db, shareLinkL.New("gallery"),
		"ShareID", "GalleryID",
		nil,
	)
}

var userGalleryL = l.New("user:gallery")

func SetupUserGalleryController(group *gin.RouterGroup, db *gorm.DB) error {
	err := gocrud.Setup(group, db, userGalleryL.New("crud"), &gocrud.Crud[model.UserGallery]{
		DisableGetOne: false,
		SearchHandlers: gocrud.BaseSearchHandlers(gocrud.SearchHandlers{
			"userId":    gocrud.KeywordEqual("user_id", nil),
			"galleryId": gocrud.KeywordEqual("gallery_id", nil),
		}),
		WillSave: func(record *model.UserGallery, context *gin.Context, db *gorm.DB) {
			if record.UserID == 0 {
				gocrud.MakeErrorResponse(context, gocrud.RestCoder.BadRequest(), "userId is required")
				return
			} else if record.GalleryID == 0 {
				gocrud.MakeErrorResponse(context, gocrud.RestCoder.BadRequest(), "galleryId is required")
				return
			}
		},
		OnDelete: func(context *gin.Context, db *gorm.DB) bool {
			userId, err := strconv.ParseInt(context.Query("userId"), 10, 64)
			if err != nil {
				gocrud.MakeErrorResponse(context, gocrud.RestCoder.BadRequest(), "userId is invalid")
				return false
			}

			galleryId, err := strconv.ParseInt(context.Query("galleryId"), 10, 64)
			if err != nil {
				gocrud.MakeErrorResponse(context, gocrud.RestCoder.BadRequest(), "galleryId is invalid")
				return false
			}

			return db.Delete(&model.UserGallery{}, "`user_id` = ? AND `gallery_id` = ?", userId, galleryId).Error != nil
		},
	})
	if err != nil {
		return err
	}

	return nil
}
