package main

import (
	"encoding/json"
	"gopkg.in/gorp.v1"
)

func getUserAnswers(user User) map[string]string {
	jsonStr0 := user.Answers
	x := map[string]string{}
	json.Unmarshal([]byte(jsonStr0), &x)
	return x
}
func compareUsers(dbMap *gorp.DbMap, user1 User, user2 User) int {
	answers1 := getUserAnswers(user1)
	answers2 := getUserAnswers(user1)

	count := 0
	for k, v := range answers1 {
		if answers2[k] == v {
			count++
		}
	}
	return count
}
