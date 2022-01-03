package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func search(slice []Article, item string) int {
	for i := range slice {
		if slice[i].Id == item {
			return i
		}
	}
	return -1
}

func getMapFromSessioN(session sessions.Session) map[string]string {
	var jsonStr = fmt.Sprintf("%s", session.Get("results"))

	if jsonStr != "" {
		x := map[string]string{}
		json.Unmarshal([]byte(jsonStr), &x)
		return x
	} else {
		m := make(map[string]string)
		m["date"] = time.Now().String()
		mJson, _ := json.Marshal(m)
		session.Set("results", string(mJson))
		return m
	}
}
func updateResults(session sessions.Session, id string, choice string) {
	m := getMapFromSessioN(session)
	m[id] = choice
	mJson, _ := json.Marshal(m)
	session.Set("results", string(mJson))
	session.Save()
}

type Result struct {
	Id       string
	Title    string
	ChoiceI  string
	Choice   string
	Answer   string
	ImageURL string
}

func web() {

	gin.ForceConsoleColor()

	dbmap := initDb()

	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("psyme", store))

	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		dbmap := initDb()
		defer dbmap.Db.Close()

		var articles []Article
		_, err := dbmap.Select(&articles, "SELECT * FROM Article ORDER BY MyIndex")
		if err != nil {
			fmt.Printf("Erro %s \n", err)
		}

		session := sessions.Default(c)
		m := getMapFromSessioN(session)
		fmt.Printf("%s", m)

		c.HTML(http.StatusOK, "index.tmpl", gin.H{"as": articles, "results": m})
	})
	router.GET("/erase", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("results", "")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/summary")
	})
	router.GET("/summary", func(c *gin.Context) {
		session := sessions.Default(c)
		m := getMapFromSessioN(session)
		// fmt.Printf("%s", m)
		var articles []Article
		dbmap.Select(&articles, "SELECT * FROM Article ORDER BY MyIndex")

		results := make([]Result, len(articles))
		notanswer := float64(0)
		for i := range articles {
			choiceI := m[articles[i].Id]
			var choice string
			var answer string
			switch choiceI {
			case "1":
				choice = articles[i].Choice1
				answer = articles[i].Answer1.Title
			case "2":
				choice = articles[i].Choice2
				answer = articles[i].Answer2.Title
			case "3":
				choice = articles[i].Choice3
				answer = articles[i].Answer3.Title
			case "4":
				choice = articles[i].Choice4
				answer = articles[i].Answer4.Title
			case "":
				// fmt.Printf("+1")
				notanswer = notanswer + 1
			}
			results[i] = Result{articles[i].Id, articles[i].Title, choiceI, choice, answer, articles[i].ImageURL}
		}

		var progress float64
		progress = (float64(len(articles)) - notanswer) / float64(len(articles)) * 100
		fmt.Printf("%.9f , %d , %d", progress, len(articles), notanswer)

		c.HTML(http.StatusOK, "summary.tmpl", gin.H{"results": results, "progress": int(progress)})
	})
	router.GET("/a/:id", func(c *gin.Context) {
		id := c.Param("id")
		a, _ := dbmap.Get(Article{}, id)

		session := sessions.Default(c)
		m := getMapFromSessioN(session)
		fmt.Printf("Selected: %s\n", m[id])

		c.HTML(http.StatusOK, "article.tmpl", a)
	})
	router.GET("/c/:id/:choice", func(c *gin.Context) {

		var articles []Article
		dbmap.Select(&articles, "Select Id from Article order by MyIndex")

		id := c.Param("id")
		loc := search(articles, id)

		var next string
		if articles[len(articles)-1].Id != id {
			next = articles[loc+1].Id
		} else {
			next = "-1"
		}

		choice := c.Param("choice")
		session := sessions.Default(c)
		updateResults(session, id, choice)

		a, _ := dbmap.Get(Article{}, id)
		c.HTML(http.StatusOK, "choice.tmpl", gin.H{"a": a, "choice": choice, "next": next})
	})

	router.Run(":8080")
}
