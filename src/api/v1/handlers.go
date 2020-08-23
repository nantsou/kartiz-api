package v1

import (
    "github.com/gin-gonic/gin"
    "kartiz/utils"
    "net/http"
)

func version(c *gin.Context) {
    c.JSON(http.StatusOK, utils.BuildJsonOutput("v1", nil, http.StatusOK))
}