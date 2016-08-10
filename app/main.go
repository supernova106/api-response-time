package main

import "github.com/gin-gonic/gin"
import "net/http"

func main() {
    // Creates a gin router with default middleware:
    // logger and recovery (crash-free) middleware
    router := gin.Default()
    
    router.GET("/api", func(c *gin.Context) {
        //firstname := c.DefaultQuery("firstname", "Guest")
        //lastname := c.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")
        firstname := "bean"
        lastname := "nguyen"
        c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
    })

    // This handler will match /user/john but will not match neither /user/ or /user
    router.GET("/api/:apiname", func(c *gin.Context) {
        apiname := c.Param("apiname")
        c.String(http.StatusOK, "Hello %s", apiname)
    })
    
    // By default it serves on :8080 unless a
    // PORT environment variable was defined.
    router.Run()
    // router.Run(":3000") for a hard coded port

}