package main

import (
	"encoding/json"
	"fmt"
	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
	"github.com/psykhi/wordclouds"
	"image"
	"image/color"
	"image/png"
	"os"
)

func answerToExtendedText(answer Answer) string {
	return answer.Title + " " + answer.Answer1 + " " + answer.Answer2 + " " + answer.Answer3
}
func getTextForUser(userName string) []string {
	user := getUserFromName(dbmap, userName)
	m := getUserAnswers(user)

	articles := getAllArticles(dbmap)

	results := make([]string, 0)
	for i := range articles {
		choiceI := m[articles[i].Id]
		temp := ""
		switch choiceI {
		case "1":
			temp = answerToExtendedText(articles[i].Answer1)
		case "2":
			temp = answerToExtendedText(articles[i].Answer2)
		case "3":
			temp = answerToExtendedText(articles[i].Answer3)
		case "4":
			temp = answerToExtendedText(articles[i].Answer4)
		}
		if temp != "" {
			results = append(results, temp)
		}

	}

	return results
}

func groupBy(arr []string) map[string]int {
	dict := make(map[string]int)
	for _, num := range arr {
		if dict[num] == 0 {
			dict[num] = 1
		} else {
			dict[num] = dict[num] + 1
		}
	}
	return dict
}

func simpleMap(wordCounts map[string]int) image.Image {
	var RedColors = []color.RGBA{
		{0xe5, 0x28, 0x2a, 0xff},
		{0x9a, 0x24, 0x1a, 0xff},
		{0xea, 0x87, 0x1a, 0xff},
		{0xff, 0x6f, 0x1a, 0xff},
		{0xff, 0x8f, 0x66, 0xff},
	}

	colors := make([]color.Color, 0)
	for _, c := range RedColors {
		colors = append(colors, c)
	}

	w := wordclouds.NewWordcloud(
		wordCounts,
		wordclouds.FontFile("rounded-l-mplus-2c-medium.ttf"),
		wordclouds.Height(2048),
		wordclouds.Width(2048),
		wordclouds.Colors(colors),
	)

	img := w.Draw()

	return img

}

func writeImgToFile(img image.Image) string {
	filePath := "output.png"
	outputFile, err := os.Create(filePath)
	if err != nil {
		// Handle error
		fmt.Println(err)
	}

	// Encode takes a writer interface and an image interface
	// We pass it the File and the RGBA
	png.Encode(outputFile, img)

	// Don't forget to close files
	outputFile.Close()
	return filePath
}

func mapForUser(userName string) image.Image {
	test := getTextForUser(userName)

	tokenizer, _ := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	tokenss := make([]string, 0)

	for _, line := range test {
		tokens := tokenizer.Tokenize(line)
		for _, token := range tokens {
			if token.Features()[0] == "名詞" {
				tokenss = append(tokenss, token.Surface)
			}
		}
	}

	by := groupBy(tokenss)

	bs, _ := json.Marshal(by)
	fmt.Println(string(bs))

	return simpleMap(by)
}
