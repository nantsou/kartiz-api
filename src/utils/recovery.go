package utils

import (
    "errors"
    "fmt"
    "github.com/gin-gonic/gin"
    "net/http"
)

func Recovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                c.JSON(http.StatusInternalServerError, BuildJsonOutput(
                    nil,
                    errors.New(fmt.Sprintf("%s", err)),
                    http.StatusInternalServerError))
            }
        }()
        c.Next()
    }
}
