package routes

import (
	"booking-app/core"

	"github.com/gin-gonic/gin"
)

func UserAuthentication(incomingRoutes *gin.Engine) {
	incomingRoutes.POST(BASEURI+"/users/register", core.Register())
	incomingRoutes.POST(BASEURI+"/users/login", core.Login())
}
