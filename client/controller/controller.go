package controller

import (
	"net/http"

	"github.com/allape/gocrud"
	"github.com/gin-gonic/gin"
)

func Make401Response(context *gin.Context) {
	gocrud.MakeErrorResponse(context, gocrud.RestCoder.FromStatus(http.StatusUnauthorized), http.StatusText(http.StatusUnauthorized))
}
