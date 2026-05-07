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

	clientBindAddr      = Project + "CLIENT_BIND_ADDR"
	clientUIFolder      = Project + "CLIENT_UI_FOLDER"
	clientDebugMode     = Project + "CLIENT_DEBUG_MODE"
	clientDebugUserID   = Project + "CLIENT_DEBUG_USER_ID"
	clientDebugUserName = Project + "CLIENT_DEBUG_USER_NAME"
)

var (
	DebugMode = goenv.Getenv(debugMode, true) // FIXME remove debug mode when production

	DatabaseDSN = goenv.Getenv(databaseDSN, "root:Root_123456@tcp(127.0.0.1:3306)/projectname?charset=utf8mb4&parseTime=True&loc=Local")

	BindAddr     = goenv.Getenv(bindAddr, ":8080")
	UIFolder     = goenv.Getenv(uiFolder, "./ui/dist/index.html")
	StaticFolder = goenv.Getenv(staticFolder, "./static")

	ClientBindAddr      = goenv.Getenv(clientBindAddr, ":8888")
	ClientUIFolder      = goenv.Getenv(clientUIFolder, "./client/ui/dist/index.html")
	ClientDebugMode     = goenv.Getenv(clientDebugMode, false)
	ClientDebugUserID   = goenv.Getenv(clientDebugUserID, "-1")
	ClientDebugUserName = goenv.Getenv(clientDebugUserName, "allape")
)
