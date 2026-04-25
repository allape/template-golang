package client

import (
	"fmt"
	"net/http"

	"github.com/allape/gocrud"
	"github.com/allape/golang/asset"
	"github.com/allape/golang/client/controller"
	"github.com/allape/golang/env"
	"github.com/allape/gophorward"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupControllers(db *gorm.DB) (*gin.Engine, error) {
	engine := gin.Default()

	if env.DebugMode {
		engine.Use(gocrud.NewCors())
	}

	engine.Use(gocrud.RecoveryHandler(env.DebugMode))

	engine.Use(gophorward.GinMiddlewareHandler(func(context *gin.Context) {
		controller.Make401Response(context)
	}))

	err := gocrud.NewSingleHTMLServe(engine.Group("/ui"), env.UIFolder, &gocrud.SingleHTMLServeConfig{
		AllowReplace: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to setup single html serve: %v", err)
	}

	engine.GET("/", func(context *gin.Context) {
		context.Redirect(http.StatusMovedPermanently, "/ui/")
	})
	engine.GET("/favicon.ico", func(context *gin.Context) {
		context.Data(http.StatusOK, asset.FaviconMIME, asset.Favicon)
	})

	apiGrp := engine.Group("/api")

	err = controller.SetupGalleryController(apiGrp.Group("/gallery"), db)
	if err != nil {
		return nil, fmt.Errorf("failed to setup gallery controller: %v", err)
	}

	return engine, nil
}
