package controller

import (
	"net/http"

	"github.com/allape/gocrud"
	"github.com/gin-gonic/gin"
)

func Make401Response(context *gin.Context) {
	gocrud.MakeErrorResponse(context, gocrud.RestCoder.FromStatus(http.StatusUnauthorized), http.StatusText(http.StatusUnauthorized))
}

func Make403Response(context *gin.Context) {
	gocrud.MakeErrorResponse(context, gocrud.RestCoder.FromStatus(http.StatusForbidden), http.StatusText(http.StatusForbidden))
}

func NoCache(context *gin.Context) {
	context.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	context.Header("Pragma", "no-cache")
	context.Header("Expires", "0")
}
