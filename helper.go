package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"os"
)

func downloadFile(URL, fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}
	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func idFromUrl(url string) string {
	return url[len(url)-7 : len(url)]
}

func parseArticle(url string) Article {

	a := Article{}

	a.Id = idFromUrl(url)
	a.OriginalLink = url

	resp, err := http.Get(url)
	if err != nil {
		print(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Title
	_title := doc.Find("Title").First().Text()
	a.Title = _title[:len(_title)-35]

	// image
	s := doc.Find("img.image-block__image").First()
	a.ImageURL, _ = s.Attr("data-src")
	a.ImageFile = fmt.Sprintf("%s.png", a.Id)
	// downloadFile(a.ImageURL, a.ImageFile)

	// Summary
	a.Summary = doc.Find(".trill-description > p:nth-child(2)").First().Text()

	//fmt.Printf("%s\n", a.Summary)

	testField := doc.Find(".trill-description > p:nth-child(10)").First().Text()
	//fmt.Printf("1\n")
	//fmt.Printf("%s [%d]\n", testField, len(testField))

	testField2 := doc.Find(".trill-description > p:nth-child(6)").First().Text()
	//fmt.Printf("2\n")
	//fmt.Printf("%s [%d]\n", testField2, len(testField2))

	pattern := 1
	if len(testField) != 0 {
		//fmt.Printf("WHYYYYY?\n")
		//fmt.Printf("%s [%d]\n", testField, len(testField))
		pattern = 2
	}
	if testField2 == "" {
		pattern = 3
	}
	//fmt.Printf("pattern %d\n", pattern)

	// question
	switch pattern {
	case 1:
		a.Choice1 = doc.Find(".trill-description > p:nth-child(6)").First().Text()[2:]
		a.Choice2 = doc.Find(".trill-description > p:nth-child(7)").First().Text()[2:]
		a.Choice3 = doc.Find(".trill-description > p:nth-child(8)").First().Text()[2:]
		a.Choice4 = doc.Find(".trill-description > p:nth-child(9)").First().Text()[2:]
	case 2:
		a.Choice1 = doc.Find(".trill-description > p:nth-child(5)").First().Text()[2:]
		a.Choice2 = doc.Find(".trill-description > p:nth-child(6)").First().Text()[2:]
		a.Choice3 = doc.Find(".trill-description > p:nth-child(7)").First().Text()[2:]
		a.Choice4 = doc.Find(".trill-description > p:nth-child(8)").First().Text()[2:]
	case 3:
		a.Choice1 = doc.Find(".trill-description > p:nth-child(7)").First().Text()[2:]
		a.Choice2 = doc.Find(".trill-description > p:nth-child(8)").First().Text()[2:]
		a.Choice3 = doc.Find(".trill-description > p:nth-child(9)").First().Text()[2:]
		a.Choice4 = doc.Find(".trill-description > p:nth-child(10)").First().Text()[2:]
	}

	// answers
	a.Answer1 = parseAnswer(doc, 1, pattern)
	a.Answer2 = parseAnswer(doc, 2, pattern)
	a.Answer3 = parseAnswer(doc, 3, pattern)
	a.Answer4 = parseAnswer(doc, 4, pattern)

	return a

	/* the whole text */
	//doc.Find(".trill-description").Each(func(i int, s *goquery.Selection) {
	//	fmt.Printf("%s", s.Text())
	//})
}

func parseAnswer(doc *goquery.Document, i int, pattern int) Answer {
	answer := Answer{}

	// 12, 16, 20 ,24
	// or
	// 11,15,19,23
	var startIndex = 8
	switch pattern {
	case 2:
		startIndex = 7
	case 3:
		startIndex = 9
	}
	//.trill-description > h3:nth - child(14)
	// startIndex+4*i
	//for i := 10; i < 30; i++ {
	//	fmt.Printf("[%d] %s\n", i, doc.Find(fmt.Sprintf(".trill-description > h3:nth-child(%d)", i)).First().Text())
	//}
	h3 := fmt.Sprintf(".trill-description > h3:nth-child(%d)", startIndex+4*i)
	// fmt.Printf("pattern : %d [%s]\n", pattern, h3)

	// 13,14,15
	// 17,18,19
	// 21,22,23
	// 25,26,27
	p1 := fmt.Sprintf(".trill-description > p:nth-child(%d)", startIndex+4*i+1)
	p2 := fmt.Sprintf(".trill-description > p:nth-child(%d)", startIndex+4*i+2)
	p3 := fmt.Sprintf(".trill-description > p:nth-child(%d)", startIndex+4*i+3)

	answer.Title = doc.Find(h3).First().Text()[2:]
	answer.Answer1 = doc.Find(p1).First().Text()
	answer.Answer2 = doc.Find(p2).First().Text()
	answer.Answer3 = doc.Find(p3).First().Text()
	return answer
}
