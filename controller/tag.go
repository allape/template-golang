package controller

import (
	"strings"

	"github.com/allape/gocrud"
	"github.com/allape/golang/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var tagL = l.New("tag")

func SetupTagController(group *gin.RouterGroup, db *gorm.DB) error {
	err := gocrud.Setup(group, db, tagL.New("crud"), &gocrud.Crud[model.Tag]{
		EnableGetAll: true,
		SearchHandlers: gocrud.BaseSearchHandlers(gocrud.SearchHandlers{
			"like_name":  gocrud.KeywordLike("name", nil),
			"like_alias": gocrud.KeywordLike("alias", nil),
			"like_keyword": func(db *gorm.DB, values []string, _ *gin.Context) (*gorm.DB, error) {
				if value, ok := gocrud.PickFirstValuableString(values); ok {
					v := "%" + value + "%"
					return db.Where("(name LIKE ? OR alias LIKE ?)", v, v), nil
				}
				return db, nil
			},
		}),
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
