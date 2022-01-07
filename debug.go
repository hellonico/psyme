package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

/**
Methods only used for debugging now
*/
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

func printArrayOfString(test []string) {
	semiformat := fmt.Sprintf("%q\n", test)  // Turn the slice into a string that looks like ["one" "two" "three"]
	tokens := strings.Split(semiformat, " ") // Split this string by spaces
	fmt.Printf(strings.Join(tokens, ", "))   // Join the Slice together (that was split by spaces) with commas
}
