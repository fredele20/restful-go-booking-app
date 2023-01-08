package routes

import (
	"booking-app/core"

	"github.com/gin-gonic/gin"
)


func EventRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/events/create/", core.PublishEvent())
}