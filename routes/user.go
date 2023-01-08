package routes

import (
	"booking-app/core"

	"github.com/gin-gonic/gin"
)



func UserAuthentication(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/register", core.Register())
	incomingRoutes.POST("/users/login", core.Login())
}