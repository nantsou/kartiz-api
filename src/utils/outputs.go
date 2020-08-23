package utils

import "github.com/gin-gonic/gin"

func BuildJsonOutput(data interface{}, err error, statusCode int) gin.H {
    output := gin.H{"data": data, "statusCode": statusCode, "error": nil}
    if err != nil {
        output["error"] = err.Error()
    }
    return output
}
