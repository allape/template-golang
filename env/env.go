package env

import (
	"github.com/allape/goenv"
)

const Project = "projectname_"

const (
	debugMode = Project + "DEBUG_MODE"

	bindAddr = Project + "BIND_ADDR"

	databaseDSN = Project + "DATABASE_DSN"

	uiFolder = Project + "UI_FOLDER"

	staticFolder = Project + "STATIC_FOLDER"
)

var (
	DebugMode = goenv.Getenv(debugMode, true)

	BindAddr = goenv.Getenv(bindAddr, ":8080")

	DatabaseDSN = goenv.Getenv(databaseDSN, "root:Root_123456@tcp(127.0.0.1:3306)/projectname?charset=utf8mb4&parseTime=True&loc=Local")

	UIFolder = goenv.Getenv(uiFolder, "./ui/dist/index.html")

	StaticFolder = goenv.Getenv(staticFolder, "./static")
)
