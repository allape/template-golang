package controller

import (
	"strconv"

	"github.com/allape/gocrud"
	"github.com/allape/golang/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var userGalleryL = l.New("user:gallery")

func SetupUserGalleryController(group *gin.RouterGroup, db *gorm.DB) error {
	err := gocrud.Setup(group, db, userGalleryL.New("crud"), &gocrud.Crud[model.UserGallery]{
		DisableGetOne: false,
		SearchHandlers: gocrud.BaseSearchHandlers(gocrud.SearchHandlers{
			"userId":    gocrud.KeywordEqual("user_id", nil),
			"galleryId": gocrud.KeywordEqual("gallery_id", nil),
		}),
		WillSave: func(record *model.UserGallery, context *gin.Context, db *gorm.DB) {
			if record.UserID == "" {
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

	group.PUT("/all", func(context *gin.Context) {
		var userGalleries []model.UserGallery
		err := context.ShouldBindJSON(&userGalleries)
		if err != nil {
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.BadRequest(), "invalid request")
			return
		}

		if len(userGalleries) == 0 {
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.BadRequest(), "empty set")
			return
		}

		res := db.Model(&model.UserGallery{}).Save(&userGalleries)
		if res.Error != nil {
			userGalleryL.Error().Printf("failed to save user galleries: %v", res)
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.InternalServerError(), res)
			return
		}

		gocrud.MakeOkayDataResponse(context, res.RowsAffected)
	})

	return nil
}
