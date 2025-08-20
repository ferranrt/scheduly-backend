package controllers

import (
	"net/http"

	"bifur.app/core/cmd/rest/helpers"
	"bifur.app/core/internal/ports"
	"github.com/gin-gonic/gin"
)

func GetCentersController(ctx *gin.Context, centersRepository ports.CentersRepository) {
	userCtx, err := helpers.GetUserIdFromRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	centers, err := centersRepository.GetAll(ctx.Request.Context(), userCtx.AsUUID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, centers)
}
