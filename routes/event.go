package routes

import (
	"booking-app/core"

	"github.com/gin-gonic/gin"
)


var BASEURI = "/api/booking"

func EventRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST(BASEURI +"/events/create/", core.PublishEvent())
}