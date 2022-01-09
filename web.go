package main

import (
	"bytes"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"image/jpeg"
	"net/http"
)

var dbmap = initDb()

func web() {

	gin.ForceConsoleColor()

	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("psyme", store))

	router.LoadHTMLGlob("templates/*")
	router.Static("/assets", "./assets")

	router.GET("/", GetRoute)
	router.GET("/erase", GetErase)
	router.GET("/presubmit", GetPresubmit)
	router.GET("/compare/:username1/:username2", GetCompare)
	router.GET("/users", GetUsers)
	router.POST("/submit", PostSubmit)
	router.GET("/summary", GetSummary)
	router.GET("/one/:user", GetOneUser)
	router.GET("/a/:id", GetArticle)
	router.GET("/c/:id/:choice", GetChoice)
	router.GET("/image/:user", GetImage)

	router.Run(":8080")
}

func GetArticle(c *gin.Context) {
	id := c.Param("id")
	a, _ := dbmap.Get(Article{}, id)

	session := sessions.Default(c)
	m := getMapFromSessioN(session)
	fmt.Printf("Selected: %s\n", m[id])

	c.HTML(http.StatusOK, "article.tmpl", gin.H{"a": a, "selected": m[id]})
}

func GetImage(context *gin.Context) {
	img := mapForUser("Nico")
	writeImgToFile(img)

	buf := new(bytes.Buffer)
	jpeg.Encode(buf, img, nil)
	context.Data(http.StatusOK, "image/jpeg", buf.Bytes())
}
