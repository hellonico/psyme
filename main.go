package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v1"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//
//func findAverage(a []int) float64 {
//	count := 4
//	sum := 0
//	for i := 0; i < count; i++ {
//		sum += (a[i])
//	}
//
//	return float64(sum) / float64(count)
//}
//

func fetch() {
	content, err := ioutil.ReadFile("urls.txt")
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(content), "\n")

	dbmap := initDb()
	defer dbmap.Db.Close()

	for i, line := range lines {
		if !strings.HasPrefix(line, "#") {
			processOne(i, line, dbmap)
		}
	}

}

func generate() {
	// READ FROM DB
	dbmap := initDb()
	defer dbmap.Db.Close()

	//file, err := os.OpenFile("Index.html", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	file, err := os.OpenFile("Index.html", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// res, err := dbmap.Get(Article{}, a.Id)
	// log.Printf("SQL: %d", res)
	// addToSummary(a, file)

	var articles []Article
	dbmap.Select(&articles, "Select * from Article")
	for i, article := range articles {
		if i == 0 {
			fmt.Fprint(file, "<h1>Index</h1><table>")
		}
		log.Printf("> %s", article.Id)
		article.html()
		addToSummary(article, file)
	}

	fmt.Fprint(file, "</table>")
}

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		fmt.Printf("Usage: psy <command>\n")
		fmt.Printf("Where: command is\n")
		fmt.Printf("--fetch \t: Fetch Articles (incremental) and populate database\n")
		fmt.Printf("--generate\t: Generate Static Website\n")
		fmt.Printf("--cleanHtml\t: Clean Static Website\n")
		fmt.Printf("--cleanDb\t: Clean Database\n")
		fmt.Printf("--regen \t: Clean and Generate Website\n")
		fmt.Printf("--one     \t: Fetch one\n")
		fmt.Printf("--web     \t: Start web server\n")
		fmt.Printf("--images  \t: Download images\n")
	} else {
		switch argsWithoutProg[0] {
		case "--fetch":
			fetch()
		case "--generate":
			generate()
		case "--cleanDb":
			cleanDb()
		case "--cleanHtml":
			cleanHtml()
		case "--regen":
			cleanHtml()
			generate()
		case "--one":
			a := parseArticle(argsWithoutProg[1])
			a.markdown()
		case "--web":
			web()
		case "--images":
			images()
		}
	}
}

func images() {
	dbmap := initDb()
	defer dbmap.Db.Close()

	var articles []Article
	dbmap.Select(&articles, "Select * from Article")
	for _, article := range articles {
		filePath := fmt.Sprintf("assets/img/%s.png", article.Id)
		fmt.Printf("< %s\n", filePath)
		downloadFile(article.ImageURL, filePath)
	}
}

func cleanDb() {
	dbmap := initDb()
	dbmap.Db.Exec("drop table if exists Article")
}

func cleanHtml() {
	fmt.Printf("Clean html files...\n")
	files, _ := filepath.Glob("./*.html")
	for _, f := range files {
		os.Remove(f)
	}
	fmt.Printf("Clean md files...\n")
	files, _ = filepath.Glob("./*.md")
	for _, f := range files {
		os.Remove(f)
	}
	//fmt.Printf("Clean png files...\n")
	//files, _ = filepath.Glob("./*.png")
	//for _, f := range files {
	//	os.Remove(f)
	//}
}

//parseArticle("https://trilltrill.jp/articles/2126816")
//parseArticle("https://trilltrill.jp/articles/2126770")

func initDb() *gorp.DbMap {
	//init db
	path, _ := os.Getwd()
	dbFileName := fmt.Sprintf("%s/articles.sqlite", path)
	log.Printf("Db File is: %s", dbFileName)
	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		log.Fatalf("Cannot open databse: %s", dbFileName)
		panic(err)
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbmap.TypeConverter = myTypeConverter{}

	dbmap.AddTable(Article{}).SetKeys(false, "Id")
	dbmap.AddTable(User{}).SetKeys(false, "Name")
	err = dbmap.CreateTablesIfNotExists()
	return dbmap
}

//func checkExistFirst(id string, i int, file *os.File, line string, dbmap *gorp.DbMap) {
//	if _, err := os.Stat(fmt.Sprintf("%s.md", id)); err == nil {
//		// exist do not process again
//		fmt.Printf("< already here\n")
//	} else if errors.Is(err, os.ErrNotExist) {
//		processOne(line, dbmap)
//	} else {
//		// Schrodinger: file may or may not exist. See err for details.
//		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
//	}
//}

func processOne(index int, line string, dbmap *gorp.DbMap) {
	id := idFromUrl(line)
	obj, _ := dbmap.Get(Article{}, id)
	fmt.Printf("Processing > %s : ", id)
	if obj != nil {
		fmt.Printf("< skip\n")
	} else {
		fmt.Printf("< download\n")
		a := parseArticle(line)
		a.MyIndex = index

		err := dbmap.Insert(&a)
		if err != nil {
			log.Fatalf("Cannot do database insert %s\n", err)
		}
	}

}

func addToSummary(a Article, file *os.File) {
	fmt.Fprint(file, fmt.Sprintf("<tr><td><a href=\"%s.html\">%s</a></td><td><img src=\"%s\"/></td></tr>\n", a.Id, a.Title, a.ImageFile))
}
