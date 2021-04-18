package main

import (
	"fmt"
	"time"

	"github.com/anaskhan96/soup"
)

type Job struct {
	Company   string
	Date      time.Time
	BodyTitle string
	Body      string
}

func main() {
	url := "https://hnhiring.com/technologies/go"

	content, err := loadPage(url)
	if err != nil {
		panic(err)
	}

	// fmt.Println(string(content))
	jobs, err := findJobs(content)
	if err != nil {
		panic(err)
	}

	print(jobs)
}

func loadPage(url string) (string, error) {
	resp, err := soup.Get(url)
	if err != nil {
		return "", fmt.Errorf("URL [%v]: %v", url, err)
	}

	return resp, nil
}

func findJobs(content string) ([]Job, error) {
	var jobs []Job

	doc := soup.HTMLParse(content)

	jobsPosts := doc.Find("ul", "class", "jobs").FindAll("li")

	for _, v := range jobsPosts {

		// parseando tiempo
		t, err := time.Parse("2006-01-02", v.Find("div", "class", "user").Find("span").Text())
		if err != nil {
			return nil, err
		}

		// agregando post
		jobs = append(jobs, Job{
			Company:   v.Find("div", "class", "user").Find("a").Text(),
			Date:      t,
			BodyTitle: v.Find("div", "class", "body").FindAll("p")[0].Text(),
			Body:      parsingBody(v.Find("div", "class", "body").FindAll("p")[1:]),
		})
	}

	return jobs, nil
}

func parsingBody(content []soup.Root) string {
	var body string

	for _, v := range content {
		body += v.Text()
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
