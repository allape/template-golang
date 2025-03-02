package controller

import (
	"github.com/allape/gocrud"
	"github.com/allape/golang/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
)

func SetupUserController(group *gin.RouterGroup, db *gorm.DB) error {
	err := gocrud.New(group, db, gocrud.Crud[model.User]{
		DefaultPageSize: DefaultPageSize,
		EnableGetAll:    true,
		SearchHandlers: map[string]gocrud.SearchHandler{
			"like_name": gocrud.KeywordLike("name", nil),
			"in_id":     gocrud.KeywordIDIn("id", gocrud.OverflowedArrayTrimmerFilter[gocrud.ID](DefaultPageSize)),
			"deleted":   gocrud.NewSoftDeleteSearchHandler(""),
		},
		OnDelete: gocrud.NewSoftDeleteHandler[model.User](gocrud.RestCoder),
		WillSave: func(record *model.User, context *gin.Context, db *gorm.DB) {
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
