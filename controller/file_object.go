package controller

import (
	"github.com/allape/gocrud"
	"github.com/allape/gogger"
	"github.com/allape/golang/env"
	"github.com/allape/golang/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupFileObjectController(group *gin.RouterGroup, db *gorm.DB) error {
	return gocrud.NewHttpFileSystemObjectController[model.FileObject](
		group, db, gogger.New("hfso:static"),
		env.StaticFolder, &gocrud.HttpFileSystemConfig{
			AllowUpload:   true,
			FileMasterKey: gocrud.SHASum256FromString(env.StaticFileMasterKey),
		},
		"",
	)
}
