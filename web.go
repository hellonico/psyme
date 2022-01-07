package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var dbmap = initDb()

func web() {

	gin.ForceConsoleColor()

	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("psyme", store))

	router.LoadHTMLGlob("templates/*")

	router.GET("/", GetRoute)
	router.GET("/erase", GetErase)
	router.GET("/presubmit", GetPresubmit)
	router.GET("/compare/:username1/:username2", GetCompare)
	router.GET("/users", GetUsers)
	router.POST("/submit", PostSubmit)
	router.GET("/summary", GetSummary)
	router.GET("/one/:user", GetOneUser)
	router.GET("/a/:id", GetOneUser)
	router.GET("/c/:id/:choice", GetChoice)

	router.Run(":8080")
}
