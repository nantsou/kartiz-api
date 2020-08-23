package api

import (
    "github.com/gin-gonic/gin"
    v1 "kartiz/api/v1"
)

func Apply(r *gin.Engine) {
    api := r.Group("/api")
    v1.Apply(api)
}