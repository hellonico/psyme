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
