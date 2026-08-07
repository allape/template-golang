package controller

import (
	"strings"

	"github.com/allape/gocrud"
	"github.com/allape/golang/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var accessGroupL = l.New("accessgroup")

func SetupAccessGroupController(group *gin.RouterGroup, db *gorm.DB) error {
	err := gocrud.Setup(group, db, itemL.New("crud"), &gocrud.Crud[model.AccessGroup]{
		EnableGetAll: true,
		SearchHandlers: gocrud.BaseSearchHandlers(gocrud.SearchHandlers{
			"like_name": gocrud.KeywordLike("name", nil),
		}),
		WillSave: func(record *model.AccessGroup, context *gin.Context, db *gorm.DB) {
			record.Name = strings.TrimSpace(record.Name)
			if record.Name == "" {
				gocrud.MakeErrorResponse(context, gocrud.RestCoder.BadRequest(), "name cannot be empty")
				return
			}
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func SetupAccessGroupGalleryController(group *gin.RouterGroup, db *gorm.DB) error {
	return gocrud.SetupM2MConnectorController[model.AccessGroupGallery](
		group, db, accessGroupL.New("gallery"),
		"GroupID", "GalleryID",
		nil,
	)
}

func SetupAccessGroupUserController(group *gin.RouterGroup, db *gorm.DB) error {
	return gocrud.SetupM2MConnectorController[model.AccessGroupUser](
		group, db, accessGroupL.New("user"),
		"GroupID", "UserID",
		nil,
	)
}
