package main

import (
	"encoding/json"
	"errors"
	"gopkg.in/gorp.v1"
)

type User struct {
	Name    string
	Answers string
}

type Answer struct {
	Title   string
	Answer1 string
	Answer2 string
	Answer3 string
}

type Article struct {
	Id           string
	MyIndex      int
	OriginalLink string
	Title        string
	ImageURL     string
	Summary      string
	Choice1      string
	Choice2      string
	Choice3      string
	Choice4      string
	ImageFile    string
	Answer1      Answer
	Answer2      Answer
	Answer3      Answer
	Answer4      Answer
}

type myTypeConverter struct{}

func (me myTypeConverter) ToDb(val interface{}) (interface{}, error) {

	switch t := val.(type) {
	case Answer:
		b, err := json.Marshal(t)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	return val, nil
}

func (me myTypeConverter) FromDb(target interface{}) (gorp.CustomScanner, bool) {
	switch target.(type) {
	case *Answer:
		binder := func(holder, target interface{}) error {
			s, ok := holder.(*string)
			if !ok {
				return errors.New("FromDb: Unable to convert Answer to *string")
			}
			b := []byte(*s)
			return json.Unmarshal(b, target)
		}
		return gorp.CustomScanner{new(string), target, binder}, true
	}

	return gorp.CustomScanner{}, false
}
