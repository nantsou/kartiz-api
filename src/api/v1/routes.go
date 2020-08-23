package v1

import (
    "github.com/gin-gonic/gin"
)

func Apply(r *gin.RouterGroup) {
    v1 := r.Group("/v1")
    v1.GET("/", version)
}
