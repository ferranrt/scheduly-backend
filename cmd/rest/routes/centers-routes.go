package routes

import (
	"bifur.app/core/cmd/rest/controllers"
	"bifur.app/core/internal/ports"
	"github.com/gin-gonic/gin"
)

type CentersRoutesDeps struct {
	CentersRepository ports.CentersRepository
}

func SetupCentersRoutes(router *gin.RouterGroup, deps *CentersRoutesDeps) {

	router.POST("/initial-setup", func(ctx *gin.Context) {

		controllers.GetCentersController(ctx, deps.CentersRepository)
	})
}
