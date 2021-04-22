package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type Job struct {
	Company   string
	Date      time.Time
	BodyTitle string
	Body      string
}

func main() {
	url := "https://hnhiring.com/technologies/vue"

	content, err := loadPage(url)
	if err != nil {
		panic(err)
	}

	jobs, err := findJobs(content)
	if err != nil {
		panic(err)
	}

	fmt.Println(len(jobs))

	print(jobs)
}

func loadPage(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("URL [%v]: %v", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error parseando body: %v", err)
	}

	return string(body), nil
}

func findJobs(content string) ([]Job, error) {
	var jobs []Job

	doc := parseHTML(content)

	jobsPosts := doc.Find("ul", "class", "jobs").FindAll("li", "class", "job")

	for _, v := range jobsPosts {

		// parseando tiempo
		t, err := time.Parse("2006-01-02", v.Find("div", "class", "user").Find("span", "class", "gray").Value())
		if err != nil {
			return nil, err
		}

		body := v.Find("div", "class", "body").FindAll("p", "", "")

		// agregando post
		jobs = append(jobs, Job{
			Company:   v.Find("div", "class", "user").Find("a", "href", "https").Value(),
			Date:      t,
			BodyTitle: body[0].Value(),
			Body:      parsingBody(body[1:]),
		})
	}

	return jobs, nil
}

func parseHTML(str string) *HTMLStruct {
	html, err := html.Parse(strings.NewReader(str))
	if err != nil {
		panic(err)
	}

	return &HTMLStruct{*html}
}

func parsingBody(content []HTMLStruct) string {
	var body string

	for _, v := range content {
		body += v.Value()
	}

	return body
}

func print(data []Job) {
	for _, v := range data {
		fmt.Printf("Company: %s \n", v.Company)
		fmt.Printf("Date: %s \n", v.Date)
		fmt.Printf("Body Title: %s \n", v.BodyTitle)
		fmt.Printf("Body: %s \n", v.Body)
	}
}

type HTMLStruct struct {
	html.Node
}

func (h HTMLStruct) Find(tag, atr, value string) HTMLStruct {
	result, ok := findOne(h, tag, atr, value)
	if !ok {
		panic("No se ha encontrado ningun elemento con el etiqueta " + tag)
	}
	return result
}

func findOne(content HTMLStruct, tag, atr, value string) (HTMLStruct, bool) {
	if content.Type == html.ElementNode && content.Data == tag {
		for _, v := range content.Attr {
			if (v.Key == atr && v.Val == value) || atrNameContain(v, atr, value) {
				return content, true
			} else {
				return content, true
			}
		}
	}
	for c := content.FirstChild; c != nil; c = c.NextSibling {
		p, ok := findOne(HTMLStruct{*c}, tag, atr, value)
		if ok {
			return p, ok
		}
	}
	return HTMLStruct{}, false
}

func atrNameContain(attr html.Attribute, atr, value string) bool {
	if attr.Key == atr {
		for _, attrVal := range strings.Fields(attr.Val) {
			if attrVal == value {
				return true
			}
		}
	}
	return false
}

func (h HTMLStruct) FindAll(tag, atr, value string) []HTMLStruct {
	var elemts []HTMLStruct

	result := findAll(&h.Node, tag, atr, value)
	if len(result) == 0 {
		panic("No se ha encontrado ningun elemento con el etiqueta " + tag)
	}

	for _, v := range result {
		elemts = append(elemts, HTMLStruct{*v})
	}
	return elemts
}

func findAll(cont *html.Node, tag, atr, value string) []*html.Node {
	var nodes []*html.Node

	var f func(n *html.Node, tag, atr, value string)
	f = func(n *html.Node, tag, atr, value string) {
		if n.Type == html.ElementNode && (n.Data == tag) {
			if atr != "" && value != "" {
				for _, v := range n.Attr {
					if (v.Key == atr && v.Val == value) || atrNameContain(v, atr, value) {
						nodes = append(nodes, n)
					}
				}
			} else {
				nodes = append(nodes, n)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c, tag, atr, value)
		}
	}

	f(cont, tag, atr, value)
	return nodes
}

func (h HTMLStruct) Value() string {
	child := h.FirstChild
checkNode:
	if child != nil && child.Type != html.TextNode {
		child = child.FirstChild
		if child == nil {
			panic("ningun valor encontrado")
		}
		goto checkNode
	}
	if child != nil {
		r, _ := regexp.Compile(`^\s+$`)
		if ok := r.MatchString(child.Data); ok {
			child = child.NextSibling
			if child == nil {
				panic("ningun valor encontrado")
			}
			goto checkNode
		}
		return child.Data
	}
	return ""
}
