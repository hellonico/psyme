package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gopkg.in/gorp.v1"
	"net/http"
	"sort"
	"strings"
	"time"
)

type SampleRequest struct {
	Name string `form:"nameField"`
}

func search(slice []Article, item string) int {
	for i := range slice {
		if slice[i].Id == item {
			return i
		}
	}
	return -1
}

func getResultsFromSessionAsJson(session sessions.Session) string {
	var jsonStr = fmt.Sprintf("%s", session.Get("results"))
	return jsonStr
}

func getMapFromSessioN(session sessions.Session) map[string]string {
	jsonStr := getResultsFromSessionAsJson(session)

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
type Score struct {
	Name  string
	Score int
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
		session.Set("results", nil)
		session.Set("user", nil)
		session.Save()
		c.Redirect(http.StatusSeeOther, "/summary")
	})

	router.GET("/presubmit", func(c *gin.Context) {
		session := sessions.Default(c)
		possibleCurrentUser := session.Get("user")
		if possibleCurrentUser != nil {
			// user already in session skip form
			name := fmt.Sprintf("%s", possibleCurrentUser)
			fmt.Printf("User %s is already logged in", name)
			persistResults(session, dbmap, name)
			c.Redirect(http.StatusSeeOther, "/users")
		} else {
			// send to user form
			c.HTML(http.StatusOK, "presend.tmpl", nil)
		}

	})

	router.GET("/users", func(c *gin.Context) {

		session := sessions.Default(c)
		possibleCurrentUser := session.Get("user")

		if possibleCurrentUser == nil {
			c.Redirect(http.StatusSeeOther, "/presubmit")
		}
		currentName := possibleCurrentUser
		scores := getScores(dbmap, currentName)

		c.HTML(http.StatusOK, "users.tmpl", gin.H{"current": currentName, "scores": scores})
	})
	router.POST("/submit", func(c *gin.Context) {
		//var request SampleRequest
		//err := c.Bind(&request)
		//if err != nil {
		//	fmt.Printf("err: %s", err)
		//}
		session := sessions.Default(c)

		// update user in session
		name := strings.Trim(c.PostForm("message"), " ")
		session.Set("user", name)
		session.Save()

		persistResults(session, dbmap, name)

		c.Redirect(http.StatusSeeOther, "/users")
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
		progress = (float64(len(results)) - notanswer) / float64(len(results)) * 100
		fmt.Printf("%.9f , %d , %d", progress, len(results), notanswer)

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

func getScores(dbmap *gorp.DbMap, currentName interface{}) []Score {
	var currentUser User
	dbmap.SelectOne(&currentUser, fmt.Sprintf("SELECT * FROM User WHERE Name='%s'", currentName))

	var others []User
	dbmap.Select(&others, fmt.Sprintf("SELECT * FROM User WHERE Name!='%s'", currentName))

	var scores = make([]Score, len(others))

	for i, other := range others {
		scores[i] = Score{other.Name, compareUsers(dbmap, currentUser, other)}
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})
	return scores
}

func persistResults(session sessions.Session, dbmap *gorp.DbMap, name string) {
	answers := getResultsFromSessionAsJson(session)

	obj, _ := dbmap.Get(User{}, name)
	if obj == nil {
		dbmap.Insert(&User{name, answers})
	} else {
		dbmap.Update(&User{name, answers})
	}
}
