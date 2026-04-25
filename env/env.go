package env

import (
	"github.com/allape/goenv"
)

const Project = "projectname_"

const (
	debugMode = Project + "DEBUG_MODE"

	databaseDSN = Project + "DATABASE_DSN"

	bindAddr     = Project + "BIND_ADDR"
	uiFolder     = Project + "UI_FOLDER"
	staticFolder = Project + "STATIC_FOLDER"

	bindClientAddr = Project + "BIND_CLIENT_ADDR"
	uiClientFolder = Project + "UI_CLIENT_FOLDER"
)

var (
	DebugMode = goenv.Getenv(debugMode, true)

	DatabaseDSN = goenv.Getenv(databaseDSN, "root:Root_123456@tcp(127.0.0.1:3306)/projectname?charset=utf8mb4&parseTime=True&loc=Local")

	BindAddr     = goenv.Getenv(bindAddr, ":8080")
	UIFolder     = goenv.Getenv(uiFolder, "./ui/dist/index.html")
	StaticFolder = goenv.Getenv(staticFolder, "./static")

	BindClientAddr = goenv.Getenv(bindClientAddr, ":8888")
	UIClientFolder = goenv.Getenv(uiClientFolder, "./uiclient/dist/index.html")
)
