package main

import (
	"encoding/json"
)

func getUserAnswers(user User) map[string]string {
	jsonStr0 := user.Answers
	x := map[string]string{}
	json.Unmarshal([]byte(jsonStr0), &x)
	return x
}

/**
Compare how many common results between two users
and returns
*/
func compareUsers(user1 User, user2 User) CompareResult {
	answers1 := getUserAnswers(user1)
	answers2 := getUserAnswers(user2)

	count := 0
	for k, v := range answers1 {
		// fmt.Sprintf("> %d, %d\n", v, answers2[k])
		if answers2[k] == v {
			count++
		}
	}
	return CompareResult{count, len(answers2)}
}

type CompareResult struct {
	Count  int
	Theirs int
}
