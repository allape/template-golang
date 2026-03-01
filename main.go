package main

import (
	"net/http"
	"time"

	"github.com/allape/gocrud"
	"github.com/allape/gogger"
	"github.com/allape/golang/asset"
	"github.com/allape/golang/controller"
	"github.com/allape/golang/env"
	"github.com/allape/golang/model"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var l = gogger.New("main")

func main() {
	err := gogger.InitFromEnv()
	if err != nil {
		l.Error().Fatalf("Failed to init logger: %v", err)
	}

	db, err := gorm.Open(mysql.Open(env.DatabaseDSN), &gorm.Config{
		Logger: logger.New(gogger.New("db").Debug(), logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      logger.Info,
			Colorful:      true,
		}),
	})
	if err != nil {
		l.Error().Fatalf("Failed to open database: %v", err)
	}

	err = db.AutoMigrate(
		&model.User{},
	)
	if err != nil {
		l.Error().Fatalf("Failed to auto migrate database: %v", err)
	}

	engine := gin.Default()

	if env.EnableCors {
		engine.Use(gocrud.NewCors())
	}

	apiGrp := engine.Group("/api")

	err = gocrud.NewSingleHTMLServe(engine.Group("/ui"), env.UIFolder, &gocrud.SingleHTMLServeConfig{
		AllowReplace: true,
	})
	if err != nil {
		l.Error().Fatalf("Failed to setup single html serve: %v", err)
	}

	err = gocrud.NewHttpFileSystem(engine.Group("/static"), env.StaticFolder, &gocrud.HttpFileSystemConfig{
		AllowOverwrite: false,
		AllowUpload:    true,
		EnableDigest:   true,
	})

	err = controller.SetupUserController(apiGrp.Group("/user"), db)
	if err != nil {
		l.Error().Fatalf("Failed to setup user controller: %v", err)
	}

	engine.GET("/", func(context *gin.Context) {
		context.Redirect(http.StatusMovedPermanently, "/ui/")
	})
	engine.GET("/favicon.ico", func(context *gin.Context) {
		context.Data(http.StatusOK, asset.FaviconMIME, asset.Favicon)
	})

	go func() {
		err := engine.Run(env.BindAddr)
		if err != nil {
			l.Error().Fatalf("Failed to start http server: %v", err)
		}
	}()

	gogger.New("ctrl-c").Info().Println("Exiting with", gocrud.Wait4CtrlC())
}
