package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/gorp.v1"
	"io"
	"os"
	"strings"
)

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

func (a Article) markdown() {

	fileName := fmt.Sprintf("%s.md", a.Id)
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fmt.Fprint(file, fmt.Sprintf("# %s\n", a.Title))
	fmt.Fprint(file, fmt.Sprintf("__[%s](%s)__\n", a.OriginalLink, a.OriginalLink))
	fmt.Fprint(file, fmt.Sprintf("![](./%s)\n", a.ImageFile))
	fmt.Fprint(file, "\n\n")

	fmt.Fprint(file, fmt.Sprintf("%s\n\n", a.Summary))
	fmt.Fprint(file, fmt.Sprintf("- %s\n", a.Choice1))
	fmt.Fprint(file, fmt.Sprintf("- %s\n", a.Choice2))
	fmt.Fprint(file, fmt.Sprintf("- %s\n", a.Choice3))
	fmt.Fprint(file, fmt.Sprintf("- %s\n", a.Choice4))

	as := []Answer{a.Answer1, a.Answer2, a.Answer3, a.Answer4}

	for _, ab := range as {
		fmt.Fprint(file, fmt.Sprintf("## %s\n\n", ab.Title))
		fmt.Fprint(file, fmt.Sprintf("%s\n\n", ab.Answer1))
		fmt.Fprint(file, fmt.Sprintf("%s\n\n", ab.Answer2))
		fmt.Fprint(file, fmt.Sprintf("%s\n\n", ab.Answer3))
	}

	// print
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b := new(strings.Builder)
	io.Copy(b, f)
	print(b.String())

}

func (a Article) html() {
	// prepare file
	file, err := os.Create(fmt.Sprintf("%s.html", a.Id))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fmt.Fprint(file, fmt.Sprintf("<h1>%s</h1>\n", a.Title))
	fmt.Fprint(file, fmt.Sprintf("<a href=\"%s\">%s</a><br/>\n", a.OriginalLink, a.OriginalLink))
	fmt.Fprint(file, fmt.Sprintf("<img src=\"%s\"/><br/>\n", a.ImageFile))
	fmt.Fprint(file, "\n\n")

	fmt.Fprint(file, fmt.Sprintf("<h2>Summary</h2>\n%s\n", a.Summary))
	fmt.Fprint(file, fmt.Sprintf("<h2>choices</h2>\n"))

	fmt.Fprint(file, fmt.Sprintf("<ul><li>%s</li><li>%s</li><li>%s</li><li>%s</li></ul>\n", a.Choice1, a.Choice2, a.Choice3, a.Choice4))

	fmt.Fprint(file, fmt.Sprintf("<h2>answers</h2>\n"))
	as := []Answer{a.Answer1, a.Answer2, a.Answer3, a.Answer4}
	for _, ab := range as {
		fmt.Fprint(file, fmt.Sprintf("<h3>%s</h3>\n\n", ab.Title))
		fmt.Fprint(file, fmt.Sprintf("<p>%s</p>\n\n", ab.Answer1))
		fmt.Fprint(file, fmt.Sprintf("<p>%s</p>\n\n", ab.Answer2))
		fmt.Fprint(file, fmt.Sprintf("<p>%s</p>\n\n", ab.Answer3))
	}

}
