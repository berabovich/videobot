package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"unicode"
)

var quary = "робинзон"

type Films struct {
	LinkText string `json:"Название"`
	Link     string `json:"Ссылка"`
}
type FilmsData struct {
	Films []Films
}

func linkScrape() []Films {
	quStr := strings.ReplaceAll(quary, " ", "+")
	doc, err := goquery.NewDocument("https://www.kinopoisk.ru/index.php?kp_query=" + quStr)
	if err != nil {
		log.Fatal(err)
	}

	var films []Films
	doc.Find("body a").Each(func(index int, item *goquery.Selection) {
		var film Films
		linkTag := item
		link, _ := linkTag.Attr("href")
		linkText := linkTag.Text()
		quarySplit := strings.Split(quary, " ")
		linkTextSplit := strings.Split(linkText, " ")
		for _, q := range quarySplit {
			for _, l := range linkTextSplit {
				if strings.EqualFold(l, q) && strings.Contains(link, "film") {
					link = strings.ReplaceAll(link, "/sr/1/", "")
					sp := strings.Split(link, "/")
					if len(sp) < 4 && sp[1] == "film" {
						link = "https://www.kinopoiskk.ru" + link
						film.LinkText = linkText
						film.Link = link
						films = append(films, film)
					}

				}
			}
		}

	})
	return films
}

func main() {
	films := controller()
	rawDataOut, err := json.MarshalIndent(films, "", "  ")
	if err != nil {
		log.Fatal("JSON marshaling failed:", err)
	}
	if films.Films == nil {
		fmt.Println("Фильмы не найдены")
	} else {
		fmt.Println(string(rawDataOut))
	}
}

func controller() FilmsData {
	var fd FilmsData
	if IsEngByLoop(quary) {
		return fd
	}
	fd = readFile(&fd, quary)
	if fd.Films == nil {
		f := linkScrape()
		fd.Films = f
		writeFile(fd)
		return fd
	}
	return fd
}

func IsEngByLoop(str string) bool {
	if _, err := strconv.Atoi(str); err == nil {
		return false
	}
	for i := 0; i < len(str); i++ {
		if str[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func readFile(m *FilmsData, quary string) FilmsData {
	rawDataIn, _ := ioutil.ReadFile("films.json")
	_ = json.Unmarshal(rawDataIn, &m)
	var fd FilmsData
	var fd2 FilmsData
	for _, j := range m.Films {
		quarySplit := strings.Split(quary, " ")
		linkTextSplit := strings.Split(j.LinkText, " ")
		for _, q := range quarySplit {
			for _, t := range linkTextSplit {
				if strings.EqualFold(q, t) {
					fd.Films = append(fd.Films, j)
					break
				}
			}
			if fd.Films != nil {
				break
			}
		}
	}

	if len(strings.Split(quary, " ")) > 1 {
		for i, _ := range strings.Split(quary, " ") {
			i++
			for j, f := range fd.Films {
				if strings.Contains(f.LinkText, strings.Split(quary, " ")[i]) {
					fd2.Films = append(fd2.Films, fd.Films[j])
				}
			}
			break
		}
		return fd2
	}

	return fd
}

func writeFile(m FilmsData) {
	var data FilmsData
	rawDataIn, _ := ioutil.ReadFile("films.json")
	_ = json.Unmarshal(rawDataIn, &data)
	for _, j := range data.Films {
		m.Films = append(m.Films, j)
	}
	rawDataOut, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		log.Fatal("JSON marshaling failed:", err)
	}
	err = ioutil.WriteFile("films.json", rawDataOut, 0644)
	if err != nil {
		log.Fatal("Cannot write updated file:", err)
	}
}
