package controller

import (
	"encoding/json"
	"path"
	"slices"
	"strings"

	"github.com/allape/gocrud"
	"github.com/allape/golang/model"
	"github.com/allape/golang/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var itemL = l.New("item")

func SetupItemController(group *gin.RouterGroup, db *gorm.DB) error {
	err := gocrud.Setup(group, db, itemL.New("crud"), &gocrud.Crud[model.Item]{
		EnableGetAll: true,
		SearchHandlers: gocrud.BaseSearchHandlers(gocrud.SearchHandlers{
			"like_name": gocrud.KeywordLike("name", nil),
			"createBy":  gocrud.KeywordEqual("created_by", nil),
			"in_galleryId": func(db *gorm.DB, values []string, _ *gin.Context) (*gorm.DB, error) {
				if value, ok := gocrud.PickFirstValuableString(values); ok {
					ids := gocrud.IDsFromCommaSeparatedString(value)
					if len(ids) == 0 {
						return db, gocrud.NotArrayError
					}
					return db.Where("`id` IN (SELECT gi.item_id FROM gallery_items gi WHERE gi.gallery_id IN ?)", ids), nil
				}
				return db, nil
			},
		}),
		WillGetAll: func(context *gin.Context, db *gorm.DB) *gorm.DB {
			handledSearch := gocrud.GetHandledSearch(context)
			if !slices.Contains(handledSearch, "in_id") {
				gocrud.MakeErrorResponse(context, gocrud.RestCoder.BadRequest(), "in_id can NOT be empty for getting all")
				return db
			}
			return db
		},
		WillSave: func(record *model.Item, context *gin.Context, db *gorm.DB) {
			record.Name = strings.TrimSpace(record.Name)
		},
	})
	if err != nil {
		return err
	}

	group.POST("/upload", func(context *gin.Context) {
		mpForm, err := context.MultipartForm()
		if err != nil {
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.BadRequest(), err)
			return
		}

		mpRecord, ok := mpForm.Value["record"]
		if !ok || len(mpRecord) == 0 {
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.BadRequest(), "record data not found")
			return
		}

		var record model.Item
		err = json.Unmarshal([]byte(mpRecord[0]), &record)
		if err != nil {
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.BadRequest(), "record data is invalid")
			return
		}

		if record.ID == 0 {
			mpFiles, ok := mpForm.File["file"]
			if !ok || len(mpFiles) == 0 {
				gocrud.MakeErrorResponse(context, gocrud.RestCoder.BadRequest(), "no file in form")
				return
			}
			mpFile := mpFiles[0]

			file, err := mpFile.Open()
			if err != nil {
				gocrud.MakeErrorResponse(context, gocrud.RestCoder.BadRequest(), "mp file is invalid")
				return
			}
			defer func() {
				_ = file.Close()
			}()

			filenameSrc, filenameThu, err := service.SaveFile(file, mpFile.Size, path.Ext(mpFile.Filename))
			if err != nil {
				gocrud.MakeErrorResponse(context, gocrud.RestCoder.BadRequest(), err)
				return
			}

			record.Src = filenameSrc
			record.Thumbnail = filenameThu
		} else {
			var old model.Item
			if err := db.Model(&old).Where("id = ?", record.ID).First(&old).Error; err != nil {
				gocrud.MakeErrorResponse(context, gocrud.RestCoder.BadRequest(), "record not found")
				return
			}

			record.Thumbnail = old.Thumbnail
			record.Src = old.Src
		}

		if err := db.Save(&record).Error; err != nil {
			itemL.Error().Printf("record save error: %s", err)
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.BadRequest(), "failed to save record [error]")
			return
		}

		gocrud.MakeOkayDataResponse(context, record)
	})

	return nil
}

func SetupItemTagController(group *gin.RouterGroup, db *gorm.DB) error {
	return gocrud.SetupM2MConnectorController[model.ItemTag](
		group, db, itemL.New("tag"),
		"ItemID", "TagID",
		nil,
	)
}
