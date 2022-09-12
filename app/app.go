package app

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type Films struct {
	LinkText string `json:"Название"`
	Link     string `json:"Ссылка"`
}
type FilmsData struct {
	Films []Films
}

func linkScrape(quary string) []Films {
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
		quary = strCompl(quary)
		quarySplit := strings.Split(quary, " ")
		linkTxt := strCompl(linkText)
		linkTextSplit := strings.Split(linkTxt, " ")
		for _, q := range quarySplit {
			for _, l := range linkTextSplit {
				if strings.EqualFold(l, q) || strings.Contains(link, "film") && l != "" {
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

func Controller(quary string) FilmsData {
	var fd FilmsData
	if IsEngByLoop(quary) {
		return fd
	}
	fd = readFile(&fd, quary)
	if fd.Films == nil {
		f := linkScrape(quary)
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

func strCompl(str string) string {
	var re = regexp.MustCompile(`[[:punct:]]`)
	str = re.ReplaceAllString(str, "")
	return str
}

func readFile(m *FilmsData, quary string) FilmsData {
	rawDataIn, _ := ioutil.ReadFile("films.json")
	_ = json.Unmarshal(rawDataIn, &m)
	var fd FilmsData
	var fd2 FilmsData
	for _, j := range m.Films {
		quary = strCompl(quary)
		quarySplit := strings.Split(quary, " ")
		linkTxt := strCompl(j.LinkText)
		linkTextSplit := strings.Split(linkTxt, " ")
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
				str := strings.Split(quary, " ")
				if strings.Contains(strings.ToLower(f.LinkText), strings.ToLower(str[i])) {
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
