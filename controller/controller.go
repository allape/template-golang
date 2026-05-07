package controller

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/allape/gocrud"
	"github.com/allape/gogger"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	DefaultPageSize = 100
)

func willGetAllInID(context *gin.Context, db *gorm.DB) *gorm.DB {
	ids := gocrud.IDsFromCommaSeparatedString(context.Query("in_id"))
	if len(ids) == 0 {
		gocrud.MakeErrorResponse(context, gocrud.RestCoder.BadRequest(), "in_id cannot be empty")
		return nil
	} else if len(ids) > DefaultPageSize {
		gocrud.MakeErrorResponse(context, gocrud.RestCoder.BadRequest(), "in_id too long")
		return nil
	}
	return db
}

func SetupDualPrimaryKeyModelController[T any](
	group *gin.RouterGroup, db *gorm.DB,
	objectFieldName1, objectFieldName2 string,
	jsonFieldName1, jsonFieldName2 string,
	databaseFieldName1, databaseFieldName2 string,
	logger *gogger.Logger, maxCount int,
) error {
	if objectFieldName1 == "" || objectFieldName2 == "" {
		return fmt.Errorf("field1 and field2 cannot be empty")
	}

	inFieldName1 := "in_" + jsonFieldName1
	inFieldName2 := "in_" + jsonFieldName2

	err := gocrud.New(group, db, gocrud.Crud[T]{
		DisableSave:   true,
		DisableCount:  true,
		DisableDelete: true,
		DisableGetOne: true,
		DisablePage:   true,
		EnableGetAll:  true,
		WillGetAll: func(context *gin.Context, db *gorm.DB) *gorm.DB {
			f1 := gocrud.IDsFromCommaSeparatedString(context.Query(inFieldName1))
			f2 := gocrud.IDsFromCommaSeparatedString(context.Query(inFieldName2))
			if len(f1) == 0 && len(f2) == 0 {
				gocrud.MakeErrorResponse(context, gocrud.RestCoder.BadRequest(), "ids cannot be empty")
				return nil
			}
			if len(f1) > maxCount || len(f2) > maxCount {
				gocrud.MakeErrorResponse(context, gocrud.RestCoder.BadRequest(), "ids too long")
				return nil
			}
			return db
		},
		SearchHandlers: map[string]gocrud.SearchHandler{
			inFieldName1: gocrud.KeywordIDIn(databaseFieldName1, gocrud.OverflowedArrayTrimmerFilter[gocrud.ID](maxCount)),
			inFieldName2: gocrud.KeywordIDIn(databaseFieldName2, gocrud.OverflowedArrayTrimmerFilter[gocrud.ID](maxCount)),
		},
	})
	if err != nil {
		return err
	}

	// /save/[jsonField1 || jsonField2]?[jsonField1]=1,2,3...&[jsonField2]=1,2,3...
	group.POST("/save/:primaryFieldName", func(c *gin.Context) {
		primaryFieldName := strings.TrimSpace(c.Param("primaryFieldName"))
		if primaryFieldName != jsonFieldName1 && primaryFieldName != jsonFieldName2 {
			gocrud.MakeErrorResponse(c, gocrud.RestCoder.BadRequest(), "field name invalid")
			return
		}

		field1Ids := gocrud.IDsFromCommaSeparatedString(c.Query(jsonFieldName1))
		field2Ids := gocrud.IDsFromCommaSeparatedString(c.Query(jsonFieldName2))

		var pid gocrud.ID
		var sids []gocrud.ID

		var objectPrimaryFieldName string
		var objectSecondaryFieldName string
		var dbFieldName string

		switch primaryFieldName {
		case jsonFieldName1:
			pid = gocrud.Pick(field1Ids, 0, 0)
			sids = field2Ids
			objectPrimaryFieldName = objectFieldName1
			objectSecondaryFieldName = objectFieldName2
			dbFieldName = databaseFieldName1
		case jsonFieldName2:
			pid = gocrud.Pick(field2Ids, 0, 0)
			sids = field1Ids
			objectPrimaryFieldName = objectFieldName2
			objectSecondaryFieldName = objectFieldName1
			dbFieldName = databaseFieldName2
		}

		if pid == 0 {
			gocrud.MakeErrorResponse(c, gocrud.RestCoder.BadRequest(), "pid cannot be empty")
			return
		}

		sids = gocrud.RemoveDuplication(sids)

		count := int64(0)

		err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Delete(new(T), dbFieldName+" = ?", pid).Error; err != nil {
				return err
			}

			if len(sids) == 0 {
				return nil
			}

			items := make([]*T, len(sids))
			for i, sid := range sids {
				record := new(T)

				reflected := reflect.ValueOf(record).Elem()

				primaryField := reflected.FieldByName(objectPrimaryFieldName)
				primaryField.SetUint(uint64(pid))

				secondaryField := reflected.FieldByName(objectSecondaryFieldName)
				secondaryField.SetUint(uint64(sid))

				items[i] = record
			}

			res := tx.Save(items)
			if res.Error != nil {
				return res.Error
			}

			count = res.RowsAffected

			return nil
		})
		if err != nil {
			logger.Error().Printf("failed to save %v for %d: %v", sids, pid, err)
			gocrud.MakeErrorResponse(c, gocrud.RestCoder.InternalServerError(), "failed to save [error]")
			return
		}

		gocrud.MakeOkayDataResponse(c, count)
	})

	return nil
}
