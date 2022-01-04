package main

import (
	"encoding/json"
	"fmt"
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
	answers2 := getUserAnswers(user2)

	// fmt.Printf("| %d\n", len(answers1))
	// fmt.Printf("| %d\n", len(answers2))
	count := 0
	for k, v := range answers1 {
		fmt.Sprintf("> %d, %d\n", v, answers2[k])
		if answers2[k] == v {
			count++
		}
	}
	return count
}
