package server

import (
    "github.com/gin-gonic/gin"
    "github.com/yuuyahypg/ssolap/server/controllers"
)

func SetRoutes(e *gin.Engine) {
    e.GET("/api/home", controllers.Home)
}
