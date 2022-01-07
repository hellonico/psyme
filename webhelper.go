package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"gopkg.in/gorp.v1"
	"sort"
)

//type SampleRequest struct {
//	Name string `form:"nameField"`
//}

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
func stringToAnswersMap(jsonStr string) map[string]string {
	x := map[string]string{}
	json.Unmarshal([]byte(jsonStr), &x)
	return x
}
func getMapFromSessioN(session sessions.Session) map[string]string {
	jsonStr := getResultsFromSessionAsJson(session)
	if jsonStr != "" {
		return stringToAnswersMap(jsonStr)
	} else {
		m := make(map[string]string)
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
	Score CompareResult
}

func getUserFromName(dbmap *gorp.DbMap, currentName string) User {
	var currentUser User
	dbmap.SelectOne(&currentUser, fmt.Sprintf("SELECT * FROM User WHERE Name='%s'", currentName))
	return currentUser
}

func getScores(dbmap *gorp.DbMap, currentName string) []Score {
	currentUser := getUserFromName(dbmap, currentName)

	var others []User
	dbmap.Select(&others, fmt.Sprintf("SELECT * FROM User WHERE Name!='%s'", currentName))

	var scores = make([]Score, len(others))

	for i, other := range others {
		scores[i] = Score{other.Name, compareUsers(currentUser, other)}
	}

	// sort per best matching
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score.Count > scores[j].Score.Count
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

type CompareUser struct {
	Article Article
	Same    bool
}

func getResultsFromMap(dbmap *gorp.DbMap, m map[string]string) ([]Result, float64) {
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
	//fmt.Printf("%.9f , %d , %d", progress, len(results), notanswer)
	return results, progress
}

// func getUserFromSession(session sessions.Session) (response *string) {
// is great but forces to return a pointer

func getUserFromSession(session sessions.Session) string {
	possibleCurrentUser := session.Get("user")
	if possibleCurrentUser == nil {
		return ""
		// return nil
	} else {
		name := fmt.Sprintf("%s", possibleCurrentUser)
		fmt.Printf("Using Username %s", name)
		return name
		// return &name
	}
}
