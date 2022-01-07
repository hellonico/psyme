package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func GetErase(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("results", nil)
	session.Set("user", nil)
	session.Save()
	c.Redirect(http.StatusSeeOther, "/summary")
}

func GetRoute(c *gin.Context) {
	dbmap := initDb()
	defer dbmap.Db.Close()

	var articles []Article
	_, err := dbmap.Select(&articles, "SELECT * FROM Article ORDER BY MyIndex")
	if err != nil {
		fmt.Printf("Erro %s \n", err)
	}

	// not used for now
	session := sessions.Default(c)
	m := getMapFromSessioN(session)

	c.HTML(http.StatusOK, "index.tmpl", gin.H{"as": articles, "results": m})
}
func GetPresubmit(c *gin.Context) {
	session := sessions.Default(c)
	possibleCurrentUser := getUserFromSession(session)

	if possibleCurrentUser != "" {
		// user already in session skip form
		persistResults(session, dbmap, possibleCurrentUser)
		c.Redirect(http.StatusSeeOther, "/users")
	} else {
		// send to user form
		c.HTML(http.StatusOK, "presend.tmpl", nil)
	}
}
func GetCompare(c *gin.Context) {
	username1 := c.Param("username1")
	username2 := c.Param("username2")

	user1 := getUserFromName(dbmap, username1)
	answers1 := stringToAnswersMap(user1.Answers)

	user2 := getUserFromName(dbmap, username2)
	answers2 := stringToAnswersMap(user2.Answers)

	var articles []Article
	dbmap.Select(&articles, "SELECT * FROM Article ORDER BY MyIndex")

	var results = make([]CompareUser, 0)
	matches := 0
	for _, a := range articles {
		same := answers1[a.Id] == answers2[a.Id] && answers1[a.Id] != ""
		results = append(results, CompareUser{a, same})
		if same {
			matches++
		}
	}

	matchScore := (float64(matches) / (float64(len(articles)))) * 100
	// for testing
	// matchScore += 50

	// ridiculous lol
	// no int range in templates
	numberHearts := make([]int, int(matchScore/10))

	//TODO to rewrite when go 1.18 can be installed (support for generics)
	//filteredItems, _ := Filter(list, func(item int) bool {
	//	return item%2 == 0
	//})
	fmt.Printf("~~ %s %s\n", username1, username2)

	c.HTML(http.StatusOK, "compare.tmpl", gin.H{"numberHearts": numberHearts, "matchScore": matchScore, "username1": username1, "username2": username2, "results": results})

}

func GetUsers(c *gin.Context) {
	session := sessions.Default(c)
	currentName := getUserFromSession(session)

	if currentName == "" {
		c.Redirect(http.StatusSeeOther, "/presubmit")
	} else {
		mainUser := getUserFromName(dbmap, currentName)
		scores := getScores(dbmap, currentName)
		c.HTML(http.StatusOK, "users.tmpl", gin.H{"current": currentName, "scores": scores, "currentTotal": len(getUserAnswers(mainUser))})
	}
}
func PostSubmit(c *gin.Context) {
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
}

func GetSummary(c *gin.Context) {
	session := sessions.Default(c)
	m := getMapFromSessioN(session)

	results, progress := getResultsFromMap(dbmap, m)

	c.HTML(http.StatusOK, "summary.tmpl", gin.H{"results": results, "progress": int(progress)})
}

func GetOneUser(c *gin.Context) {
	userName := c.Param("user")
	user := getUserFromName(dbmap, userName)
	m := getUserAnswers(user)

	results, progress := getResultsFromMap(dbmap, m)

	c.HTML(http.StatusOK, "one.tmpl", gin.H{"userName": userName, "results": results, "progress": int(progress)})
}

func GetArticles(c *gin.Context) {
	id := c.Param("id")
	a, _ := dbmap.Get(Article{}, id)

	session := sessions.Default(c)
	m := getMapFromSessioN(session)
	fmt.Printf("Selected: %s\n", m[id])

	c.HTML(http.StatusOK, "article.tmpl", gin.H{"a": a, "selected": m[id]})
}
func GetChoice(c *gin.Context) {

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

	// if user is defined, persist
	possibleCurrentUser := getUserFromSession(session)
	if possibleCurrentUser != "" {
		persistResults(session, dbmap, possibleCurrentUser)
	}

	a, _ := dbmap.Get(Article{}, id)
	c.HTML(http.StatusOK, "choice.tmpl", gin.H{"a": a, "choice": choice, "next": next})
}
