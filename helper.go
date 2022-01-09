package main

type StrSlice []string

func (list StrSlice) Has(a string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
