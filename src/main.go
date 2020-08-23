package main

import (
    "errors"
    "github.com/gin-gonic/gin"
    "kartiz/api"
    "kartiz/utils"
    "log"
    "net/http"
)

func getServer() *gin.Engine {
    r := gin.New()

    // Set middleware
    r.Use(gin.Logger())
    r.Use(utils.Recovery())

    // Set no route handler
    r.NoRoute(func(c *gin.Context) {
        c.JSON(http.StatusNotFound, utils.BuildJsonOutput(nil, errors.New("resource not found"), http.StatusNotFound))
    })

    // Set common routes
    r.GET("/", func(c *gin.Context) {c.Status(http.StatusOK)})

    // Set api routes
    api.Apply(r)
    return r
}

func main () {
    log.Fatal(getServer().Run())
}