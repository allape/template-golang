package main

import (
	"fmt"
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
	err := SetupLogger()
	if err != nil {
		l.Error().Fatalf("failed to init logger: %v", err)
		return
	}

	db, err := SetupDatabase()
	if err != nil {
		l.Error().Fatalf("failed to setup database: %v", err)
		return
	}

	engine, err := SetupControllers(db)
	if err != nil {
		l.Error().Fatalf("failed to setup controllers: %v", err)
		return
	}

	go func() {
		err := engine.Run(env.BindAddr)
		if err != nil {
			l.Error().Fatalf("failed to start http server: %v", err)
		}
	}()

	gogger.New("ctrl-c").Info().Println("exiting with", gocrud.Wait4CtrlC())
}

func SetupControllers(db *gorm.DB) (*gin.Engine, error) {
	engine := gin.Default()

	if env.DebugMode {
		engine.Use(gocrud.NewCors())
	}

	err := gocrud.NewSingleHTMLServe(engine.Group("/ui"), env.UIFolder, &gocrud.SingleHTMLServeConfig{
		AllowReplace: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to setup single html serve: %v", err)
	}

	err = gocrud.NewHttpFileSystem(engine.Group("/static"), env.StaticFolder, &gocrud.HttpFileSystemConfig{
		AllowOverwrite: false,
		AllowUpload:    true,
		EnableDigest:   true,
	})

	engine.GET("/", func(context *gin.Context) {
		context.Redirect(http.StatusMovedPermanently, "/ui/")
	})
	engine.GET("/favicon.ico", func(context *gin.Context) {
		context.Data(http.StatusOK, asset.FaviconMIME, asset.Favicon)
	})

	apiGrp := engine.Group("/api")

	err = controller.SetupTagController(apiGrp.Group("/tag"), db)
	if err != nil {
		return nil, fmt.Errorf("failed to setup tag controller: %v", err)
	}

	err = controller.SetupItemController(apiGrp.Group("/item"), db)
	if err != nil {
		return nil, fmt.Errorf("failed to setup item controller: %v", err)
	}

	err = controller.SetupItemTagController(apiGrp.Group("/item-tag"), db)
	if err != nil {
		return nil, fmt.Errorf("failed to setup item tag controller: %v", err)
	}

	err = controller.SetupGalleryController(apiGrp.Group("/gallery"), db)
	if err != nil {
		return nil, fmt.Errorf("failed to setup gallery controller: %v", err)
	}

	err = controller.SetupGalleryItemController(apiGrp.Group("/gallery-item"), db)
	if err != nil {
		return nil, fmt.Errorf("failed to setup gallery item controller: %v", err)
	}

	return engine, nil
}

func SetupLogger() error {
	err := gogger.InitFromEnv()
	if err != nil {
		return err
	}

	if env.DebugMode {
		gogger.Level = gogger.Debug
	}

	return nil
}

func SetupDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(env.DatabaseDSN), &gorm.Config{
		Logger: logger.New(gogger.New("db").Debug(), logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      logger.Info,
			Colorful:      true,
		}),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(
		&model.Tag{},
		&model.Item{},
		&model.ItemTag{},
		&model.Gallery{},
		&model.GalleryItem{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate models: %v", err)
	}

	// MyISAM tables in MySQL
	//err = db.Set("gorm:table_options", "ENGINE=MyISAM CHARSET=utf8mb4").AutoMigrate(
	//	&model.AccessLog{},
	//)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to auto migrate MyISAM tables: %v", err)
	//}

	return db, nil
}
