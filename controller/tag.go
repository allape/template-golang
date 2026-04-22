package controller

import (
	"strings"

	"github.com/allape/gocrud"
	"github.com/allape/golang/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupTagController(group *gin.RouterGroup, db *gorm.DB) error {
	err := gocrud.New(group, db, gocrud.Crud[model.Tag]{
		DefaultPageSize: DefaultPageSize,
		SearchHandlers: map[string]gocrud.SearchHandler{
			"in_id":             gocrud.KeywordIDIn("id", gocrud.OverflowedArrayTrimmerFilter[gocrud.ID](DefaultPageSize)),
			"like_name":         gocrud.KeywordLike("name", nil),
			"like_alias":        gocrud.KeywordLike("alias", nil),
			"deleted":           gocrud.NewSoftDeleteSearchHandler(""),
			"orderBy_priority":  gocrud.SortBy("priority"),
			"orderBy_createdAt": gocrud.SortBy("created_at"),
			"orderBy_updatedAt": gocrud.SortBy("updated_at"),
		},
		OnDelete: gocrud.NewSoftDeleteHandler[model.Tag](gocrud.RestCoder),
		WillSave: func(record *model.Tag, context *gin.Context, db *gorm.DB) {
			record.Name = strings.TrimSpace(record.Name)
			if record.Name == "" {
				gocrud.MakeErrorResponse(context, gocrud.RestCoder.BadRequest(), "name cannot be empty")
				return
			}

			record.Alias = strings.TrimSpace(record.Alias)
		},
	})
	if err != nil {
		return err
	}

	return nil
}
