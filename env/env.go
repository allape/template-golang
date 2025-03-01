package env

import (
	"github.com/allape/goenv"
)

const Project = "GOLANG_"

const (
	bindAddr     = Project + "BIND_ADDR"
	enableCors   = Project + "ENABLE_CORS"
	databaseDSN  = Project + "DATABASE_DSN"
	uiFolder     = Project + "UI_FOLDER"
	staticFolder = Project + "STATIC_FOLDER"
)

var (
	BindAddr     = goenv.Getenv(bindAddr, ":8080")
	EnableCors   = goenv.Getenv(enableCors, true)
	DatabaseDSN  = goenv.Getenv(databaseDSN, "root:Root_123456@tcp(127.0.0.1:3306)/golang?charset=utf8mb4&parseTime=True&loc=Local")
	UIFolder     = goenv.Getenv(uiFolder, "./ui/dist/index.html")
	StaticFolder = goenv.Getenv(staticFolder, "./static")
)
